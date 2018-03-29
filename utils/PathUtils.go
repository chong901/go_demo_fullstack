package utils

import "github.com/kardianos/osext"

func GetProjectRoot() string {
	dir, _ := osext.ExecutableFolder()
	return dir
}
