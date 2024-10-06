package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/memberlist"
	log "github.com/sirupsen/logrus"
)

func GetMemberlistConfig(configType string) (*memberlist.Config, error) {
	switch configType {
	case "local":
		return memberlist.DefaultLocalConfig(), nil

	case "lan":
		return memberlist.DefaultLANConfig(), nil

	case "wan":

		return memberlist.DefaultWANConfig(), nil
	}

	return nil, fmt.Errorf("failed to get a memerblist default config from %s", configType)
}

func WaitSignal() {
	signal_chan := make(chan os.Signal, 2)
	signal.Notify(signal_chan, syscall.SIGTERM)
	signal.Notify(signal_chan, syscall.SIGINT)

	for {
		select {
		case s := <-signal_chan:
			log.Printf("signal %s happen", s.String())
			return
		}
	}
}
