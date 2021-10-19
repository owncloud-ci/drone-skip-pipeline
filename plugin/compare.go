package plugin

import (
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

// Release holds ties the drone env data and github client together.
type compare struct {
	gitPath string

	commitSHAafter  string
	commitSHAbefore string

	repo *git.Repository

	changed []string

	disallowSkipChanged []string
	allowSkipChanged    []string
}

func newCompare(
	gitPath, commitSHAafter, commitSHAbefore string,
	disallowSkipChanged, allowSkipChanged []string,
) compare {
	return compare{
		gitPath:             gitPath,
		commitSHAafter:      commitSHAafter,
		commitSHAbefore:     commitSHAbefore,
		disallowSkipChanged: disallowSkipChanged,
		allowSkipChanged:    allowSkipChanged,
	}
}

func (c *compare) open() (err error) {
	c.repo, err = git.PlainOpen(c.gitPath)
	if err != nil {
		return err
	}
	return nil
}

func (c *compare) getChanged() error {

	c.commitSHAbefore = "c2c198c3f38034b62138618796ef1af0a47a86d5"
	c.commitSHAafter = "debaf34326530e8dce844215140b53c36557fe39"

	beforeHash := plumbing.NewHash(c.commitSHAbefore)
	beforeCommit, err := c.repo.CommitObject(beforeHash)
	if err != nil {
		return errors.Wrap(err, "commit from DRONE_COMMIT_BEFORE not found")
	}

	parentCommit, err := beforeCommit.Parent(0)
	if err != nil {
		return errors.Wrap(err, "parent not found")
	}

	afterHash := plumbing.NewHash(c.commitSHAafter)
	afterCommit, err := c.repo.CommitObject(afterHash)
	if err != nil {
		return errors.Wrap(err, "commit from DRONE_COMMIT_AFTER not found")
	}

	diff, err := parentCommit.Patch(afterCommit)
	if err != nil {
		return errors.Wrap(err, "could not get diff")
	}

	fmt.Println("### changed files ###")
	fmt.Println(diff.Stats().String())

	fileStats := diff.Stats()

	changed := []string{}

	for _, f := range fileStats {
		changed = append(changed, f.Name)
	}

	c.changed = changed
	return nil
}

func (c *compare) isSkip() (skip bool, err error) {

	if len(c.changed) == 0 {
		return true, nil
	}

	allowSkipChanged := false
	disallowSkipChanged := false

	if len(c.allowSkipChanged) == 0 {
		allowSkipChanged = true
	} else {
		allowSkipChanged, err = allowSkipCompare(c.changed, c.allowSkipChanged)
		if err != nil {
			return false, err
		}
	}

	if len(c.disallowSkipChanged) == 0 {
		disallowSkipChanged = false
	} else {
		disallowSkipChanged, err = disallowSkipCompare(c.changed, c.disallowSkipChanged)
		if err != nil {
			return false, err
		}
	}

	skip = allowSkipChanged && !disallowSkipChanged

	return skip, nil
}

func disallowSkipCompare(strings, regexes []string) (bool, error) {
	fmt.Println("### check if disallowed file was changed ###")
	res := []*regexp.Regexp{}
	for _, r := range regexes {
		re, err := regexp.Compile(r)
		if err != nil {
			return false, errors.Wrap(err, fmt.Sprintf("no valid regex expression: '%s'", re))
		}
		res = append(res, re)
	}

	disallowSkip := false

	for _, re := range res {
		for _, s := range strings {
			if re.MatchString(s) {
				fmt.Println(" - '", s, "'is not allowed to be skipped because of '", re.String(), "'")
				disallowSkip = disallowSkip || true
			}
		}

	}
	if !disallowSkip {
		fmt.Println(" - no file disallowed to be skipped was changed")
	}

	return disallowSkip, nil
}

func allowSkipCompare(strings, regexes []string) (bool, error) {
	fmt.Println("### check if changed files are on allowed skip list ###")
	res := []*regexp.Regexp{}
	for _, r := range regexes {
		re, err := regexp.Compile(r)
		if err != nil {
			return false, errors.Wrap(err, fmt.Sprintf("no valid regex expression: '%s'", re))
		}
		res = append(res, re)
	}

	skip := true

	for _, s := range strings {
		fileSkip := false
		for _, re := range res {
			fileSkip = fileSkip || re.MatchString(s)
		}

		// one file change was not skipable
		if !fileSkip {
			fmt.Println(" - '", s, "' is not allowed to be skipped because it didn't match any skip rule")
			skip = skip && false
		}
	}

	if skip {
		fmt.Println(" - all changed files are allowed to be skipped")
	}

	return skip, nil
}
