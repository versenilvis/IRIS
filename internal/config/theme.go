package config

import (
	"errors"
	"os"
	"reflect"

	"github.com/BurntSushi/toml"
)

var defaultTheme = ThemeStyles{
	Border:     "#A277FF",
	Accent:     "#61FFCA",
	Muted:      "#77738A",
	Text:       "#EDECEE",
	TextSel:    "#FFFFFF",
	Match:      "#72FFD3",
	Desc:       "#A7A3B8",
	DescSel:    "#FFFFFF",
	SelBg:      "#4A426D",
	ScrollInfo: "#B58CFF",
	GhostText:  "#66636F",
	History:    "#3BAE96",
	HistorySel: "#61FFCA",
	Sys:        "#B9576C",
	SysSel:     "#FF677D",
	Alias:      "#74A965",
	AliasSel:   "#A6FF8F",
}

var Theme = defaultTheme

type ThemeStyles struct {
	Border     string `toml:"border"`
	Accent     string `toml:"accent"`
	Muted      string `toml:"muted"`
	Text       string `toml:"text"`
	TextSel    string `toml:"text_sel"`
	Match      string `toml:"match"`
	Desc       string `toml:"desc"`
	DescSel    string `toml:"desc_sel"`
	SelBg      string `toml:"sel_bg"`
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
		if !errors.Is(err, os.ErrNotExist) {
			// unreadable for other reasons, still use defaults
		}
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
