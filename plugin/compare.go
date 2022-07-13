package plugin

import (
	"fmt"
	"regexp"
	"strings"

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

	fmt.Println("DRONE_COMMIT_BEFORE: ", c.commitSHAbefore)

	beforeHash := plumbing.NewHash(c.commitSHAbefore)
	beforeCommit, err := c.repo.CommitObject(beforeHash)
	if err != nil {
		return errors.Wrap(err, "commit from DRONE_COMMIT_BEFORE not found")
	}

	fmt.Println("DRONE_COMMIT_AFTER: ", c.commitSHAafter)

	afterHash := plumbing.NewHash(c.commitSHAafter)
	afterCommit, err := c.repo.CommitObject(afterHash)
	if err != nil {
		return errors.Wrap(err, "commit from DRONE_COMMIT_AFTER not found")
	}

	// https://stackoverflow.com/a/7256391
	// get merge base to produce a PR style diff
	mergeBaseCommits, err := beforeCommit.MergeBase(afterCommit)
	if err != nil || len(mergeBaseCommits) < 1 {
		return errors.Wrap(err, "could not find common merge base")
	}

	// this equals an `git diff before...after`
	diff, err := mergeBaseCommits[0].Patch(afterCommit)
	if err != nil {
		return errors.Wrap(err, "could not get diff")
	}

	fileStats := diff.Stats()
	changed := []string{}

	fmt.Println("### changed files ###")
	for _, f := range fileStats {
		if f.Name != "" {
			fmt.Print(f.String())
			changed = append(changed, strings.Trim(f.Name, "\""))
		}
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
				fmt.Println(" - '" + s + "'is not allowed to be skipped because of '" + re.String() + "'")
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
			fmt.Println(" - '" + s + "' is not allowed to be skipped")
			skip = skip && false
		}
	}

	if skip {
		fmt.Println(" - all changed files are allowed to be skipped")
	}

	return skip, nil
}
