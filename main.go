package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"os/exec"
	"strings"
)

var mainCmds, resetCmds []string

var (
	dirPath = flag.String("f", "", "Relative path to directory for command file lookup")
)

func main() {
	setup()
	loop()
}

func setup() {
	flag.Parse()
	mainCmdsPath := *dirPath + "/mainCmds"
	resetCmdsPath := *dirPath + "/resetCmds"
	var err error
	mainCmds, err = cmdsFromPath(mainCmdsPath)
	if err != nil {
		panic(err.Error())
	}
	resetCmds, err = cmdsFromPath(resetCmdsPath)
	if err != nil {
		panic(err.Error())
	}
}

func loop() {
	in := bufio.NewReader(os.Stdin)
	cmdIdx := 0
	for {
		if cmdIdx >= len(mainCmds) {
			// TODO
		}
		runAll(resetCmds)
		run(mainCmds, 0, cmdIdx)
		switch pause(in) {
		case 'j':
			cmdIdx++
		case 'k':
			cmdIdx--
		case 'e':
			os.Exit(0)
		default:
			continue
		}
	}
}

func pause(in *bufio.Reader) byte {
	c, err := in.ReadByte()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	return c
}

func cmdsFromPath(fName string) ([]string, error) {
	f, err := os.Open(fName)
	if err != nil {
		return nil, err
	}
	cmds, err := cmdsFromFile(f)
	if err != nil {
		return nil, err
	}
	return cmds, nil
}

func cmdsFromFile(f *os.File) ([]string, error) {
	cmds := make([]string, 0)
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			return cmds, nil
		}
		if err != nil {
			return cmds, err
		}
		cmds = append(cmds, line)
	}
	panic("Unreachable")
}

func runAll(cmds []string) {
	run(cmds, 0, len(cmds))
}

func run(cmds []string, from, to int) {
	for i := from; i < to; i++ {
		name, args := fmtCommand(cmds[i])
		out, err := exec.Command(name, args...).Output()
		if err != nil {
			panic(err.Error())
		}
		println(out)
	}
}

func fmtCommand(cmd string) (string, []string) {
	ss := strings.Split(cmd, " ")
	return ss[0], ss[1:len(ss)]
}
