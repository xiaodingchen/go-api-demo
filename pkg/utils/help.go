package utils

import (
    "github.com/spf13/viper"
    "path/filepath"
)

func TemplateFile(filename string) string{
    dir := viper.GetString("template.dir")
    return filepath.Join(dir, filename)
}