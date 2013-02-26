package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const cmdsName = "/cmds.sh"
const resetName = "reset.sh"

var mainCmds, resetCmds []string
var workingDir string
var in *bufio.Reader

var (
	dirPath = flag.String("f", "./", "Relative path to directory for command file lookup")
)

func main() {
	setup()
	loop()
}

func setup() {
	flag.Parse()
	mainCmdsPath := *dirPath + cmdsName
	resetCmdsPath := *dirPath + resetName
	var err error
	mainCmds, err = cmdsFromPath(mainCmdsPath)
	if err != nil {
		panic(err.Error())
	}
	resetCmds, err = cmdsFromPath(resetCmdsPath)
	if err != nil {
		panic(err.Error())
	}
	workingDir, err = os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	in = bufio.NewReader(os.Stdin)
}

func loop() {
	cmdIdx := 0
	OUTER_LOOP:
	for {
		if cmdIdx > len(mainCmds) {
			goto RESET
		}
		os.Chdir(workingDir)
		runAll(resetCmds)
		println("----------------------------")
		run(mainCmds, 0, cmdIdx, true)
		for {
			c := pause(in)
			switch c {
			case 'j', '\n':
				cmdIdx++
				goto OUTER_LOOP
			case 'k':
				cmdIdx--
				goto OUTER_LOOP
			case 'e':
				goto RESET
			default:
				println(c)
				continue
			}
		}
	}
	RESET:
	shouldReset()
}

func shouldReset() {
	println("Run reset.sh one last time(y/n)?")
	c := pause(in)
	if c == 'y' || c == 'Y' {
		os.Chdir(workingDir)
		runAll(resetCmds)
	}
}

func pause(in *bufio.Reader) byte {
	c, err := in.ReadBytes('\n')
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	return c[0]
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
	run(cmds, 0, len(cmds), false)
}

func run(cmds []string, from, to int, vocal bool) {
	for i := from; i < to; i++ {
		name, args := fmtCommand(cmds[i])
		print(name, " ")
		for _, arg := range args {
			print(arg, " ")
		}
		println()
		if name == "cd" {
			os.Chdir(args[0])
		} else {
			out, err := exec.Command(name, args...).Output()
			print(fmt.Sprintf("%s", out))
			if err != nil {
				panic(err.Error())
			}
		}
		println()
	}
}

func fmtCommand(cmd string) (string, []string) {
	trimCmd := strings.Trim(cmd, "\n")
	ss := strings.Split(trimCmd, " ")
	return ss[0], ss[1:len(ss)]
}
