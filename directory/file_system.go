package directory

import (
	"os"
	"strings"
	"fmt"
	"io/ioutil"
)

type FileDirectory struct {
	Directories map[string]bool
	Directory   string
	Source      map[string]string
}

func NewFileStructure(directory string) *FileDirectory {
	l := NewLocationScan(true)
	l.AddIgnoreFile("en_US", ".sql", ".key", "simplemde.js")
	l.AddWhitelistFile(".go")
	l.AddDirectory(directory)

	return &FileDirectory{
		Directories: l.Directories,
		Directory:   directory,
		Source:      make(map[string]string),
	}
}

func (f *FileDirectory) SaveDirectory(directory string) {
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (f *FileDirectory) ClipBasePath() {
	for directory, value := range f.Directories {
		truncDir := strings.Replace(directory, f.Directory, "", -1)
		if truncDir != directory {
			f.Directories[truncDir] = value
			delete(f.Directories, directory)
		}
	}
	for file, value := range f.Source {
		truncDir := strings.Replace(file, f.Directory, "", -1)
		if truncDir != file {
			f.Source[truncDir] = value
			delete(f.Source, file)
		}
	}

}

func (f *FileDirectory) AddFiles(m map[string]string) {
	f.Source = m
	f.ClipBasePath()
	fmt.Println("added files", m)
}

func (f *FileDirectory) CreateFolderStructure(prefix string) {
	for key, _ := range f.Directories {
		f.SaveDirectory(prefix + key)
	}
}

func (f *FileDirectory) SaveFiles(prefix string) {
	var err error
	for directory, source := range f.Source {
		//f.SaveDirectory(prefix + directory)
		err = ioutil.WriteFile(prefix + directory,  []byte(source), 0644)
		check(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}