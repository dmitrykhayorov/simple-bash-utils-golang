package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
)

type options struct {
	lines, chars, words bool
}

func countLines(file *os.File) (count uint64) {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		count++
	}
	return count
}

func countWords(file *os.File) (count uint64) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		count++
	}

	return count
}

func countChars(file *os.File) (count uint64) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		count++
	}
	return count
}

func countStatistics(fileName string, fl options, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	count := uint64(0)
	file, err := os.Open(fileName)

	if err != nil {
		return err
	}

	if fl.lines {
		count = countLines(file)
	} else if fl.words {
		count = countWords(file)
	} else if fl.chars {
		count = countChars(file)
	}

	fmt.Println(count, "\t", fileName)

	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func checkFlags(fl *options) (err error) {
	count := 0

	if fl.lines {
		count++
	}
	if fl.words {
		count++
	}
	if fl.chars {
		count++
	}

	if count == 0 {
		fl.words = true
	} else if count > 1 {
		return errors.New("more than one flag specified")
	}

	return nil
}

func readFlags(fl *options) {
	flag.BoolVar(&fl.lines, "l", false, "count lines")
	flag.BoolVar(&fl.chars, "m", false, "count characters")
	flag.BoolVar(&fl.words, "w", false, "count words")
	flag.Parse()
}

func main() {
	var fl options
	var err error
	var wg sync.WaitGroup
	readFlags(&fl)
	err = checkFlags(&fl)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	for _, filename := range flag.Args() {
		wg.Add(1)
		go func(filename string) {
			err := countStatistics(filename, fl, &wg)
			if err != nil {
				fmt.Println(err)
				return
			}
		}(filename)
	}

	wg.Wait()
}
