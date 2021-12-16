package end2end

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

var (
	upgrade_plan = []UpgradeStep{
		{
			OldVersion: "v0.6.1",
			NewVersion: "v0.6.2",
			Command: []string{
				"index", "rebuild",
			},
		},
	}
)

type UpgradeStep struct {
	OldVersion string
	NewVersion string
	Command    []string
}

func Upgrade(
	ctx context.Context,
	test_os, arch, datastore string,
	old_version, new_version string) error {

	for _, step := range upgrade_plan {
		// Check if the step applies - and if so run it.
		if step.OldVersion == getBareRelease(old_version) &&
			step.NewVersion == getBareRelease(new_version) {
			new_binary, err := GetRelease(test_os, arch, new_version, datastore)
			if err != nil {
				return err
			}

			config_file := filepath.Join(datastore, "server.config.yaml")

			argv := append(
				[]string{"-v", "--config", config_file},
				step.Command...)
			cmd := exec.CommandContext(ctx, new_binary, argv...)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error %v\n", string(out))
				return err
			}
			fmt.Printf("Output: %v\n", string(out))
		}
	}
	return nil
}
