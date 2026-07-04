package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "expo",
		Description: "Tools for creating, running, and deploying Universal Expo and React Native apps",
		Options: []core.Option{
			{Name: "-h", Description: "Output usage information"},
			{Name: "-V", Description: "Output the version number"},
			{Name: "-u", Description: "Username"},
			{Name: "-p", Description: "Password"},
			{Name: "--otp", Description: "One-time password from your 2FA device"},
			{Name: "-c", Description: "Clear all credentials stored on Expo servers"},
			{Name: "--clear-dist-cert", Description: "Remove Distribution Certificate stored on Expo servers"},
			{Name: "--clear-push-key", Description: "Remove Push Notifications Key stored on Expo servers"},
			{Name: "--clear-push-cert", Description: "Remove Provisioning Profile stored on Expo servers"},
			{Name: "-r", Description: "Type of build: [archive|simulator]"},
			{Name: "--release-channel", Description: "Pull from specified release channel"},
			{Name: "--no-publish", Description: "Disable automatic publishing before building"},
			{Name: "--no-wait", Description: "Exit immediately after scheduling build"},
			{Name: "--team-id", Description: "Apple Team ID"},
			{Name: "--dist-p12-path", Description: "Push Key ID (ex: 123AB4C56D)"},
			{Name: "--push-p8-path", Description: "Path to your Push Key .p8 file"},
			{Name: "--provisioning-profile-path", Description: "Path to your Provisioning Profile"},
			{Name: "--public-url", Description: "The URL of an externally hosted manifest (for self-hosted apps)"},
			{Name: "--skip-credentials-check", Description: "Skip checking credentials"},
			{Name: "--skip-workflow-check", Description: "Skip warning about build service bare workflow limitations"},
			{Name: "--config", Description: "Deprecated: Use app.config.js to switch config files instead"},
			{Name: "--keystore-path", Description: "Path to your Keystore: *.jks"},
			{Name: "--keystore-alias", Description: "Keystore Alias"},
			{Name: "--generate-keystore", Description: "[deprecated] Generate Keystore if one does not exist"},
			{Name: "-t", Description: "Type of build: [app-bundle|apk]"},
			{Name: "--no-pwa", Description: "Turns dev flag on before bundling"},
			{Name: "--apple-id", Description: "Deprecated: Use app.config.js to switch config files instead"},
			{Name: "--latest", Description: "Install the latest version of Expo Go, ignoring the current project version"},
			{Name: "-d", Description: "Device name to install the client on"},
			{Name: "-f", Description: "Allows replacing existing files"},
			{Name: "--offline", Description: "Allows this command to run while offline"},
			{Name: "--no-install", Description: "Skip installing npm packages and CocoaPods"},
			{Name: "--npm", Description: "Use npm to install dependencies. (default when Yarn is not installed)"},
			{Name: "--clean", Description: "Delete the native folders and regenerate them before applying changes"},
			{Name: "--template", Description: "Platforms to sync: ios, android, all. Default: all"},
			{Name: "--skip-dependency-update", Description: "Preserves versions of listed packages in package.json (comma separated list)"},
			{Name: "--dest", Description: "Destination directory for assets"},
			{Name: "--platform", Description: "Detached project platform"},
			{Name: "--skipXcodeConfig", Description: "[iOS only] if true, do not configure Xcode project"},
			{Name: "--output-dir", Description: "The directory to export the static files to"},
		},
	})
}
