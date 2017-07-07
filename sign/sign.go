package sign

import (
	"os"
	"os/signal"
	"syscall"
)

func RegSignINT(f func() error) {
	regSign(f, syscall.SIGINT)
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
