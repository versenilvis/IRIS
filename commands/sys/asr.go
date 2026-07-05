package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "asr",
		Description: "Can be a disk image, /dev entry, or volume mountpoint",
		Subcommands: []spec.Subcommand{
			{Name: "help", Description: "Displays asr usage information"},
			{Name: "version", Description: "Displays asr version"},
			{Name: "restore", Description: "Restores a disk image or volume to another volume"},
			{Name: "server", Description: "Multicasts source over the network"},
			{Name: "source", Description: "UDIF disk image local/remote path"},
			{Name: "interface", Description: "The network interface to be used for multicasting"},
			{Name: "configuration", Description: "Configuration file in standard property list format"},
		},
		Options: []spec.Option{
			{Name: "--source", Description: "Can be a disk image, /dev entry, or volume mountpoint"},
			{Name: "--target", Description: "Can be a /dev entry, or volume mountpoint"},
			{Name: "--file", Description: "When performing a multicast restore, --file can be specified instead of --target"},
			{Name: "--erase", Description: "Specifies the destination filesystem format"},
			{Name: "--noprompt", Description: "Suppresses the prompt which usually occurs before target"},
			{Name: "--timeout", Description: "Number of seconds that a multicast client should wait"},
			{Name: "--puppetstrings", Description: "Provide progress output that is easy for another program to parse"},
			{Name: "--noverify", Description: "Allows restores to proceed even if the source's catalog file is fragmented"},
			{Name: "--SHA256", Description: "Forces the restore to use the SHA-256 hash in the image during verification"},
			{Name: "--sourcevolumename", Description: "Forces asr to use replication for restoring APFS volumes"},
			{Name: "--useInverter", Description: "Forces asr to use the inverter for restoring APFS volumes"},
			{Name: "--toSnapshot", Description: "One of the options that control how asr uses memory"},
			{Name: "--buffersize", Description: "One of the options that control how asr uses memory"},
			{Name: "--csumbuffers", Description: "One of the options that control how asr uses memory"},
			{Name: "--csumbuffersize", Description: "One of the options that control how asr uses memory"},
			{Name: "--verbose", Description: "Enables verbose progress and error messages"},
			{Name: "--debug", Description: "Enables other progress and error messages"},
			{Name: "--interface", Description: "The network interface to be used for multicasting"},
			{Name: "--config", Description: "Server requires a configuration file to be passed"},
			{Name: "--plist", Description: "Writes its output as an XML-formatted plist"},
		},
	})
}
