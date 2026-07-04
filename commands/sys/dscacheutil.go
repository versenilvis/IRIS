package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "dscacheutil",
		Description: "Utility for managing the Directory Service cache",
		Subcommands: []core.Subcommand{
			{Name: "category", Description: "Category to query"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "List the options for calling dscacheutil"},
			{Name: "-q", Description: "Query the Directory Service cache"},
			{Name: "-a", Description: "Attribute to query"},
			{Name: "-cachedump", Description: "Get an overview of the cache by default"},
			{Name: "-buckets", Description: "Get an overview of the cache by default"},
			{Name: "-entries", Description: "Dump detailed information about cache entries"},
			{Name: "-configuration", Description: "Flush the entire cache"},
			{Name: "-statistics", Description: "Get statistics from the cache, including an overview an detailed call statistics"},
		},
	})
}
