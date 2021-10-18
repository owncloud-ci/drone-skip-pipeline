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

	runChanged  []string
	skipChanged []string
}

func NewCompare(gitPath, targetBranch string, runChanged, skipChanged []string) compare {
	return compare{
		gitPath:      gitPath,
		targetBranch: targetBranch,
		runChanged:   runChanged,
		skipChanged:  skipChanged,
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

	fileStats := diff.Stats()

	changed := []string{}

	for _, f := range fileStats {
		changed = append(changed, f.Name)
	}

	c.changed = changed
	return nil
}

func (c *compare) isSkip() (skip bool, err error) {

	skipChanged := false
	runChanged := false

	if len(c.skipChanged) == 0 {
		skipChanged = true
	} else {
		skipChanged, err = changed(c.changed, c.skipChanged)
		if err != nil {
			return false, err
		}
	}

	if len(c.runChanged) == 0 {
		runChanged = false
	} else {
		runChanged, err = changed(c.changed, c.runChanged)
		if err != nil {
			return false, err
		}
	}

	skip = skipChanged && !runChanged

	return skip, nil
}

func changed(strings, regexes []string) (bool, error) {
	for _, r := range regexes {
		re, err := regexp.Compile(r)
		if err != nil {
			return false, errors.Wrap(err, fmt.Sprintf("no valid regex expression: '%s'", re))
		}
		for _, s := range strings {
			if re.MatchString(s) {
				return true, nil
			}
		}
	}
	return false, nil
}
