package config

import (
	"os"
	"reflect"

	"github.com/BurntSushi/toml"
)

var defaultTheme = ThemeStyles{
	Border:     "#a277ff",
	Accent:     "#61ffca",
	Muted:      "#6d6a7f",
	Text:       "#edecee",
	TextSel:    "#ffffff",
	Key:        "#a277ff",
	Match:      "#61ffca",
	Desc:       "#9692a8",
	DescSel:    "#edecee",
	SelBg:      "#3d375e",
	SelText:    "#110f18",
	ScrollInfo: "#a277ff",
	GhostText:  "#4B4A4C",
	History:    "#1a2d36",
	HistorySel: "#61ffca",
	Sys:        "#1e1d28",
	SysSel:     "#a277ff",
	Alias:      "#2a2342",
	AliasSel:   "#a277ff",
}

var Theme = defaultTheme

type ThemeStyles struct {
	Border     string `toml:"border"`
	Accent     string `toml:"accent"`
	Muted      string `toml:"muted"`
	Text       string `toml:"text"`
	TextSel    string `toml:"text_sel"`
	Key        string `toml:"key"`
	Match      string `toml:"match"`
	Desc       string `toml:"desc"`
	DescSel    string `toml:"desc_sel"`
	SelBg      string `toml:"sel_bg"`
	SelText    string `toml:"sel_text"`
	ScrollInfo string `toml:"scroll_info"`
	GhostText  string `toml:"ghost_text"`
	Sys        string `toml:"sys"`
	SysSel     string `toml:"sys_sel"`
	History    string `toml:"hist"`
	HistorySel string `toml:"hist_sel"`
	Alias      string `toml:"alias"`
	AliasSel   string `toml:"alias_sel"`
}

// LoadTheme reads theme.toml at filePath
// Missing file - use defaults silently
// Missing/empty fields - fall back to the default value for that field
func LoadTheme(filePath string) {
	Theme = defaultTheme

	data, err := os.ReadFile(filePath)
	if err != nil {
		// missing or unreadable file: use defaults
		return
	}

	var t ThemeStyles
	if _, err := toml.Decode(string(data), &t); err != nil {
		return
	}

	applyThemeWithFallback(&Theme, t, defaultTheme)
}

// applyThemeWithFallback copies non-empty string fields from src into dst,
// using def for any field that is empty in src
func applyThemeWithFallback(dst *ThemeStyles, src, def ThemeStyles) {
	dv := reflect.ValueOf(dst).Elem()
	sv := reflect.ValueOf(src)
	dv2 := reflect.ValueOf(def)

	for i := range dv.NumField() {
		sf := sv.Field(i).String()
		if sf != "" {
			dv.Field(i).SetString(sf)
		} else {
			dv.Field(i).SetString(dv2.Field(i).String())
		}
	}
}
