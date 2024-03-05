package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type params struct {
	symlink, dir, file bool
	ext                string
}

var options params

func CheckIfDir(path string, d fs.DirEntry, err error) error {
	// Если не удалось считать директорию, функция вызывается повторно для
	// того, чтобы передать ошибку

	if err != nil {
		return fs.SkipDir
	}

	if d.IsDir() && options.dir {
		fmt.Println(path)
	}
	info, err := os.Lstat(path)
	if err != nil {
		return fs.SkipDir
	}

	if info.Mode()&os.ModeSymlink == os.ModeSymlink && options.symlink {
		linkSource, err := os.Readlink(path)
		if err != nil {
			linkSource = "[broken]"
		}
		fmt.Println(path, "->", linkSource)
	}

	if info.Mode().IsRegular() && options.file {
		if filepath.Ext(path) == ("." + options.ext) {
			fmt.Println(path)
			return nil
		} else if len(options.ext) > 0 {
			return nil
		}
		fmt.Println(path)
	}

	return nil
}

func readFlags() {
	flag.BoolVar(&options.file, "f", false, "specify to show only files")
	flag.BoolVar(&options.dir, "d", false, "specify to show only files")
	flag.BoolVar(&options.symlink, "sl", false, "specify to show only symbolic links")
	flag.StringVar(&options.ext, "ext", "", "specify to show only files with certain extension")
	flag.Parse()
}

func checkParams() (err error) {
	if flag.Arg(0) == "" {
		return errors.New("path was not specified")
	}
	if options.ext != "" && !options.file {
		return errors.New("ext flag without f flag")
	}
	//if options.ext != "" && filepath.Ext(path) == ("."+options.ext) {
	//	fmt.Println(path)
	//	return nil
	//}
	if !options.file && !options.dir && !options.symlink {
		options.file = true
		options.dir = true
		options.symlink = true
	}

	return nil
}

func main() {
	var fileDir string
	var err error
	readFlags()
	err = checkParams()

	if err != nil {
		fmt.Println(err)
		return
	}
	fileDir = flag.Arg(0)

	//fmt.Println("filepath in main: ", fileDir)
	err = filepath.WalkDir(fileDir, CheckIfDir)
	if err != nil {
		fmt.Println(err)
	}
}
