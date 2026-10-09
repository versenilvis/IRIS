package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/internal/updater"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Iris to the latest release",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("checking for updates (current: %s)...\n", Version)

		latest, err := updater.FetchLatestVersion()
		if err != nil {
			fmt.Printf("\033[31m[IRIS] could not reach update server: %v\033[0m\n", err)
			return
		}

		if Version != "dev" && Version != "" && !updater.IsNewer(Version, latest) {
			fmt.Printf("\033[32m[IRIS] already up to date (%s)\033[0m\n", Version)
			state := config.LoadState()
			state.Updater.SeenVersion = ""
			_ = config.SaveState(state)
			return
		}

		fmt.Printf("\033[36m[IRIS] updating %s → %s\033[0m\n", Version, latest)

		runningPrefix := ""
		if config.Get().Updater.Channel == "nightly" {
			runningPrefix = fmt.Sprintf("IRIS_RELEASE_TAG=%s ", latest)
		}
		fmt.Printf("running: %scurl -sSL %s | sh\n\n", runningPrefix, updater.ResolveInstallScriptURL())

		if _, err := updater.PerformUpdate(latest, true); err != nil {
			fmt.Printf("\n\033[31m[IRIS] update failed: %v\033[0m\n", err)
			return
		}

		state := config.LoadState()
		state.Updater.SeenVersion = ""
		_ = config.SaveState(state)

		fmt.Printf("\n\033[32m[IRIS] restart your terminal to use the new version\033[0m\n")
	},
}
