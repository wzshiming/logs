package split

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Mv(fn string) {
	mv(fn, mklogpath(fn))
}

func mv(old string, news string) {
	os.MkdirAll(filepath.Dir(news), 0666)
	os.Rename(old, news)
}

func mklogpath(f string) string {
	ba := filepath.Base(f)
	di := strings.Index(ba, ".")
	p := ba[:di]
	ext := ba[di:]

	return logpath(p, ext, time.Now())
}

func logpath(name, ext string, ti time.Time) string {
	pa := ti.Format(time.RFC3339)
	n := filepath.Join(name, name+"-"+pa[:7], name+"-"+pa[:10]+ext)
	return n
}
