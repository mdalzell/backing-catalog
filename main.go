package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hudl/fargo"
	"github.com/mdalzell/backing-catalog/service"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	eurekaUrl := os.Getenv("EUREKA_URL")
	if len(eurekaUrl) == 0 {
		log.Fatal("Missing ENV variable: EUREKA_URL")
	}

	c := fargo.NewConn(eurekaUrl)

	i := fargo.Instance{
		HostName:         "http://127.0.0.10:3000",
		Port:             3000,
		App:              "BACKING_CATALOG",
		IPAddr:           "127.0.0.10",
		VipAddress:       "127.0.0.10",
		SecureVipAddress: "127.0.0.10",
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
		Status:           fargo.UP,
	}

	c.RegisterInstance(&i)
	f, _ := c.GetApps()

	for key, theApp := range f {
		fmt.Println("Registered App:", key, " First Host Name:", theApp.Instances[0].HostName)
	}

	server := service.NewServer(&c)
	server.Run(":" + port)
}
