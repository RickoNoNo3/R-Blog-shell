package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"./myerror"
	"./processor"
)

var processors = []processor.Processor{
	processor.LinkProcessor{},
	processor.ReadProcessor{},
	processor.DefaultProcessor{},
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	cycle := true
	fmt.Println("欢迎使用R-Blog命令工具, 输入help查看帮助, 输入exit退出.")
	for cycle {
		fmt.Print(">> ")
		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(strings.Split(input, "\n")[0])
		res := ""
		switch input {
		case "":
			continue
		case "exit":
			cycle = false
		default:
			for _, p := range processors {
				if res, err = p.Process(input); err == myerror.CannotProcessError {
					continue
				} else {
					break
				}
			}
		}
		if err == nil {
			fmt.Println(res)
		} else {
			fmt.Println(err)
		}
	}
}
