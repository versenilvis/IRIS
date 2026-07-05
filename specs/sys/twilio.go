package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "twilio",
		Description: "Level of logging messages",
		Subcommands: []core.Subcommand{
			{Name: "list", Description: "List Twilio CLI configurations"},
			{Name: "set", Description: "Update Twilio CLI configurations"},
		},
		Options: []core.Option{
			{Name: "-l", Description: "Level of logging messages"},
			{Name: "-o", Description: "Format of command output"},
			{Name: "--silent", Description: "Manage Twilio CLI configurations"},
			{Name: "-e", Description: "Sets an Edge configuration"},
			{Name: "--require-profile-input", Description: "Show a list of log events for the account"},
			{Name: "--attachment", Description: "Path for the file that you want to attach"},
			{Name: "--from", Description: "Email address of the sender"},
			{Name: "--no-attachment", Description: "Do not include or prompt for an attachment"},
			{Name: "--subject", Description: "The subject line for an email"},
			{Name: "--text", Description: "Text to send within the email body"},
			{Name: "--to", Description: "Email address of recipient (comma-separated)"},
			{Name: "--country-code", Description: "The ISO-3166-1 country code of the country from which to read phone numbers"},
			{Name: "--account-sid", Description: "The SID of the Account requesting the AvailablePhoneNumber resources"},
			{Name: "--address-sid", Description: "Whether the phone numbers can receive faxes. Can be: `true` or `false`"},
			{Name: "--no-fax-enabled", Description: "Whether the phone numbers can receive faxes. Can be: `true` or `false`"},
			{Name: "--friendly-name", Description: "The maximum number of resources to return. Use '--no-limit' to disable"},
			{Name: "--mms-enabled", Description: "Whether the phone numbers can receive MMS messages. Can be: `true` or `false`"},
			{Name: "--no-mms-enabled", Description: "Whether the phone numbers can receive MMS messages. Can be: `true` or `false`"},
			{Name: "--near-lat-long", Description: "Skip including of headers while listing the data"},
			{Name: "--page-size", Description: "Whether the phone numbers can receive text messages. Can be: `true` or `false`"},
			{Name: "--no-sms-enabled", Description: "Whether the phone numbers can receive text messages. Can be: `true` or `false`"},
			{Name: "--sms-fallback-method", Description: "Whether the phone numbers can receive calls. Can be: `true` or `false`"},
			{Name: "--no-voice-enabled", Description: "Whether the phone numbers can receive calls. Can be: `true` or `false`"},
			{Name: "--voice-fallback-method", Description: "Manage Twilio phone numbers"},
			{Name: "--no-header", Description: "Skip including of headers while listing the data"},
			{Name: "--properties", Description: "Update the properties of a Twilio phone number"},
			{Name: "--sms-fallback-url", Description: "The HTTP method Twilio will use when making requests to the SmsUrl"},
			{Name: "--sms-url", Description: "The HTTP method Twilio will use when requesting the VoiceFallbackUrl"},
			{Name: "--voice-fallback-url", Description: "The HTTP method Twilio will use when making requests to the VoiceUrl"},
			{Name: "--voice-url", Description: "The URL that Twilio should request when somebody dials the phone number"},
			{Name: "-v", Description: "Installs a plugin into the CLI"},
			{Name: "-f", Description: "`yarn install` with force flag"},
			{Name: "-p", Description: "Shorthand identifier for your profile"},
			{Name: "--auth-token", Description: "Your Twilio Auth Token for your Twilio Account or Subaccount"},
			{Name: "--region", Description: "Twilio region to use"},
			{Name: "--from-local", Description: "Interactively choose an already installed version"},
			{Name: "-r", Description: "Refresh cache (ignores displaying instructions)"},
			{Name: "-h", Description: "Show a help message"},
			{Name: "--core", Description: "List core plugins"},
		},
	})
}
