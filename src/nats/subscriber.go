//+build !test

package nats

import (
	"log"
	"os"
	"sync"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
	"github.com/nats-io/nats.go"
)

//printMsg : To print when a msg is recieved
func printMsg(m *nats.Msg, i int) {
	defer util.Panic()
	log.Printf("[#%d] Received on [%s] Queue[%s] Pid[%d]: '%s'", i, m.Subject, m.Sub.Queue, os.Getpid(), string(m.Data))
}

//Subscribe : This function is used to initiate subscriber
func Subscribe() {
	defer util.Panic()
	var wg sync.WaitGroup
	nc := config.NC
	i := 0
	wg.Add(1)
	var payload []byte
	subject := config.Conf.NatsServer.Subject
	queue := config.Conf.NatsServer.Queue
	nc.QueueSubscribe(subject+".accounts.*", queue, func(msg *nats.Msg) {
		i++
		printMsg(msg, i)
		go func(count int) {
			switch msg.Subject {
			case subject + ".accounts.checkemails":
				payload = CheckEmails(msg.Data)
			case subject + ".accounts.department":
				payload = GetDepartmentName(msg.Data)
			}
			err := msg.Respond(payload)
			if err != nil {
				log.Println(err)
			}
		}(i)
	})
	log.Println("listening on XSOnBoarding accounts service")
	wg.Wait()
}
