package main

import (
	"github.com/owncloud-ci/drone-skip-pipeline/plugin"
	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "git-path",
			Value:       ".",
			EnvVars:     []string{"PLUGIN_GIT_PATH"},
			Usage:       "path to the git repository",
			Destination: &settings.GitPath,
		},
		&cli.StringSliceFlag{
			Name:        "pattern-run-changed",
			EnvVars:     []string{"PLUGIN_RUN_ON_CHANGED_PATTERN"},
			Usage:       "pattern to run on if changed",
			Destination: &settings.RunChangedPattern,
		},
		&cli.StringSliceFlag{
			Name:        "pattern-skip-changed",
			EnvVars:     []string{"PLUGIN_SKIP_ON_CHANGED_PATTERN"},
			Usage:       "pattern to skip if changed",
			Destination: &settings.SkipChangedPattern,
		},
	}
}
