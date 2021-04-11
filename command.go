package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mattn/go-colorable"
)

var stdout io.Writer = colorable.NewColorableStdout()

type Command struct {
	Args []string
}

type Config struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	PreScripts  []string          `json:"preScripts"`
	Scripts     map[string]string `json:"scripts"`
	Ignore      []string          `json:"ignore"`
	Delimiter   string            `json:"delimiter"`
	ParseFiles  string            `json:"parseFiles"`
}

func (self Config) ProcessIgnoreMatches() {
	for _, ignore := range config.Ignore {
		ignoreMatches = append(ignoreMatches, parseIgnoreMatch(ignore))
	}
}

func (self Command) Cwd() string {
	cwd, err := os.Getwd()

	if err != nil {
		self.LogError("Failed getting your current working directory with the reason: " + err.Error())
	}

	return cwd
}

func (self Command) Dirname() string {
	return filepath.Dir(os.Args[0])
}

func (self Command) ProcessConfig() Config {
	rawContent, readErr := ioutil.ReadFile(path.Join(self.Cwd(), "templatify.config.json"))
	var config Config

	if readErr != nil {
		self.LogWarn("Could not find a \"templatify.config.json\" file to read configuration. Taking default configuration.")
		var name string
		args := self.Args

		if len(args) < 1 {
			name = filepath.Base(self.Dirname())
		} else {
			name = args[0]
		}

		config = Config{
			Name: name,
		}
	} else {
		parseErr := json.Unmarshal(rawContent, &config)

		if parseErr != nil {
			self.LogError("Failed parsing json content of the file: " + parseErr.Error())
		}
	}

	if strings.Contains(config.Name, " ") {
		oldName := config.Name
		config.Name = strings.Split(config.Name, " ")[0]
		command.LogWarn("Improper template name containing a space. Auto change from \"" + oldName + "\" to \"" + config.Name + "\"")
	}

	for _, ignoreMatch := range config.Ignore {
		ignoreMatches = append(ignoreMatches, parseIgnoreMatch(ignoreMatch))
	}

	return config
}

func (self Command) Flags() map[string]string {

	flags := map[string]string{}

	for _, v := range self.Args {
		if strings.HasPrefix(v, "--") {
			flag := strings.Split(v, "=")
			flagName := flag[0][2:]
			if len(flag) == 2 {
				flags[flagName] = flag[1]
			} else {
				flags[flagName] = ""
			}
		}
	}

	return flags

}

func (self Command) LogError(message string) {
	fmt.Fprintln(stdout, "[31mERROR[0m "+message)
	os.Exit(0)
}

func (self Command) LogWarn(message string) {
	fmt.Fprintln(stdout, "[31mWARN[0m "+message)
}

func (self Command) Print(a ...interface{}) {
	fmt.Fprint(stdout, a...)
}

func (self Command) Println(a ...interface{}) {
	fmt.Fprintln(stdout, a...)
}

func (self Command) CreateTemplateManager(dirname string) TemplateManager {
	manager := TemplateManager{dirname}
	manager.Init()
	return manager
}

func (self Command) Confirm(text string) bool {
	return self.Prompt(text) == "y"
}

func (self Command) Prompt(text string) string {
	scanner := bufio.NewScanner(os.Stdin)
	self.Print(text)

	if scanner.Scan() {
		return scanner.Text()
	}

	return ""
}

func UnknownCommand(cmd string) {
	command.Args = os.Args[1:]
	flags := command.Flags()
	_, isVersionFlag := flags["version"]

	if isVersionFlag {
		command.Println("[1mCurrent templatify version:[0m 1.1.0")
		return
	} else {
		command.LogError("Unknown command: \"" + cmd + "\"")
	}
}

func CreateCommand() Command {
	Args := os.Args[2:]
	return Command{Args}
}
