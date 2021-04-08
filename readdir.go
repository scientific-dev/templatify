package main

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/gobwas/glob"
)

type FolderAsset struct {
	IsDir  bool
	Name   string
	Assets []FolderAsset
	Mode   func() os.FileMode
}

var ignoreMatches []glob.Glob = []glob.Glob{}

func ReadDir(dirname string) []FolderAsset {
	assets := []FolderAsset{}
	files, readErr := ioutil.ReadDir(dirname)

	if readErr != nil {
		command.LogError("Failed reading cwd: " + readErr.Error())
	}

	for _, file := range files {
		isFolder := file.IsDir()
		name := file.Name()
		asset := FolderAsset{
			IsDir:  isFolder,
			Name:   name,
			Assets: []FolderAsset{},
			Mode:   file.Mode,
		}

		if isFolder {
			asset.Assets = ReadDir(path.Join(dirname, name))
		}

		assets = append(assets, asset)
	}

	return assets
}

func parseIgnoreMatch(match string) glob.Glob {
	g, err := glob.Compile(match)
	if err != nil {
		command.LogError("Failed compiling ignore match \"" + match + "\": " + err.Error())
		os.Exit(0)
	}
	return g
}

func isIgnoredFile(name string) bool {
	for _, ignoreMatch := range ignoreMatches {
		if ignoreMatch.Match(name) {
			return true
		}
	}

	return false
}

func resolveTemplatesInDir(targetPath string, offsetPath string, files []os.FileInfo) []string {
	templates := []string{}

	for _, file := range files {
		if file.IsDir() {
			name := file.Name()
			newOffset := offsetPath + "/" + name
			if dirHasLock(path.Join(targetPath, newOffset)) {
				templates = append(templates, newOffset)
			} else {
				if files, err := ioutil.ReadDir(path.Join(targetPath, newOffset)); err == nil {
					templates = append(templates, resolveTemplatesInDir(targetPath, newOffset, files)...)
				}
			}
		}
	}

	return templates
}

func dirHasLock(pathname string) bool {
	if _, err := os.Stat(path.Join(pathname, "templatify.lock.json")); err != nil {
		return false
	}

	return true
}
