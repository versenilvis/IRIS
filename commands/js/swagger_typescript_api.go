package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "swagger-typescript-api",
		Description: "Generate api via swagger scheme",
		Options: []spec.Option{
			{Name: "--output", Description: "Output path of generated templates"},
			{Name: "--modular", Description: "Http client type"},
			{Name: "--clean-output", Description: "Clean output folder before generate template"},
			{Name: "--rewrite", Description: "Rewrite content in existing templates"},
			{Name: "--version", Description: "Output the current version"},
			{Name: "--path", Description: "Path/url to swagger scheme"},
			{Name: "--default-as-success", Description: "Generate additional information about request responses"},
			{Name: "--union-enums", Description: "Generate readonly properties"},
			{Name: "--route-types", Description: "Generate type definitions for API routes"},
			{Name: "--no-client", Description: "Do not generate an API class"},
			{Name: "--enum-names-as-values", Description: "Use values in 'x-enumNames' as enum values (not only as keys)"},
			{Name: "--extract-request-params", Description: "Extract request params to data contract"},
			{Name: "--extract-request-body", Description: "Extract request body type to data contract"},
			{Name: "--extract-response-body", Description: "Extract response body type to data contract"},
			{Name: "--extract-response-error", Description: "Extract response error type to data contract"},
			{Name: "--js", Description: "Generate js api module with declaration file"},
			{Name: "--module-name-index", Description: "Determines which path index should be used for routes separation"},
			{Name: "--module-name-first-tag", Description: "Splits routes based on the first tag"},
			{Name: "--disableStrictSSL", Description: "Disabled strict SSL"},
			{Name: "--disableProxy", Description: "Disabled proxy"},
			{Name: "--axios", Description: "Generate axios http client"},
			{Name: "--unwrap-response-data", Description: "Unwrap the data item from the response"},
			{Name: "--disable-throw-on-error", Description: "Do not throw an error when response.ok is not true (default: false)"},
			{Name: "--single-http-client", Description: "Ability to send HttpClient instance to Api constructor (default: false)"},
			{Name: "--silent", Description: "Output only errors to console"},
			{Name: "--default-response", Description: "Default type for empty response schema"},
			{Name: "--type-prefix", Description: "Data contract name prefix"},
			{Name: "--type-suffix", Description: "Data contract name suffix"},
			{Name: "--patch", Description: "Fix up small errors in the swagger source definition (default: false)"},
			{Name: "--debug", Description: "Additional information about processes inside this tool (default: false)"},
			{Name: "--another-array-type", Description: "Generate array types as Array<Type> (by default Type[]) (default: false)"},
			{Name: "--sort-types", Description: "Sort fields and types (default: false)"},
			{Name: "--extract-enums", Description: "Show help for swagger-typescript-api"},
		},
	})
}
