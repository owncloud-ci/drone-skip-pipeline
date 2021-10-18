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
		skipChanged: []string{
			"*",
		},
		runChanged: []string{},
	}

	isSkip, err := c.isSkip()
	assert.NotNil(t, err)
	assert.False(t, isSkip)
}

func TestCodeChangedSkipRule(t *testing.T) {
	c := compare{
		changed: []string{
			".drone.star",
			"cmd/foo.go",
			"pkg/internal/bar.go",
		},
		skipChanged: []string{
			`^docs/.*`,
		},
		runChanged: []string{},
	}

	isSkip, err := c.isSkip()
	assert.Nil(t, err)
	assert.False(t, isSkip)
}

func TestCodeChangedRunRule(t *testing.T) {
	c := compare{
		changed: []string{
			".drone.star",
			"cmd/foo.go",
			"pkg/internal/bar.go",
			"docs/config.hugo",
		},
		skipChanged: []string{},
		runChanged: []string{
			`.*\.go$`,
		},
	}

	isSkip, err := c.isSkip()
	assert.Nil(t, err)
	assert.False(t, isSkip)
}

func TestDocChangedSkipRule(t *testing.T) {
	c := compare{
		changed: []string{
			"docs/index.md",
		},
		skipChanged: []string{
			`^docs/.*`,
		},
		runChanged: []string{},
	}
	isSkip, err := c.isSkip()
	assert.Nil(t, err)
	assert.True(t, isSkip)
}

func TestDocChangedRunRule(t *testing.T) {
	c := compare{
		changed: []string{
			"docs/index.md",
			"docs/config.hugo",
		},
		skipChanged: []string{},
		runChanged: []string{
			`.*\.go$`,
		},
	}

	isSkip, err := c.isSkip()
	assert.Nil(t, err)
	assert.True(t, isSkip)
}
