package directory

import (
	"errors"
	"io/ioutil"
	"strings"
)

type LocationListInterface interface {
	AddIgnoreFile(...string)
	AddDirectory(filepath string) (map[string]string, error)
	InitIgnore()
	GetSources() map[string]string
	GetContents(location string) ([]byte, error)
	GetFileList() []string
	GetDirectoryList() []string
	AddWhitelistFile(...string)
	EnableWhitelist(enable bool)
}

type LocationScan struct {
	Files           map[string]string
	Directories     map[string]bool
	ignore          map[string]struct{}
	isWhitelistOnly bool
	whitelist       map[string]struct{}
	IsFullPath bool
}

func NewLocationScan(isFullPath bool) *LocationScan {
	l := LocationScan{
		Files:       make(map[string]string),
		Directories: make(map[string]bool),
		ignore:      make(map[string]struct{}),
		IsFullPath:  isFullPath,
	}
	l.InitIgnore()
	return &l
}

func (l *LocationScan) AddDirectory(filepath string) (map[string]string, error) {
	var tmp string
	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			if l.IsIgnoredFile(f.Name()) {
				tmp = filepath + "/" + f.Name()
				if _, ok := l.Files[tmp]; ok != true {
					if l.IsFullPath {
						l.Directories[filepath+"/"+f.Name()] = false
					} else {
						l.Directories[f.Name()] = false
					}
					l.AddDirectory(tmp)
				}
			}
		} else {
			if l.CheckValidFile(f.Name()) {
				if l.IsFullPath {
					l.Files[filepath+"/"+f.Name()] = ""
				} else {
					l.Files[f.Name()] = ""
				}
			}
		}
	}
	return l.Files, nil
}

func (l *LocationScan) InitIgnore() {
	l.ignore = map[string]struct{}{"vendor": struct{}{},
		"node_modules": struct{}{},
		".git":         struct{}{},
		"fonts":        struct{}{},
		"images":       struct{}{},
		".css":         struct{}{},
		".idea":        struct{}{},
		".less":        struct{}{},
		"min.js":       struct{}{},
		"package.json": struct{}{},
		"package-lock": struct{}{},
		"Gopkg":        struct{}{},
		".DS_Store":    struct{}{},
		"LICENSE":      struct{}{},
		"README.md":    struct{}{},
		".example":     struct{}{},
	}
}

func (l *LocationScan) IsIgnoredFile(path string) bool {

	for skip, _ := range l.ignore {
		if strings.Contains(path, skip) {
			return false
		}
	}
	return true
}

func (l *LocationScan) IsWhitelistedFile(path string) bool {
	for skip, _ := range l.whitelist {
		if strings.Contains(path, skip) {
			return true
		}
	}
	return false
}

func (l *LocationScan) CheckValidFile(path string) bool {
	if l.isWhitelistOnly {
		return l.IsWhitelistedFile(path)
	}
	return l.IsIgnoredFile(path)
}

func (l *LocationScan) AddIgnoreFile(ignore ...string) {
	for _, file := range ignore {
		l.ignore[file] = struct{}{}
	}
}
func (l *LocationScan) AddWhitelistFile(ignore ...string) {
	l.isWhitelistOnly = true
	l.whitelist = make(map[string]struct{})

	for _, file := range ignore {
		l.whitelist[file] = struct{}{}
	}
}

func (l *LocationScan) EnableWhitelist(enable bool) {
	l.isWhitelistOnly = enable
}

func (l *LocationScan) GetSources() map[string]string {
	for filename, _ := range l.Files {
		content, err := l.GetContents(filename)
		if err != nil {
			panic(err)
		}
		l.Files[filename] = string(content)
	}
	return l.Files
}

func (l *LocationScan) GetContents(location string) ([]byte, error) {
	dat, err := ioutil.ReadFile(location)
	if err != nil {
		panic(errors.New("wrong filename"))
	}
	return dat, err
}

func (l *LocationScan) GetFileList() []string {
	j := make([]string, 0)
	for file, _ := range l.Files {
		j = append(j, file)
	}
	return j
}
func (l *LocationScan) GetDirectoryList() []string {
	j := make([]string, 0)
	for dir, _ := range l.Directories {
		j = append(j, dir)
	}
	return j
}
