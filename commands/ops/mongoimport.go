package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "mongoimport",
		Description: "Import data from a JSON, CSV, or TSV file into a MongoDB instance",
		Options: []core.Option{
			{Name: "--help", Description: "Returns information on the options and use of mongoimport"},
			{Name: "--verbose", Description: "Runs mongoimport in a quiet mode that attempts to limit the amount of output"},
			{Name: "--version", Description: "Returns the mongoimport release number"},
			{Name: "--config", Description: "Specifies the resolvable URI connection string of the MongoDB deployment"},
			{Name: "--host", Description: "Specifies the resolvable hostname of the MongoDB deployment"},
			{Name: "--port", Description: "Default port"},
			{Name: "--ssl", Description: "Enables connection to a mongod or mongos that has TLS/SSL support enabled"},
			{Name: "--sslCAFile", Description: "Specifies the .pem file that contains both the TLS/SSL certificate and key"},
			{Name: "--sslPEMKeyPassword", Description: "Password"},
			{Name: "--sslCRLFile", Description: "Specifies the .pem file that contains the Certificate Revocation List"},
			{Name: "--sslAllowInvalidCertificates", Description: "Disables the validation of the hostnames in TLS/SSL certificates"},
			{Name: "--username", Description: "Specifies the session token for MONGODB-AWS authentication mechanism"},
			{Name: "--authenticationDatabase", Description: "Database name"},
			{Name: "--authenticationMechanism", Description: "Default"},
			{Name: "--gssapiServiceName", Description: "Specifies the name of the database on which to run the mongoimport"},
			{Name: "--collection", Description: "Specifies the name of the collection to import"},
			{Name: "--fields", Description: "Comma separated list of fields"},
			{Name: "--fieldFile", Description: "Ignores empty fields in CSV and TSV exports"},
			{Name: "--type", Description: "Specifies the file type to import"},
			{Name: "--file", Description: "Specifies the location and name of a file containing the data to import"},
			{Name: "--drop", Description: "Treats the first line of the input file as a header line"},
			{Name: "--useArrayIndexFields", Description: "Insert the documents in the import file"},
			{Name: "--upsertFields", Description: "Specifies a list of fields for the query portion of the import process"},
			{Name: "--stopOnError", Description: "Indicates that the import data is in Extended JSON v1 format"},
			{Name: "--maintainInsertionOrder", Description: "Specifies the number of insertion workers to run concurrently"},
			{Name: "--writeConcern", Description: "Specifies the write concern for each write operation that mongoimport performs"},
			{Name: "--bypassDocumentValidation", Description: "Enables mongoimport to bypass document validation during the operation"},
		},
	})
}
