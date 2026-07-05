package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "screen",
		Description: "Screen manager with VT100/ANSI terminal emulation",
		Subcommands: []spec.Subcommand{
			{Name: "name", Description: "Name of the screen session"},
		},
		Options: []spec.Option{
			{Name: "-d", Description: "Does not start screen, but detaches the elsewhere running screen session"},
			{Name: "-r", Description: "Reattach a session and if necessary detach it first"},
			{Name: "-R", Description: "Reattach a session and if necessary detach or even create it first"},
			{Name: "-RR", Description: "Does not start screen, but detaches the elsewhere running screen session"},
			{Name: "-dmS", Description: "Start as daemon: Screen session in detached mode"},
			{Name: "-a", Description: "Force all capabilities into each window's termcap"},
			{Name: "-A", Description: "Adapt all windows to the new display width & height"},
			{Name: "-c", Description: "Read configuration file instead of '.screenrc'"},
			{Name: "-e", Description: "Change command characters"},
			{Name: "-f", Description: "Flow control on"},
			{Name: "-fn", Description: "Flow control off"},
			{Name: "-fa", Description: "Flow control automatic"},
			{Name: "-h", Description: "Set the size of the scrollback history buffer"},
			{Name: "-i", Description: "Interrupt output sooner when flow control is on"},
			{Name: "-list", Description: "Do nothing, just list our SockDir"},
			{Name: "-L", Description: "Turn on output logging"},
			{Name: "-m", Description: "Ignore $STY variable, do create a new screen session"},
			{Name: "-O", Description: "Choose optimal output rather than exact vt100 emulation"},
			{Name: "-p", Description: "Preselect the named window if it exists"},
			{Name: "-q", Description: "Quiet startup. Exits with non-zero return code if unsuccessful"},
			{Name: "-s", Description: "Shell to execute rather than $SHELL"},
			{Name: "-S", Description: "Name this session <pid>.sockname instead of <pid>.<tty>.<host>"},
			{Name: "-t", Description: "Set title. (window's name)"},
			{Name: "-T", Description: "Use term as $TERM for windows, rather than 'screen'"},
			{Name: "-U", Description: "Tell screen to use UTF-8 encoding"},
			{Name: "-v", Description: "Print 'Screen version 4.00.03 (FAU) 23-Oct-06'"},
			{Name: "-wipe", Description: "Do nothing, just clean up SockDir"},
			{Name: "-x", Description: "Attach to a not detached screen. (Multi display mode)"},
			{Name: "-X", Description: "Execute <cmd> as a screen command in the specified session"},
		},
	})
}
