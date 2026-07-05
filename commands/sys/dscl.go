package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "dscl",
		Description: "Prompt for password",
		Subcommands: []spec.Subcommand{
			{Name: "read", Description: "Prints a directory"},
			{Name: "readall", Description: "Prints all the records of a given type"},
			{Name: "readpl", Description: "Prints the contents of plist_path"},
			{Name: "readpli", Description: "Prints the contents of plist_path for the plist at value_index of the key"},
			{Name: "list", Description: "Lists the subdirectories of the given directory"},
			{Name: "search", Description: "Searches for records that match a pattern"},
			{Name: "create", Description: "Creates a new record"},
			{Name: "createpl", Description: "Creates a string, or array of strings at plist_path"},
			{Name: "append", Description: "Appends one or more values to a property in a given record"},
			{Name: "delete", Description: "Delete a directory, property, or value"},
			{Name: "deletepl", Description: "Deletes a value in a plist"},
			{Name: "deletepli", Description: "Deletes a value for the plist at value_index of the key"},
			{Name: "passwd", Description: "Changes the password of a user"},
			{Name: "new password", Description: "New password of the user"},
		},
		Options: []spec.Option{
			{Name: "-p", Description: "Prompt for password"},
			{Name: "-u", Description: "Authenticate as user"},
			{Name: "-P", Description: "Authenticate with password"},
			{Name: "-f", Description: "Targeted local node database file path"},
			{Name: "-raw", Description: "Don't strip off prefix from DirectoryService API constants"},
			{Name: "-plist", Description: "Print out record(s) or attribute(s) in XML plist format"},
			{Name: "-url", Description: "Print record attribute values in URL-style encoding"},
			{Name: "-q", Description: "Quiet - no interactive prompt"},
		},
	})
}
