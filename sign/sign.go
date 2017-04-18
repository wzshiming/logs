// +build !windows

package sign

import (
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
)

func SendSignUSR1(pn string, sig os.Signal) {
	sendSign(pn, syscall.SIGUSR1)
}

func sendSign(pn string, sig os.Signal) {
	exec.Command("killall", "-"+strconv.Itoa(int(sig)), pn)
}

func RegSignUSR1(f func() error) {
	regSign(f, syscall.SIGUSR1)
}

func regSign(f func() error, sig os.Signal) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, sig)
	go func() {
		for {
			<-sigs
			f()
		}
	}()
}
