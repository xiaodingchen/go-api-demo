package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var configWatcher *fsnotify.Watcher
var configPath = "/Users/xiao/www/goproject/github.com/xiaodingchen/go-api-demo/configs"

func main() {
	ReadConfig(configPath)
	log.Println(viper.GetStringMapString("redis"))
	log.Println(viper.GetString("service_name"))
	go WatchConfig(configPath)
	ctx, cannel := context.WithCancel(context.Background())
	cannel()
	<-ctx.Done()
}

func ReadConfig(path string) {
	fiels, err := ioutil.ReadDir(path)
	if err != nil {
		panic(fmt.Sprintf("Fatal error config path: %s, err: %v", path, err))
	}
	ff := []string{}
	for _, f := range fiels {
		if f.IsDir() {
			continue
		}
		if strings.Index(f.Name(), ".") == 0 {
			continue
		}
		// name := f.Name()
		// ff = append(ff, name[0:strings.Index(name, ".")])
		ff = append(ff, f.Name())
	}
	if len(ff) < 1 {
		panic(fmt.Sprintf("Fatal error config path: %s, err: %v", path, "no such config file."))
	}

	viper.AddConfigPath(path)
	for i, cf := range ff {
		// viper.SetConfigFile(cf)
		name := cf[0:strings.Index(cf, ".")]
		viper.SetConfigName(name)
		if i == 0 {
			err = viper.ReadInConfig()
		} else {
			err = viper.MergeInConfig()
		}
		if err != nil {
			panic(fmt.Sprintf("Fatal error config file:%s, err: %v", filepath.Join(path, cf), err))
		}
	}

	log.Println(ff)
}

func WatchConfig(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprintf("Fatal error watch config path: %s, err: %v", path, err))
	}
	defer watcher.Close()
	err = watcher.Add(path)
	if err != nil {
		panic(fmt.Sprintf("Fatal error watch config path: %s, err: %v", path, err))
	}
	// log.Println(configWatcher)

	// go func() {
	for {
		select {
		case event := <-watcher.Events:
			// log.Printf("config is change :%s \n", event.String())
			if (event.Op&fsnotify.Write == fsnotify.Write) ||
				(event.Op&fsnotify.Create == fsnotify.Create) ||
				(event.Op&fsnotify.Remove == fsnotify.Remove) ||
				(event.Op&fsnotify.Rename == fsnotify.Rename) {
				log.Printf("config is change :%s \n", event.String())
				ReadConfig(configPath)
				log.Println(viper.GetStringMapString("redis"))
				log.Println(viper.GetString("service_name"))
			}
		case err = <-watcher.Errors:
			panic(fmt.Sprintf("Fatal error watch config path: %s, err: %v", path, err))
		}
	}
	// }()
	// select {}
}
