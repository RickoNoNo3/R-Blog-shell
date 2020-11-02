package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Config struct {
	CustomBlogCoreLoc string
	CustomBlogCoreDir string
}

var GlobalConfig Config

func init() {
	GlobalConfig.CustomBlogCoreLoc = "blog.exe"
	GlobalConfig.CustomBlogCoreDir = ".///"
	if file, err := os.Open("config.json"); err == nil {
		if jsonStr, err := ioutil.ReadAll(file); err == nil {
			if err := json.Unmarshal(jsonStr, &GlobalConfig); err == nil {
				GlobalConfig.CustomBlogCoreLoc = strings.ReplaceAll(GlobalConfig.CustomBlogCoreLoc, "\\", "/")
				GlobalConfig.CustomBlogCoreDir = path.Dir(GlobalConfig.CustomBlogCoreLoc)
			}
		}
	}
}