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
			Name:        "disallow-skip-changed",
			EnvVars:     []string{"PLUGIN_DISALLOW_SKIP_CHANGED"},
			Usage:       "files never allowed to be skipped if changed",
			Destination: &settings.DisallowSkipChanged,
		},
		&cli.StringSliceFlag{
			Name:        "allow-skip-changed",
			EnvVars:     []string{"PLUGIN_ALLOW_SKIP_CHANGED"},
			Usage:       "files allowed to be skipped, even if changed",
			Destination: &settings.AllowSkipChanged,
		},
	}
}
