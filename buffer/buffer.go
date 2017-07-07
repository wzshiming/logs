package buffer

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"

	"github.com/wzshiming/ffmt"
	"github.com/wzshiming/tee/split"
)

type Buffer struct {
	stdin    *bufio.Reader
	stdout   *bufio.Writer
	fileout  []*bufio.Writer
	closes   []func() error
	fn       []string
	flag     int
	markSwit chan int
	sync.Mutex
}

func NewBuffer(flag int, fn ...string) (*Buffer, error) {
	l := &Buffer{
		stdin:    bufio.NewReader(os.Stdin),
		stdout:   bufio.NewWriter(os.Stdout),
		fn:       fn,
		flag:     flag,
		markSwit: make(chan int, 1),
	}

	return l, l.Mklogs()
}

func (l *Buffer) On() {
	select {
	case l.markSwit <- 1:
		if len(l.markSwit) == 1 {
			time.Sleep(time.Second / 10)
			go func() {
				l.Flush()
				<-l.markSwit
			}()
		}
	default:
	}
}

func (l *Buffer) Write(p []byte) (n int, err error) {
	l.Lock()
	defer l.Unlock()
	n, err = l.stdout.Write(p)
	if err != nil {
		ffmt.Mark(err)
		return
	}
	for _, v := range l.fileout {
		n, err = v.Write(p)
		if err != nil {
			ffmt.Mark(err)
			return
		}
	}
	l.On()
	return
}

func (l *Buffer) Flush() (err error) {
	l.Lock()
	defer l.Unlock()
	if l.stdout.Buffered() == 0 {
		return nil
	}
	err = l.stdout.Flush()
	if err != nil {
		ffmt.Mark(err)
		return
	}
	for _, v := range l.fileout {
		err = v.Flush()
		if err != nil {
			ffmt.Mark(err)
			return
		}
	}
	return
}

func (l *Buffer) Copy() error {
	_, err := io.Copy(l, l.stdin)
	return err
}

func (l *Buffer) Mklogs() error {
	l.Lock()
	defer l.Unlock()
	fs := make([]*bufio.Writer, 0, len(l.fn))
	efs := make([]func() error, 0, len(l.fn))

	// 移动老的日志
	for _, v := range l.fn {
		split.Mv(v)
	}

	// 创建新日志文件
	for _, v := range l.fn {
		f, err := os.OpenFile(v, l.flag, 0666)
		if err != nil {
			return err
		}
		fs = append(fs, bufio.NewWriter(f))
		efs = append(efs, f.Close)
	}

	l.closes, efs = efs, l.closes
	l.fileout, fs = fs, l.fileout

	// 将剩余未输入的给输入
	for _, v := range fs {
		v.Flush()
	}
	for _, v := range efs {
		v()
	}
	return nil
}

// 写入缓冲区
func (l *Buffer) Run() error {
	l.Copy()
	return nil
}
