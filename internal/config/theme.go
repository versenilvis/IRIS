package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/versenilvis/iris/assets"
)

var Theme ThemeStyles

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

func WriteThemeFiles(path string) error {
	embedThemeDir := "themes"
	files, err := assets.Themes.ReadDir(embedThemeDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		data, err := assets.Themes.ReadFile(embedThemeDir + "/" + file.Name())
		if err != nil {
			return err
		}

		themeFile, err := os.Create(filepath.Join(path, file.Name()))
		if err != nil {
			return err
		}

		defer themeFile.Close()
		_, err = themeFile.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetThemeFile(path string, theme string) error {
	var themeStyles ThemeStyles
	_, err := toml.DecodeFile(filepath.Join(path, theme+".toml"), &themeStyles)
	if err != nil {
		return err
	}

	Theme = themeStyles

	return nil
}
