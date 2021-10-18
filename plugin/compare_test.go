package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidRegex(t *testing.T) {
	c := compare{
		changed: []string{
			".drone.star",
			"cmd/foo.go",
			"pkg/internal/bar.go",
		},
		allowSkipChanged: []string{
			"*", // this is invalid regex syntax
		},
		disallowSkipChanged: []string{},
	}

	isSkip, err := c.isSkip()
	assert.NotNil(t, err)
	assert.False(t, isSkip)
}

func TestCodeAllowSkipDocs1(t *testing.T) {
	c := compare{
		changed: []string{
			".drone.star",
			"cmd/foo.go",
			"pkg/internal/bar.go",
			"docs/config.hugo",
		},
		allowSkipChanged: []string{
			`^docs/.*`,
		},
		disallowSkipChanged: []string{},
	}

	isSkip, err := c.isSkip()
	assert.Nil(t, err)
	assert.False(t, isSkip)
}

func TestCodeAllowSkipDocs2(t *testing.T) {
	c := compare{
		changed: []string{
			"docs/index.md",
		},
		allowSkipChanged: []string{
			`^docs/.*`,
		},
		disallowSkipChanged: []string{},
	}
	isSkip, err := c.isSkip()
	assert.Nil(t, err)
	assert.True(t, isSkip)
}

func TestCodeDisallowSkipGoChanges1(t *testing.T) {
	c := compare{
		changed: []string{
			".drone.star",
			"cmd/foo.go",
			"pkg/internal/bar.go",
			"docs/config.hugo",
		},
		allowSkipChanged: []string{},
		disallowSkipChanged: []string{
			`.*\.go$`,
		},
	}

	isSkip, err := c.isSkip()
	assert.Nil(t, err)
	assert.False(t, isSkip)
}

func TestCodeDisallowSkipGoChanges2(t *testing.T) {
	c := compare{
		changed: []string{
			"docs/index.md",
			"docs/config.hugo",
		},
		allowSkipChanged: []string{},
		disallowSkipChanged: []string{
			`.*\.go$`,
		},
	}

	isSkip, err := c.isSkip()
	assert.Nil(t, err)
	assert.True(t, isSkip)
}

func TestNoChangesDisallowList(t *testing.T) {
	c := compare{
		changed:          []string{},
		allowSkipChanged: []string{},
		disallowSkipChanged: []string{
			`.*\.go$`,
		},
	}

	isSkip, err := c.isSkip()
	assert.Nil(t, err)
	assert.True(t, isSkip)
}

func TestNoChangesAllowList(t *testing.T) {
	c := compare{
		changed: []string{},
		allowSkipChanged: []string{
			`.*\.go$`,
		},
		disallowSkipChanged: []string{},
	}

	isSkip, err := c.isSkip()
	assert.Nil(t, err)
	assert.True(t, isSkip)
}
