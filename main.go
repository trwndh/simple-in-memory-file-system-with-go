package main

import (
	"errors"
	"fmt"
	"strings"
)

type FileSystem struct {
	IsFile  bool
	Content string
	Child   map[string]*FileSystem
}

func (fs *FileSystem) ls(path string) ([]string, error) {
	list := []string{}

	if path == "/" { // is root
		c := getAndPrintChild(fs)
		return c, nil
	}

	paths := strings.Split(strings.Trim(path, "/"), "/")
	currFs := fs // set current fs to root
	for _, p := range paths {
		fsTarget, found := currFs.Child[p]
		if !found {
			fmt.Println("target not found")
			return []string{}, errors.New("not found")
		}

		currFs = fsTarget
	}

	list = getAndPrintChild(currFs)

	return list, nil
}

func getAndPrintChild(fs *FileSystem) []string {
	c := []string{}
	for child := range fs.Child {
		c = append(c, child)
	}

	fmt.Println(c)
	return c
}

func (fs *FileSystem) save(path string, content string) error {
	paths := strings.Split(strings.Trim(path, "/"), "/")

	fileName := paths[len(paths)-1]
	targetDir := paths[:len(paths)-1]

	if len(targetDir) == 0 { // write in root
		fs.Child[fileName] = &FileSystem{
			IsFile:  true,
			Content: content,
		}

		return nil
	}

	dirFs, err := fs.mkdir(strings.Join(targetDir, "/"))
	if err != nil {
		return err
	}

	if dirFs.Child == nil {
		dirFs.Child = make(map[string]*FileSystem)
	}
	dirFs.Child[fileName] = &FileSystem{
		IsFile:  true,
		Content: content,
	}

	return nil
	// TODO implement
}

func (fs *FileSystem) read(path string) (string, error) {
	paths := strings.Split(strings.Trim(path, "/"), "/")
	targetFs := fs
	for _, p := range paths {
		objFs, found := targetFs.Child[p] // cek adakah child dari current fs yg namanya var p
		if !found {
			fmt.Println("not found")
		} else {
			targetFs = objFs
		}
	}

	fmt.Println(targetFs.Content)
	return targetFs.Content, nil
}

func (fs *FileSystem) mkdir(path string) (*FileSystem, error) { // acting like mkdir -p
	paths := strings.Split(strings.Trim(path, "/"), "/")

	targetFs := fs
	for _, p := range paths {
		objFs, found := targetFs.Child[p] // cek adakah child dari current fs yg namanya var p
		if !found {
			if targetFs.IsFile {
				errMsg := "Cannot write directory under file"
				fmt.Println(errMsg)
				return nil, errors.New(errMsg)
			}
			if targetFs.Child == nil {
				targetFs.Child = make(map[string]*FileSystem)
			}
			targetFs.Child[p] = &FileSystem{}

			targetFs = targetFs.Child[p]
		} else {
			targetFs = objFs
		}
	}

	return targetFs, nil
}

func main() {
	fs := &FileSystem{} // root
	fs.ls("/")          // returns: []
	fs.mkdir("/foo/bar")
	fs.mkdir("/foo")
	fs.save("/file.txt", "File content")
	fs.save("/foo/file.txt", "File contentxxx")
	fs.save("/foo/bar/file-2.txt", "123 content")
	fs.ls("/foo")                  // returns: [bar, file.txt]
	fs.ls("/foo/m")                // returns: target not found
	fs.read("/foo/bar/file-2.txt") // returns "123 content"
	fs.ls("/foo/bar")              // returns ["file-2.txt"]
}
