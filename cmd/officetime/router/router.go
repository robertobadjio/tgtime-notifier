package router

import (
	"cloud-time-tracker/cmd/officetime/api"
	"flag"
	"gopkg.in/routeros.v2"
	"log"
)

var (
	properties = flag.String("properties", "mac-address", "Properties")
)

func Router(router api.Router) []string { // TODO: Переименовать
	var macAddresses []string
	flag.Parse()

	c, err := routeros.Dial(router.Address, router.Login, router.Password)
	if err != nil {
		log.Fatal(err)
	}

	reply, err := c.Run("/interface/wireless/registration-table/print", "=.proplist="+*properties)
	if err != nil {
		log.Fatal(err)
	}

	for _, re := range reply.Re {
		macAddresses = append(macAddresses, re.List[0].Value)
	}

	return macAddresses
}
