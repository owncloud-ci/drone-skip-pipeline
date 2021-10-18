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
	gitPath      string
	targetBranch string

	repo *git.Repository

	changed []string

	disallowSkipChanged []string
	allowSkipChanged    []string
}

func newCompare(gitPath, targetBranch string, disallowSkipChanged, allowSkipChanged []string) compare {
	return compare{
		gitPath:             gitPath,
		targetBranch:        targetBranch,
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

	headRef, err := c.repo.Head()
	if err != nil {
		return errors.Wrap(err, "could not open git repo")
	}

	headCommit, err := c.repo.CommitObject(headRef.Hash())
	if err != nil {
		return errors.Wrap(err, "could not find HEAD")
	}

	targetRefName := plumbing.NewBranchReferenceName(c.targetBranch)
	if err != nil {
		return errors.Wrap(err, "could not reference target branch")
	}

	targetRef, err := c.repo.Reference(targetRefName, false)
	if err != nil {
		return errors.Wrap(err, "could not resolve target branch")
	}

	targetCommit, err := c.repo.CommitObject(targetRef.Hash())
	if err != nil {
		return errors.Wrap(err, "target commit not found")
	}

	diff, err := headCommit.Patch(targetCommit)
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

	for _, re := range res {
		disallowRule := true
		for _, s := range strings {
			disallowRule = disallowRule && re.MatchString(s)

			if re.MatchString(s) {
				fmt.Printf(" - '%s' is not allowed to be skipped because of '%s'", s, re.String())
			}
		}

		// one disallow rule triggered
		if !disallowRule {
			return false, nil
		}
	}

	fmt.Println("   - no disallowed file was changed")
	// no disallow rule triggered
	return true, nil
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
			fmt.Printf(" - '%s' is not allowed to be skipped because it didn't match any skip rule", s)
			skip = skip && false
		}
	}

	if skip {
		fmt.Println(" - all changed files are allowed to be skipped")
	}

	// all file changes are skipable
	return skip, nil
}
