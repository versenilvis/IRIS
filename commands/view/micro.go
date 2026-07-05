package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "micro",
		Description: "True/false",
		Subcommands: []spec.Subcommand{
			{Name: "remove", Description: "Remove plugin(s)"},
			{Name: "update", Description: "Update plugin(s) (if no argument is given, updates all plugins)"},
			{Name: "search", Description: "Search for a plugin"},
			{Name: "list", Description: "List installed plugins"},
			{Name: "available", Description: "List available plugins"},
		},
		Options: []spec.Option{
			{Name: "--plugin", Description: "Manage plugins"},
			{Name: "--clean", Description: "Cleans the configuration directory"},
			{Name: "--config-dir", Description: "Specify a custom location for the configuration directory"},
			{Name: "--options", Description: "Show all option help"},
			{Name: "--debug", Description: "Enable debug mode"},
			{Name: "--version", Description: "Show the version number and information"},
			{Name: "--autoindent", Description: "When creating a new line, use the same indentation as the previous line"},
			{Name: "--autosave", Description: "Seconds"},
			{Name: "--autosu", Description: "Automatically attempt to use super user privileges to save without asking"},
			{Name: "--backup", Description: "Automatically keep backups of all open buffers"},
			{Name: "--backupdir", Description: "The directory to place backups in"},
			{Name: "--basename", Description: "Apecifies how the system clipboard should be accessed"},
			{Name: "--colorcolumn", Description: "Display a color at the specified column if not set to 0"},
			{Name: "--colorscheme", Description: "Loads the colorscheme stored in $(configDir)/colorschemes/<scheme>.micro"},
			{Name: "--cursorline", Description: "Highlight the line that the cursor is on in a different color"},
			{Name: "--diffgutter", Description: "Display diff indicators before lines"},
			{Name: "--divchars", Description: "Reverse colors specified by the colorscheme"},
			{Name: "--encoding", Description: "The encoding to open and save files with"},
			{Name: "--eofnewline", Description: "Automatically add a newline to the end of the file if one does not exist"},
			{Name: "--fastdirty", Description: "Use 'fast dirty' algorithm to determine if a buffer is modified or not"},
			{Name: "--fileformat", Description: "Type of line endings to be used for the file"},
			{Name: "--filetype", Description: "File type for the current buffer"},
			{Name: "--ignorecase", Description: "Perform case-insensitive searches"},
			{Name: "--incsearch", Description: "Sets the indentation character"},
			{Name: "--infobar", Description: "Enables the line at the bottom of the editor where messages are printed"},
			{Name: "--keepautoindent", Description: "Display the nano-style key menu at the bottom of the screen"},
			{Name: "--matchbrace", Description: "Enable mouse support"},
			{Name: "--paste", Description: "Cause backups to be permanently saved"},
			{Name: "--pluginchannels", Description: "List of URLs pointing to plugin channels for downloading and installing plugins"},
			{Name: "--pluginrepos", Description: "A list of links to plugin repositories"},
			{Name: "--readonly", Description: "Disallow edits to the buffer"},
			{Name: "--relativeruler", Description: "Make line numbers display relatively"},
			{Name: "--rmtrailingws", Description: "Automatically trim trailing whitespaces at ends of lines"},
			{Name: "--ruler", Description: "Display line numbers"},
			{Name: "--savecursor", Description: "Remember command history between closing and re-opening micro"},
			{Name: "--saveundo", Description: "Remember undo state between closing and re-opening a file"},
			{Name: "--scrollbar", Description: "Display a scroll bar"},
			{Name: "--scrollmargin", Description: "Number of lines to scroll for one scroll event"},
			{Name: "--smartpaste", Description: "Add leading whitespace when pasting multiple lines"},
			{Name: "--softwrap", Description: "Wrap lines that are too long to fit on the screen"},
		},
	})
}
