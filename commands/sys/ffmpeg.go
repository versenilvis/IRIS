package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "ffmpeg",
		Description: "Play, record, convert, and stream audio and video",
		Options: []spec.Option{
			{Name: "-i", Description: "Input file"},
			{Name: "-L", Description: "Show license"},
			{Name: "-h", Description: "Show help"},
			{Name: "-help", Description: "Show help"},
			{Name: "--help", Description: "Show help"},
			{Name: "-version", Description: "Show version"},
			{Name: "-buildconf", Description: "Show build configuration"},
			{Name: "-formats", Description: "Show available formats"},
			{Name: "-muxers", Description: "Show available muxers"},
			{Name: "-demuxers", Description: "Show available demuxers"},
			{Name: "-devices", Description: "Show available devices"},
			{Name: "-codecs", Description: "Show available codecs"},
			{Name: "-decoders", Description: "Show available decoders"},
			{Name: "-encoders", Description: "Show available encoders"},
			{Name: "-bsfs", Description: "Show available bit stream filters"},
			{Name: "-protocols", Description: "Show available protocols"},
			{Name: "-filters", Description: "Show available filters"},
			{Name: "-pix_fmts", Description: "Show available pixel formats"},
			{Name: "-layouts", Description: "Show standard channel layouts"},
			{Name: "-sample_fmts", Description: "Show available audio sample formats"},
			{Name: "-dispositions", Description: "Show available stream dispositions"},
			{Name: "-colors", Description: "Show available color names"},
			{Name: "-sources", Description: "List sources of the input device"},
			{Name: "-sinks", Description: "List sinks of the output device"},
			{Name: "-hwaccels", Description: "Show available HW acceleration methods"},
			{Name: "-loglevel", Description: "Set logging level"},
			{Name: "-v", Description: "Set logging level"},
			{Name: "-report", Description: "Generate a report"},
			{Name: "-max_alloc", Description: "Set maximum size of a single allocated block"},
			{Name: "-y", Description: "Overwrite output files"},
			{Name: "-n", Description: "Never overwrite output files"},
			{Name: "-ignore_unknown", Description: "Ignore unknown stream types"},
			{Name: "-filter_threads", Description: "Number of non-complex filter threads"},
			{Name: "-filter_complex_threads", Description: "Number of threads for -filter_complex"},
			{Name: "-stats", Description: "Print progress report during encoding"},
			{Name: "-max_error_rate", Description: "Change audio volume (256=normal)"},
			{Name: "-cpuflags", Description: "Force specific cpu flags"},
			{Name: "-cpucount", Description: "Force specific cpu count"},
			{Name: "-hide_banner", Description: "Do not show program banner"},
			{Name: "-copy_unknown", Description: "Copy unknown stream types"},
		},
	})
}
