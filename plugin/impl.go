package plugin

import (
	"errors"
	"os"

	"github.com/urfave/cli/v2"
)

// Settings for the Plugin.
type Settings struct {
	GitPath            string
	SkipChangedPattern cli.StringSlice
	RunChangedPattern  cli.StringSlice
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {

	if len(p.settings.SkipChangedPattern.Value()) == 0 && len(p.settings.RunChangedPattern.Value()) == 0 {
		return errors.New("you must at least set a skip or run pattern")
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {

	c := newCompare(
		p.settings.GitPath,
		p.pipeline.Build.TargetBranch,
		p.settings.RunChangedPattern.Value(),
		p.settings.SkipChangedPattern.Value(),
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
