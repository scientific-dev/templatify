package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gobwas/glob"
)

var (
	lockData       Config
	parseFilesGlob glob.Glob
	parserRegex    *regexp.Regexp
	noGit          bool
)

func Use() {
	if len(command.Args) == 0 {
		command.LogError("No template name has been provided!")
	}

	flags := command.Flags()
	timeStart := time.Now().UnixNano() / 1000000
	dirname := command.Dirname()
	templateDir = path.Join(dirname, "templates", command.Args[0])
	manager := command.CreateTemplateManager(dirname)
	lockData = manager.GetTemplateInfo(command.Args[0])
	writePath := command.Cwd()
	customPath, hasCustomPath := flags["custom-path"]
	_, noGit = flags["no-git"]

	if len(lockData.Delimiter) == 0 {
		lockData.Delimiter = "%"
	}

	parserRegex = regexp.MustCompile("\\" + lockData.Delimiter + "\\{(.*?)\\}")

	if hasCustomPath {
		writePath = path.Join(writePath, customPath)
	}

	if glb, err := glob.Compile(lockData.ParseFiles); err != nil {
		command.LogError("Invalid delimiter provided!")
	} else {
		parseFilesGlob = glb
	}

	manager.CheckAbsoluteDir(writePath)
	command.Println("[36mINFO[0m Copying template \"" + command.Args[0] + "\" to \"" + writePath + "\".")
	createTemplateDir("", writePath, ReadDir(templateDir))
	if _, removeLockFile := flags["remove-lock"]; removeLockFile {
		if err := os.Remove(path.Join(writePath, "templatify.lock.json")); err != nil {
			command.LogError("Failed removing templatify.lock.json file: " + err.Error())
		}
	}

	command.Println("[36mINFO[0m Cloned files.")

	for _, preScript := range lockData.PreScripts {
		executeCommand(strings.Split(preScript, " "))
	}

	command.Println("[36mINFO[0m Ran all preScripts.")
	timeEnd := time.Now().UnixNano() / 1000000
	command.Println("[32mSUCCESS[0m Finished in " + strconv.Itoa(int(timeEnd-timeStart)/1000) + "s")
}

func createTemplateDir(pathname string, writePath string, assets []FolderAsset) {
	for _, asset := range assets {
		if asset.IsDir {
			if asset.Name == ".git" && noGit {
				continue
			}
			name := path.Join(pathname, asset.Name)
			os.Mkdir(path.Join(writePath, name), asset.Mode())
			createTemplateDir(name, writePath, asset.Assets)
		} else {
			createTemplateFile(path.Join(pathname, asset.Name), writePath, asset)
		}
	}
}

func createTemplateFile(pathname string, writePath string, asset FolderAsset) {
	content, readErr := ioutil.ReadFile(path.Join(templateDir, pathname))
	if readErr != nil {
		command.LogError("Failed reading a file: " + readErr.Error())
	}

	err := ioutil.WriteFile(path.Join(writePath, pathname), parseFile(pathname, content), asset.Mode())
	if err != nil {
		command.LogError("Error originated while storing the file in templates: " + err.Error())
	}
}

func executeCommand(args []string) {
	if len(args) == 0 {
		command.LogError("Looks like a script provided in \"preScripts\" field in config is invalid!")
	}

	process := exec.Command(args[0], args[1:]...)
	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr
	command.Println("[90mPRE-SCRIPT[0m " + strings.Join(args, " "))

	if err := process.Run(); err != nil {
		command.LogError("Failed running a preScript \"" + strings.Join(args, " ") + "\": " + err.Error())
	}
}

func parseFile(name string, data []byte) []byte {
	if parseFilesGlob.Match(name) {
		return []byte(cleanString(string(data)))
	} else {
		return data
	}
}

func cleanString(str string) string {
	result := ""
	lastIndex := 0
	values := map[string]string{}

	for _, v := range parserRegex.FindAllSubmatchIndex([]byte(str), -1) {
		var value string
		key := str[(v[0] + 2):(v[1] - 1)]

		if val, found := values[key]; found {
			value = val
		} else {
			value = command.Prompt("[90mCONFIG[0m Value for [1m" + key + "[0m: ")
		}

		result += str[lastIndex:v[0]] + value
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}
