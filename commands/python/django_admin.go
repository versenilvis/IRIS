package python

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "django-admin",
		Description: "Show this help message and exit",
		Subcommands: []spec.Subcommand{
			{Name: "help", Description: "Usage and help information for django-admin"},
		},
		Options: []spec.Option{
			{Name: "-h", Description: "Show this help message and exit"},
			{Name: "--version", Description: "Show program's version number and exit"},
			{Name: "-v", Description: "Raise on CommandError exceptions"},
			{Name: "--no-color", Description: "Don't colorize the command output"},
			{Name: "--force-color", Description: "Force colorization of the command output"},
			{Name: "--database", Description: "Change a user's password for django.contrib.auth"},
			{Name: "--username", Description: "Specifies the login for the superuser"},
			{Name: "--noinput", Description: "Tells Django to NOT prompt the user for input of any kind"},
			{Name: "--skip-checks", Description: "Skip system checks"},
			{Name: "--tag", Description: "Run only checks labeled with given tag"},
			{Name: "--list-tags", Description: "List available tags"},
			{Name: "--deploy", Description: "Check deployment settings"},
			{Name: "--fail-level", Description: "Run database related checks against these aliases"},
			{Name: "--locale", Description: "Locales to exclude. Default is none. Can be used multiple times"},
			{Name: "--use-fuzzy", Description: "Use fuzzy translations"},
			{Name: "--ignore", Description: "Compiles .po files to .mo files for use with builtin gettext support"},
			{Name: "--all", Description: "Specifies the output serialization format for fixtures"},
			{Name: "--indent", Description: "Specifies the indent level to use when pretty-printing output"},
			{Name: "--natural-primary", Description: "Use natural primary keys if they are available"},
			{Name: "-a", Description: "Specifies file to which the output is written"},
			{Name: "--include-views", Description: "Also output models for database views"},
			{Name: "--ignorenonexistent", Description: "An app_label or app_label.ModelName to exclude. Can be used multiple times"},
			{Name: "--format", Description: "Format of serialized data when reading from stdin"},
			{Name: "--domain", Description: "Updates the message files for all existing locales"},
			{Name: "--extension", Description: "Don't ignore the common glob-style patterns 'CVS', '.*', '*~' and '*.pyc'"},
			{Name: "--no-wrap", Description: "Don't break long message lines into several lines"},
			{Name: "--no-location", Description: "Don't write '#: filename:line' lines"},
			{Name: "--add-location", Description: "The lines include both file name and line number"},
			{Name: "--no-obsolete", Description: "Remove obsolete message strings"},
			{Name: "--keep-pot", Description: "Keep .pot file after making messages. Useful when debugging"},
			{Name: "--dry-run", Description: "Just show what migrations would be made; don't actually write them"},
			{Name: "--merge", Description: "Enable fixing of migration conflicts"},
			{Name: "--empty", Description: "Create an empty migration"},
			{Name: "-n", Description: "Use this name for migration file(s)"},
			{Name: "--no-header", Description: "Do not add header comments to new migration file(s)"},
			{Name: "--check", Description: "Exit with a non-zero status if model changes are missing migrations"},
			{Name: "--fake-initial", Description: "Shows a list of the migration actions that will be performed"},
			{Name: "--run-syncdb", Description: "Creates tables for apps without migrations"},
			{Name: "--managers", Description: "Send a test email to the addresses specified in settings.MANAGERS"},
			{Name: "--admins", Description: "Send a test email to the addresses specified in settings.ADMINS"},
		},
	})
}
