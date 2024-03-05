package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func readParams() (params []string) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		params = append(params, line)
	}
	//fmt.Println(params)
	return
}

func readCommands() (commands []string) {
	flag.Parse()

	commands = flag.Args()

	if len(commands) == 0 {
		panic("no command was specified")
	}

	return commands
}

func main() {
	commands := readCommands()
	params := readParams()
	for _, p := range params {
		tmpComm := make([]string, len(commands))
		copy(tmpComm, commands)
		tmpComm = append(tmpComm, p)
		res, _ := exec.Command(commands[0], tmpComm[1:]...).Output()
		fmt.Printf("%s", res)

	}
	return
}
