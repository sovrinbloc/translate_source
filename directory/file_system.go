package directory

import (
	"os"
	"strings"
	"fmt"
	"io/ioutil"
	"translate_source/config"
)

type FileDirectory struct {
	Directories map[string]bool
	Directory   string
	Source      map[string]string
}

func NewFileStructure(directory string) *FileDirectory {
	l := NewLocationScan(true)
	l.AddIgnoreFile(config.Env.Ignore...)
	l.AddWhitelistFile(config.Env.Whitelist...)
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
		fmt.Println(err)
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
	fmt.Println("added files")
}

func (f *FileDirectory) CreateFolderStructure(prefix string) {
	fmt.Println("Creating Directories")

	for key, _ := range f.Directories {
		f.SaveDirectory(prefix + key)
		fmt.Println(">>>", prefix+key)

	}
}

func (f *FileDirectory) SaveFiles(prefix string) {
	var err error
	fmt.Println("Writing Files")
	for file, source := range f.Source {
		//f.SaveDirectory(prefix + file)
		fmt.Println(">>>", prefix+file)
		err = ioutil.WriteFile(prefix +file,  []byte(source), 0644)
		check(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}