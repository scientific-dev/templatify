package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func RunPreScripts() {
	config := command.ProcessConfig()

	for _, preScript := range config.PreScripts {
		executeCommand(strings.Split(preScript, " "))
	}

	command.Println("[36mINFO[0m Ran all pre-scripts.")
}

func MakeInit() {
	cwd = command.Cwd()
	pathname := path.Join(cwd, "templatify.config.json")
	if _, err := os.Stat(pathname); err == nil {
		command.LogError("Config file already exists.")
	}

	name := filepath.Base(cwd)
	str := "{\n    \"name\": \"" + name + "\",\n    \"description\": \"\",\n    \"ignore\": [\"test\"]\n}"
	if err := ioutil.WriteFile(pathname, []byte(str), 0644); err != nil {
		command.LogError("Failed creating the config file: " + err.Error())
	}
}

func GetGHBaseUrl(url string) string {
	if strings.HasPrefix(url, "https://github.com") {
		return url
	}

	split := strings.Split("https://github.com/"+url, "/")

	if len(split) < 5 {
		command.LogError("Could not parse the github url: https://github.com/" + url)
	}

	return split[3] + "/" + split[4]
}
