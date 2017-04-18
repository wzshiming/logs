package main

import (
	"flag"
	"fmt"

	"github.com/wzshiming/logs/sign"
	"github.com/wzshiming/logs/tee"
)

func main() {
	t := flag.String("f", "", "tee.log")
	s := flag.String("s", "", "nginx")

	flag.Parse()

	if *s != "" {
		sign.SendSignUSR1(*s)
	}

	if *t != "" {
		lo, err := tee.NewTee(*t)
		if err != nil {
			fmt.Println(err)
			return
		}
		sign.RegSignUSR1(lo.Mklogs)
		lo.Run()
	}

}
