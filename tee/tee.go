package tee

import (
	"bufio"
	"os"
	"time"
)

type Tee struct {
	stdin   *bufio.Reader
	stdout  *bufio.Writer
	fileout *bufio.Writer
	fn      string
}

func NewTee(fn string) (*Tee, error) {
	l := &Tee{
		stdin:  bufio.NewReader(os.Stdin),
		stdout: bufio.NewWriter(os.Stdout),
		fn:     fn,
	}

	return l, l.Mklogs()
}

func (l *Tee) Run() error {
	go func() {
		for {
			time.Sleep(time.Second / 10)
			l.stdout.Flush()
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second / 1)
			l.fileout.Flush()
		}
	}()
	for {
		c, _, _ := l.stdin.ReadRune()
		l.stdout.WriteRune(c)
		l.fileout.WriteRune(c)
	}
	return nil
}

func (l *Tee) Mklogs() error {
	f, err := os.OpenFile(l.fn, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	l.fileout = bufio.NewWriter(f)
	return nil
}
