package plugin

import (
	"errors"
	"os"

	"github.com/urfave/cli/v2"
)

// Settings for the Plugin.
type Settings struct {
	GitPath             string
	DisallowSkipChanged cli.StringSlice
	AllowSkipChanged    cli.StringSlice
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {

	if len(p.settings.DisallowSkipChanged.Value()) == 0 && len(p.settings.AllowSkipChanged.Value()) == 0 {
		return errors.New("you must at least set a allow skip or disallow skip pattern")
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {

	c := newCompare(
		p.settings.GitPath,
		"origin",
		p.pipeline.Commit.SHA,
		p.pipeline.Build.TargetBranch,
		p.settings.DisallowSkipChanged.Value(),
		p.settings.AllowSkipChanged.Value(),
	)

	err := c.open()
	if err != nil {
		return err
	}

	err = c.getChanged()
	if err != nil {
		return err
	}

	isSkip, err := c.isSkip()
	if err != nil {
		return err
	}

	if isSkip {
		// https://discourse.drone.io/t/how-to-exit-a-pipeline-early-without-failing/3951
		os.Exit(78)
	}

	return nil // don't report errors properly, because then we can't control the code
}
