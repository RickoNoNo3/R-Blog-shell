package processor

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"../config"
)

// 所有Processor的统一接口, 要实现处理函数
type Processor interface {
	Process(input string) (res string, err error)
}

// 用来判断在调用BlogCore时是否需要一个stdin绑定,
// 如果需要, 则会打开一条从本程序stdin到BlogCore stdin的管道.
// 管道关闭后, stdin恢复正常
var NeedStdinTable = map[string]bool{
	"draw":     false,
	"drawCore": false,
	"new":      true,
	"drag":     false,
	"edit":     true,
	"move":     true,
	"remove":   true,
	"read":     false,
}

// 把输入的一行乱七八糟字符串处理后按空格切分
func splitArgs(input string) (args []string) {
	input = strings.TrimSpace(input)
	input = regexp.MustCompile(` +`).ReplaceAllString(input, ` `)
	args = strings.Split(input, " ")
	return
}

// 调用BlogCore主程序
func runBlogCore(args []string) (output string, err error) {
	// 是否需要打开stdin通道
	needStdin := NeedStdinTable[args[0]]

	// 接受stdout和stderr输出, 以及其他的调用配置
	outputBuf := bytes.Buffer{}
	cmd := exec.Cmd{
		Dir:    config.GlobalConfig.CustomBlogCoreDir,
		Path:   config.GlobalConfig.CustomBlogCoreLoc,
		Args:   append([]string{config.GlobalConfig.CustomBlogCoreLoc}, args...),
		Stdout: &outputBuf,
		Stderr: os.Stderr,
	}

	if needStdin {
		// 需要stdin, 打开管道, 开启BlogCore, 传递输入数据.
		fmt.Println("你现在处于本命令的数据输入模式.\n在新行中输入EXIT后回车即可离开.\n--------------------------------------")
		var cmdStdin io.WriteCloser
		cmdStdin, err = cmd.StdinPipe()
		if err != nil {
			return
		}
		// 读取本程序的stdin, 注意这里会自动屏蔽main里的reader
		reader := bufio.NewReader(os.Stdin)
		if err = cmd.Start(); err != nil {
			return
		}
		// 不断向管道内转发输入
		for {
			// 读取一行
			var lineBytes []byte
			lineBytes, _, err = reader.ReadLine()
			line := string(lineBytes)
			if err != nil {
				return
			}
			// 输入EXIT时结束,
			if line == "EXIT" {
				_ = cmdStdin.Close()
				break
			}
			// 转发
			_, _ = fmt.Fprintln(cmdStdin, line)
		}
	} else {
		// 不需要stdin, 直接调用BlogCore
		if err = cmd.Start(); err != nil {
			return
		}
	}
	// 等待BlogCore运行结束
	if err = cmd.Wait(); err != nil {
		return
	}
	// 处理输出
	output = strings.TrimSpace(outputBuf.String())
	if output == "-1" {
		err = errors.New("blog: bad command arguments")
	}
	return
}
