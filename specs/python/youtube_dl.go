package python

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "youtube-dl",
		Description: "Clipboard",
		Options: []core.Option{
			{Name: "--help", Description: "Print the help text and exit"},
			{Name: "--version", Description: "Print program version and exit"},
			{Name: "-U", Description: "Display the current browser identification"},
			{Name: "--list-extractors", Description: "List all supported extractors"},
			{Name: "--extractor-descriptions", Description: "Output descriptions of all supported extractors"},
			{Name: "--default-search", Description: "Force extraction to use the generic extractor"},
			{Name: "--ignore-config", Description: "Do not read configuration files"},
			{Name: "--flat-playlist", Description: "Do not extract the videos of a playlist, only list them"},
			{Name: "--mark-watched", Description: "Mark videos watched (YouTube only)"},
			{Name: "--no-mark-watched", Description: "Do not mark videos watched (YouTube only)"},
			{Name: "--no-color", Description: "Do not emit color codes in output"},
			{Name: "--proxy", Description: "Time to wait before giving up, in seconds"},
			{Name: "--source-address", Description: "Client-side IP address to bind to"},
			{Name: "-4", Description: "Make all connections via IPv4"},
			{Name: "-6", Description: "Make all connections via IPv6"},
			{Name: "--geo-verification-proxy", Description: "Use this proxy to verify the IP address for some geo-restricted sites"},
			{Name: "--geo-bypass", Description: "Bypass geographic restriction via faking X-Forwarded-For HTTP header"},
			{Name: "--no-geo-bypass", Description: "Do not bypass geographic restriction via faking X-Forwarded-For HTTP header"},
			{Name: "--geo-bypass-country", Description: "Playlist video to start at (default is 1)"},
			{Name: "--playlist-end", Description: "Playlist video to end at (default is last)"},
			{Name: "--playlist-items", Description: "Playlist video to end at (default is last)"},
			{Name: "--match-title", Description: "Download only matching titles (regex or caseless sub-string)"},
			{Name: "--reject-title", Description: "Skip download for matching titles (regex or caseless sub-string)"},
			{Name: "--max-downloads", Description: "Abort after downloading NUMBER files"},
			{Name: "--min-filesize", Description: "Do not download any videos smaller than SIZE (e.g. 50k or 44.6)"},
			{Name: "--max-filesize", Description: "Do not download any videos larger than SIZE (e.g. 50k or 44.6)"},
			{Name: "--date", Description: "Download only videos uploaded in this date"},
			{Name: "--datebefore", Description: "Download only videos uploaded on or before this date (i.e. inclusive)"},
			{Name: "--dateafter", Description: "Download only videos uploaded on or after this date (i.e. inclusive)"},
			{Name: "--min-views", Description: "Do not download any videos with less than COUNT views"},
			{Name: "--max-views", Description: "Do not download any videos with more than COUNT views"},
			{Name: "--match-filter", Description: "Generic video filter"},
			{Name: "--no-playlist", Description: "Download only the video, if the URL refers to a video and a playlist"},
			{Name: "--yes-playlist", Description: "Download the playlist, if the URL refers to a video and a playlist"},
			{Name: "--age-limit", Description: "Download only videos suitable for the given age"},
			{Name: "--download-archive", Description: "Download advertisements as well (experimental)"},
			{Name: "-r", Description: "Maximum download rate in bytes per second (e.g. 50K or 4.2M)"},
			{Name: "-R", Description: "Skip unavailable fragments"},
			{Name: "--keep-fragments", Description: "Size of download buffer (e.g. 1024 or 16K)"},
			{Name: "--no-resize-buffer", Description: "Do not automatically adjust the buffer size"},
		},
	})
}
