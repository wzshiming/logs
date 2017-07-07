package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wzshiming/tee/buffer"
	"github.com/wzshiming/tee/sign"
)

func main() {
	a := flag.Bool("a", false, "Append to the given FILEs, don't overwrite")
	i := flag.Bool("i", false, "Ignore interrupt signals (SIGINT)")
	u := flag.Bool("u", false, "Accept the signals (SIGUSR1) to regenerate the log file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "      tee [-ai] [FILE]...")
		flag.PrintDefaults()
	}
	if *i {
		// 忽略中断信号
		sign.RegSignINT(func() error { return nil })
	}
	flag.Parse()

	fl := os.O_CREATE | os.O_WRONLY
	if *a {
		fl |= os.O_APPEND
	} else {
		fl |= os.O_TRUNC
	}
	lo, err := buffer.NewBuffer(fl, flag.Args()...)
	if err != nil {
		flag.PrintDefaults()
		fmt.Println(err)
		return
	}
	// 接收到usr1 信号 重新创建 文件
	if *u {
		sign.RegSignUSR1(lo.Mklogs)
	}

	lo.Run()
	return
}
