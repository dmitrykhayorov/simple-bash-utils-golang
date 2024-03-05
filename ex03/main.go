package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func readFlag() (dir string) {
	flag.StringVar(&dir, "a", "", "set directory to save logs")
	flag.Parse()
	return dir
}

func archivateFile(fileName string, dir string, wg *sync.WaitGroup) {
	defer wg.Done()
	shortname := filepath.Base(fileName)
	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	if dir == "" {
		dir = filepath.Dir(fileName)
	}

	info, _ := os.Stat(fileName)
	timestamp := strconv.FormatInt(info.ModTime().Unix(), 10)

	archiveName := shortname + "_" + timestamp + ".tar.gz"
	archive, err := os.Create(dir + "/" + archiveName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer archive.Close()

	gw := gzip.NewWriter(archive)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	fileInfo, _ := os.Stat(fileName)
	header, _ := tar.FileInfoHeader(fileInfo, shortname)
	header.Name = fileName
	err = tw.WriteHeader(header)

	_, err = io.Copy(tw, file)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	dir := readFlag()
	args := flag.Args()
	var wg sync.WaitGroup

	for _, log := range args {
		wg.Add(1)
		go func(log string) {
			archivateFile(log, dir, &wg)
		}(log)
	}
	wg.Wait()
}
