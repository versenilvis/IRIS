package integration

import "strings"

var iconMap = map[string]string{
	"git":            "󰊢",
	"docker":         "",
	"docker-compose": "",
	"go":             "",
	"golang":         "",
	"python":         "",
	"python3":        "",
	"pip":            "",
	"node":           "",
	"npm":            "",
	"npx":            "",
	"bun":            "",
	"yarn":           "",
	"rust":           "",
	"cargo":          "",
	"java":           "",
	"mvn":            "",
	"gradle":         "",
	"nvim":           "",
	"vim":            "",
	"vi":             "",
	"cd":             "",
	"ls":             "",
	"eza":            "",
	"tree":           "",
	"pwd":            "",
	"cat":            "",
	"less":           "",
	"more":           "",
	"bat":            "",
	"grep":           "",
	"ripgrep":        "",
	"find":           "",
	"alias":          "",
	"history":        "",
	"system":         "",
	"root":           "",
}

func lookupIcon(key string) string {
	key = strings.ToLower(strings.TrimSpace(key))
	if icon, ok := iconMap[key]; ok {
		return icon
	}
	if len(key) > 0 && key[0] >= '0' && key[0] <= '9' {
		return ""
	}
	return "❯"
}
