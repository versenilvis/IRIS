package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "neofetch",
		Description: "The most complete system information CLI tool",
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for neofetch"},
			{Name: "--disable", Description: "Disable information line"},
			{Name: "--title_fqdn", Description: "Hide/Show Fully Qualified Domain Name in title"},
			{Name: "--package_managers", Description: "Hide/Show Package Manager names"},
			{Name: "--os_arch", Description: "Hide/Show OS architecture"},
			{Name: "--speed_type", Description: "Change the type of cpu speed to display"},
			{Name: "--speed_shorthand", Description: "Whether or not to show decimals in CPU speed"},
			{Name: "--cpu_brand", Description: "Enable/Disable CPU brand in output"},
			{Name: "--cpu_cores", Description: "Whether or not to display the number of CPU cores"},
			{Name: "--cpu_speed", Description: "Hide/Show cpu speed"},
			{Name: "--cpu_temp", Description: "Hide/Show cpu temperature"},
			{Name: "--distro_shorthand", Description: "Shorten the output of distro"},
			{Name: "--kernel_shorthand", Description: "Shorten the output of kernel"},
			{Name: "--uptime_shorthand", Description: "Shorten the output of uptime"},
			{Name: "--refresh_rate", Description: "Whether to display the refresh rate of each monitor"},
			{Name: "--gpu_brand", Description: "Enable/Disable GPU brand in output"},
			{Name: "--gpu_type", Description: "Which GPU to display"},
			{Name: "--de_version", Description: "Show/Hide Desktop Environment version"},
			{Name: "--gtk_shorthand", Description: "Shorten output of gtk theme/icons"},
			{Name: "--gtk2", Description: "Enable/Disable gtk2 theme/font/icons output"},
			{Name: "--gtk3", Description: "Enable/Disable gtk3 theme/font/icons output"},
			{Name: "--shell_path", Description: "Enable/Disable showing $SHELL path"},
			{Name: "--shell_version", Description: "Enable/Disable showing $SHELL version"},
			{Name: "--disk_show", Description: "Which disks to display"},
			{Name: "--disk_subtitle", Description: "What information to append to the Disk subtitle"},
			{Name: "--disk_percent", Description: "Hide/Show disk percent"},
			{Name: "--ip_host", Description: "URL to query for public IP"},
			{Name: "--ip_timeout", Description: "Public IP timeout"},
			{Name: "--memory_percent", Description: "Display memory percentage"},
			{Name: "--memory_unit", Description: "Memory output unit"},
			{Name: "--colors", Description: "Changes the text colors"},
			{Name: "--underline", Description: "Enable/Disable the underline"},
			{Name: "--underline_char", Description: "Character to use when underlining title"},
			{Name: "--bold", Description: "Enable/Disable bold text"},
			{Name: "--separator", Description: "Changes the default ':' separator to the specified string"},
			{Name: "--color_blocks", Description: "Enable/Disable the color blocks"},
			{Name: "--col_offset", Description: "Left-padding of color blocks"},
			{Name: "--block_width", Description: "Width of color blocks in spaces"},
			{Name: "--block_height", Description: "Height of color blocks in lines"},
			{Name: "--block_range", Description: "Range of colors to print as blocks"},
		},
	})
}
