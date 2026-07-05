package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "spotify",
		Description: "CLI to use Spotify from the terminal",
		Subcommands: []spec.Subcommand{
			{Name: "play", Description: "Resume playback where Spotify last left off"},
			{Name: "song name", Description: "The name of the song to start playing"},
			{Name: "album", Description: "Find an album by name and play it"},
			{Name: "album name", Description: "The album name you want to play"},
			{Name: "artist", Description: "Find an artist by name and play it"},
			{Name: "artist name", Description: "The artist name you want to play"},
			{Name: "list", Description: "Find a playlist by name and plays it"},
			{Name: "playlist name", Description: "The playlist name you want to play"},
			{Name: "uri", Description: "Play songs from specific uri"},
			{Name: "next", Description: "Skip to the next song in a playlist"},
			{Name: "prev", Description: "Return to the previous song in a playlist"},
			{Name: "replay", Description: "Replay the current track from the beginning"},
			{Name: "pos", Description: "Jump to a time (in secs) in the current song"},
			{Name: "time", Description: "Exact time in secs to jump in"},
			{Name: "pause", Description: "Pause (or resume) Spotify playback"},
			{Name: "stop", Description: "Stop playback"},
			{Name: "quit", Description: "Stop playback and quit Spotify"},
			{Name: "vol", Description: "Show the current Spotify volume"},
			{Name: "amount", Description: "Set the volume to an amount between 0 and 100"},
			{Name: "up", Description: "Increase the volume by 10%"},
			{Name: "down", Description: "Decrease the volume by 10%"},
			{Name: "status", Description: "Show the current player status"},
			{Name: "track", Description: "Show the currently playing track"},
			{Name: "share", Description: "Display the current song's Spotify URL and URI"},
			{Name: "url", Description: "Display the current song's Spotify URL and copies it to the clipboard"},
			{Name: "shuffle", Description: "Toggle shuffle playback mode"},
			{Name: "repeat", Description: "Toggle repeat playback mode"},
		},
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for spotify"},
		},
	})
}
