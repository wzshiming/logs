// +build !windows

package sign

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func SendSignUSR1(pn string) {
	sendSign(pn, syscall.SIGUSR1)
}

func sendSign(pn string, sig os.Signal) {
	exec.Command("killall", "-"+fmt.Sprint(sig), pn)
}

func RegSignUSR1(f func() error) {
	regSign(f, syscall.SIGUSR1)
}
