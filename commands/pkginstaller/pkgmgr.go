package pkginstaller

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/versenilvis/iris/commands/core"
)

func installedPackageGenerator(pm string) core.GeneratorFunc {
	return func(tokens []string, _ string, _ string) []core.Suggestion {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		var cmd *exec.Cmd
		switch pm {
		case "pacman", "yay", "paru":
			cmd = exec.CommandContext(ctx, pm, "-Qq")
		case "apt", "apt-get":
			cmd = exec.CommandContext(ctx, "dpkg-query", "-W", "-f=${Package}\n")
		case "dnf", "yum":
			cmd = exec.CommandContext(ctx, "rpm", "-qa", "--qf", "%{NAME}\n")
		case "brew":
			cmd = exec.CommandContext(ctx, "brew", "list")
		default:
			return nil
		}

		out, err := cmd.Output()
		if err != nil {
			return nil
		}

		var results []core.Suggestion
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			results = append(results, core.Suggestion{Cmd: line, Desc: "installed"})
		}
		return results
	}
}

func init() {
	// pacman
	core.Register(&core.Spec{
		Name:        "pacman",
		Description: "Arch package manager",
		Options: []core.Option{
			{Name: "-S", Description: "install package"},
			{Name: "-Syu", Description: "full system upgrade"},
			{Name: "-Sy", Description: "sync database"},
			{Name: "-Su", Description: "upgrade installed"},
			{Name: "-Ss", Description: "search in repos"},
			{Name: "-Si", Description: "show package info"},
			{Name: "-R", Description: "remove package"},
			{Name: "-Rs", Description: "remove + unused deps"},
			{Name: "-Rns", Description: "remove + deps + configs"},
			{Name: "-Q", Description: "query installed"},
			{Name: "-Qq", Description: "query names only"},
			{Name: "-Qs", Description: "search installed"},
			{Name: "-Qi", Description: "installed package info"},
			{Name: "-Qo", Description: "which package owns file"},
			{Name: "-Ql", Description: "list package files"},
			{Name: "-U", Description: "install local package"},
			{Name: "-D", Description: "set install reason"},
			{Name: "--noconfirm", Description: "skip confirmations"},
			{Name: "--needed", Description: "skip up-to-date"},
			{Name: "--asdeps", Description: "mark as dependency"},
			{Name: "--asexplicit", Description: "mark as explicit"},
		},
		Generator: installedPackageGenerator("pacman"),
	})

	// yay
	core.Register(&core.Spec{
		Name:        "yay",
		Description: "AUR helper (pacman wrapper)",
		Options: []core.Option{
			{Name: "-S", Description: "install package (AUR + repos)"},
			{Name: "-Syu", Description: "full system upgrade + AUR"},
			{Name: "-Ss", Description: "search AUR + repos"},
			{Name: "-Si", Description: "show package info"},
			{Name: "-R", Description: "remove package"},
			{Name: "-Rs", Description: "remove + unused deps"},
			{Name: "-Rns", Description: "remove + deps + configs"},
			{Name: "-Q", Description: "query installed"},
			{Name: "-Qs", Description: "search installed"},
			{Name: "-Qi", Description: "installed package info"},
			{Name: "-Yc", Description: "remove unneeded deps"},
			{Name: "-Y", Description: "yay-specific options"},
			{Name: "--aur", Description: "only AUR results"},
			{Name: "--repo", Description: "only repo results"},
			{Name: "--devel", Description: "update devel packages"},
			{Name: "--noconfirm", Description: "skip confirmations"},
			{Name: "--needed", Description: "skip up-to-date"},
		},
		Generator: installedPackageGenerator("yay"),
	})

	// paru
	core.Register(&core.Spec{
		Name:        "paru",
		Description: "AUR helper (feature-rich)",
		Options: []core.Option{
			{Name: "-S", Description: "install package"},
			{Name: "-Syu", Description: "full upgrade"},
			{Name: "-Ss", Description: "search"},
			{Name: "-R", Description: "remove"},
			{Name: "-Rs", Description: "remove + deps"},
			{Name: "-Q", Description: "query installed"},
			{Name: "--fm", Description: "file manager for PKGBUILDs"},
			{Name: "--noconfirm", Description: "skip confirmations"},
			{Name: "--aur", Description: "AUR only"},
			{Name: "--repo", Description: "repo only"},
		},
		Generator: installedPackageGenerator("paru"),
	})

	// apt
	core.Register(&core.Spec{
		Name:        "apt",
		Description: "Debian/Ubuntu package manager",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install packages"},
			{Name: "remove", Description: "remove packages", Generator: installedPackageGenerator("apt")},
			{Name: "purge", Description: "remove + config files", Generator: installedPackageGenerator("apt")},
			{Name: "autoremove", Description: "remove unused deps"},
			{Name: "update", Description: "refresh package lists"},
			{Name: "upgrade", Description: "upgrade installed"},
			{Name: "full-upgrade", Description: "upgrade with auto remove"},
			{Name: "dist-upgrade", Description: "smart upgrade"},
			{Name: "search", Description: "search packages"},
			{Name: "show", Description: "show package info"},
			{Name: "list", Description: "list packages", Options: []core.Option{
				{Name: "--installed", Description: "installed only"},
				{Name: "--upgradable", Description: "upgradable only"},
			}},
			{Name: "download", Description: "download package (no install)"},
			{Name: "clean", Description: "clear local repo"},
			{Name: "autoclean", Description: "clear old cached packages"},
			{Name: "depends", Description: "show dependencies"},
			{Name: "rdepends", Description: "show reverse dependencies"},
		},
		Options: []core.Option{
			{Name: "-y", Description: "yes to all"},
			{Name: "--no-install-recommends", Description: "skip recommends"},
			{Name: "--dry-run", Description: "simulate"},
			{Name: "-q", Description: "quiet"},
			{Name: "--fix-broken", Description: "fix broken deps"},
		},
	})

	// apt-get
	core.Register(&core.Spec{
		Name:        "apt-get",
		Description: "Debian/Ubuntu package manager (low-level)",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install packages"},
			{Name: "remove", Description: "remove packages", Generator: installedPackageGenerator("apt")},
			{Name: "purge", Description: "remove + configs", Generator: installedPackageGenerator("apt")},
			{Name: "autoremove", Description: "remove unused deps"},
			{Name: "update", Description: "refresh package lists"},
			{Name: "upgrade", Description: "upgrade installed"},
			{Name: "dist-upgrade", Description: "smart upgrade"},
			{Name: "clean", Description: "clean cache"},
			{Name: "autoclean", Description: "clean old packages"},
			{Name: "download", Description: "download only"},
			{Name: "source", Description: "get source package"},
		},
		Options: []core.Option{
			{Name: "-y", Description: "yes to all"},
			{Name: "--no-install-recommends", Description: "skip recommends"},
			{Name: "-q", Description: "quiet"},
			{Name: "-s", Description: "simulate"},
			{Name: "--fix-missing", Description: "fix missing"},
		},
	})

	// dnf (Fedora/RHEL)
	core.Register(&core.Spec{
		Name:        "dnf",
		Description: "Fedora/RHEL package manager",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install packages"},
			{Name: "remove", Description: "remove packages", Generator: installedPackageGenerator("dnf")},
			{Name: "update", Description: "update packages"},
			{Name: "upgrade", Description: "upgrade system"},
			{Name: "search", Description: "search packages"},
			{Name: "info", Description: "show package info"},
			{Name: "list", Description: "list packages"},
			{Name: "autoremove", Description: "remove unneeded"},
			{Name: "clean", Description: "clean cache", Subcommands: []core.Subcommand{
				{Name: "all", Description: "clean all"},
				{Name: "packages", Description: "clean packages"},
				{Name: "metadata", Description: "clean metadata"},
			}},
			{Name: "repolist", Description: "list repos"},
			{Name: "history", Description: "transaction history"},
			{Name: "group", Description: "manage package groups"},
		},
		Options: []core.Option{
			{Name: "-y", Description: "yes to all"},
			{Name: "-q", Description: "quiet"},
			{Name: "--best", Description: "install best"},
			{Name: "--allowerasing", Description: "allow package replacement"},
			{Name: "--no-best", Description: "allow older versions"},
		},
	})

	// yum (RHEL/CentOS legacy)
	core.Register(&core.Spec{
		Name:        "yum",
		Description: "RHEL/CentOS package manager (legacy)",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install packages"},
			{Name: "remove", Description: "remove packages", Generator: installedPackageGenerator("yum")},
			{Name: "update", Description: "update packages"},
			{Name: "search", Description: "search packages"},
			{Name: "info", Description: "show package info"},
			{Name: "list", Description: "list packages"},
			{Name: "clean", Description: "clean cache"},
			{Name: "repolist", Description: "list repos"},
		},
		Options: []core.Option{
			{Name: "-y", Description: "yes to all"},
			{Name: "-q", Description: "quiet"},
		},
	})

	// brew (macOS / Linux Homebrew)
	core.Register(&core.Spec{
		Name:        "brew",
		Description: "Homebrew package manager",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install formula/cask"},
			{Name: "uninstall", Description: "remove formula", Generator: installedPackageGenerator("brew")},
			{Name: "reinstall", Description: "reinstall formula", Generator: installedPackageGenerator("brew")},
			{Name: "update", Description: "update Homebrew"},
			{Name: "upgrade", Description: "upgrade packages"},
			{Name: "search", Description: "search formulas"},
			{Name: "info", Description: "show package info"},
			{Name: "list", Description: "list installed"},
			{Name: "outdated", Description: "list outdated"},
			{Name: "doctor", Description: "check system"},
			{Name: "cleanup", Description: "remove old versions"},
			{Name: "tap", Description: "add tap repo"},
			{Name: "untap", Description: "remove tap"},
			{Name: "services", Description: "manage services", Subcommands: []core.Subcommand{
				{Name: "start", Description: "start service"},
				{Name: "stop", Description: "stop service"},
				{Name: "restart", Description: "restart service"},
				{Name: "list", Description: "list services"},
			}},
			{Name: "cask", Description: "manage GUI apps", Subcommands: []core.Subcommand{
				{Name: "install", Description: "install cask"},
				{Name: "uninstall", Description: "remove cask"},
				{Name: "list", Description: "list casks"},
			}},
			{Name: "pin", Description: "pin formula version", Generator: installedPackageGenerator("brew")},
			{Name: "unpin", Description: "unpin formula", Generator: installedPackageGenerator("brew")},
			{Name: "link", Description: "create symlinks", Generator: installedPackageGenerator("brew")},
			{Name: "unlink", Description: "remove symlinks", Generator: installedPackageGenerator("brew")},
		},
		Options: []core.Option{
			{Name: "--cask", Description: "target GUI app"},
			{Name: "--formula", Description: "target formula"},
			{Name: "--force", Description: "force operation"},
			{Name: "--verbose", Description: "verbose"},
			{Name: "--dry-run", Description: "simulate"},
		},
	})

	// snap
	core.Register(&core.Spec{
		Name:        "snap",
		Description: "snap package manager",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install snap"},
			{Name: "remove", Description: "remove snap"},
			{Name: "refresh", Description: "update snaps"},
			{Name: "list", Description: "list installed"},
			{Name: "find", Description: "search snaps"},
			{Name: "info", Description: "show snap info"},
			{Name: "run", Description: "run snap app"},
			{Name: "connect", Description: "connect snap plug"},
			{Name: "disconnect", Description: "disconnect snap plug"},
			{Name: "revert", Description: "revert to previous"},
			{Name: "enable", Description: "enable snap"},
			{Name: "disable", Description: "disable snap"},
		},
		Options: []core.Option{
			{Name: "--classic", Description: "classic confinement"},
			{Name: "--beta", Description: "beta channel"},
			{Name: "--edge", Description: "edge channel"},
			{Name: "--channel", Description: "specific channel"},
		},
	})

	// flatpak
	core.Register(&core.Spec{
		Name:        "flatpak",
		Description: "flatpak package manager",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install application"},
			{Name: "uninstall", Description: "remove application"},
			{Name: "update", Description: "update applications"},
			{Name: "list", Description: "list installed"},
			{Name: "search", Description: "search applications"},
			{Name: "info", Description: "show app info"},
			{Name: "run", Description: "run application"},
			{Name: "remote-add", Description: "add remote"},
			{Name: "remote-delete", Description: "delete remote"},
			{Name: "remote-list", Description: "list remotes"},
			{Name: "repair", Description: "repair installation"},
		},
		Options: []core.Option{
			{Name: "--user", Description: "user installation"},
			{Name: "--system", Description: "system installation"},
			{Name: "-y", Description: "non-interactive"},
		},
	})
}
