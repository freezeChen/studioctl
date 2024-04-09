package util

import (
	"golang.org/x/mod/modfile"
	"os"
	"path"
	"path/filepath"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ModulePath(filename string) (string, error) {
	modBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(modBytes), nil
}

func GetGoModuleName() (string, error) {

	var module string
	var err error

	if IsExist("./go.mod") {
		module, err = ModulePath("./go.mod")
		if err != nil {
			return "", err
		}

		module = filepath.Join(module, "internal")
	} else if IsExist("../../go.mod") {
		module, err = ModulePath("../../go.mod")
		if err != nil {
			return "", err
		}
		currentPath, _ := filepath.Abs("./")
		goModPath, _ := filepath.Abs("../../go.mod")

		rel, _ := filepath.Rel(filepath.Dir(goModPath), currentPath)
		module = path.Join(module, rel, "internal")

	}
	return module, err
}
