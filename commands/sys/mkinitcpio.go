package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "mkinitcpio",
		Description: "Create an initial ramdisk environment",
		Options: []core.Option{
			{Name: "--help", Description: "Show help for mkinitcpio"},
			{Name: "--version", Description: "Display version information"},
			{Name: "--addhooks", Description: "Add additional hooks to the image"},
			{Name: "--config", Description: "Use config file to generate the ramdisk"},
			{Name: "--generatedir", Description: "Set directory as the location where the initramfs is built"},
			{Name: "--hookdir", Description: "Generate a CPIO image as filename"},
			{Name: "--hookhelp", Description: "Output help for hookname"},
			{Name: "--kernel", Description: "Use kernelversion, instead of the current running kernel"},
			{Name: "--listhooks", Description: "List all available hooks"},
			{Name: "--automods", Description: "Display modules found via autodetection"},
			{Name: "--nocolor", Description: "Disable color output"},
			{Name: "--uefi", Description: "Generate a UEFI executable as filename"},
			{Name: "--allpresets", Description: "Process all presets contained in /etc/mkinitcpio.d"},
			{Name: "--preset", Description: "Build initramfs image(s) according to specified preset"},
			{Name: "--moduleroot", Description: "Specifies the root directory to find modules in"},
			{Name: "--skiphooks", Description: "Skip hooks when generating the image"},
			{Name: "--save", Description: "Saves the build directory for the initial ramdisk"},
			{Name: "--builddir", Description: "Use tmpdir as the temporary build directory instead of /tmp"},
			{Name: "--verbose", Description: "Verbose output"},
			{Name: "--compress", Description: "Override the compression method with the compress program"},
			{Name: "--cmdline", Description: "Use kernel command line with UEFI executable"},
			{Name: "--splash", Description: "UEFI executables can show a bitmap file on boot"},
			{Name: "--uefistub", Description: "UEFI stub image used for UEFI executable generation"},
			{Name: "--kernelimage", Description: "Include a kernel image for the UEFI executable"},
			{Name: "--microcode", Description: "Include microcode into the UEFI executable"},
			{Name: "--osrelease", Description: "Include a os-release file for the UEFI executable"},
		},
	})
}
