//+build !test

package nats

import (
	"log"
	"os"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"github.com/nats-io/nats.go"
)

// InitConnection is a function to initalize a nats connection with setup options
func InitConnection() {
	nc, err := nats.Connect(config.Conf.NatsServer.URL, setupOptions()...)
	if err != nil {
		log.Println("error while connecting to", err)
		os.Exit(1)
	}
	log.Println("NATS connected successfuly")
	config.NC = nc
	//log.Println(config.Conf)
}

func setupOptions() []nats.Option {
	opts := make([]nats.Option, 0)

	opts = append(opts, nats.Name("XS-Onboarding-Accounts"))

	opts = append(opts, nats.UserInfo(config.Conf.NatsServer.Username, config.Conf.NatsServer.Password))
	opts = append(opts, nats.Token(config.Conf.NatsServer.Token))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		if !nc.IsClosed() {
			log.Println("Disconnected from NATS will attempt reconnects in ", nats.DefaultReconnectWait.Seconds())
		}
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Println("Reconnected to NATS on URL :" + nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		if !nc.IsClosed() {
			log.Println("Exiting : No servers available")
			os.Exit(1)
		} else {
			log.Println("Exiting")
			os.Exit(1)
		}
	}))
	return opts
}
