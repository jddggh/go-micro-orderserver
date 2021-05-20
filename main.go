package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	httpServer "github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/selector"
	"github.com/asim/go-micro/v3/server"
	"github.com/orderserver/routers"
)

const (
	SERVER_NAME = "order-server" // server name
)

var consulReg registry.Registry

func init() {
	consulReg = consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)
}

func main() {
	srv := httpServer.NewServer(
		server.Name(SERVER_NAME),
		server.Address(":18003"),
	)

	ginRouter := routers.InitRouters()

	hd := srv.NewHandler(ginRouter)
	if err := srv.Handle(hd); err != nil {
		log.Fatalln(err)
	}

	service := micro.NewService(
		micro.Server(srv),
		micro.Registry(consulReg),
	)

	userServerAddr := GetUserServerAddr("user-server")
	if len(userServerAddr) <= 0 {
		fmt.Println("user server address is empty")
	} else {
		url := "http://"+userServerAddr+"/users"
		response, _ := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer([]byte("")))
		fmt.Println(response)
	}

	service.Init()
	service.Run()
}

func GetUserServerAddr(serverName string) (address string) {
	var retryCount int
	for {
		servers, err := consulReg.GetService(serverName)
		if err != nil {
			log.Println(err.Error())
		}
		var services []*registry.Service

		for _, value := range servers {
			fmt.Println(value.Name + ":" + value.Version)
			services = append(services, value)
		}
		next := selector.RoundRobin(services)
		if node, err := next(); err == nil {
			fmt.Println("node address:" + node.Address)
			address = node.Address
		}
		if len(address) > 0 {
			return
		}
		retryCount++
		time.Sleep(time.Second * 1)
		if retryCount >= 5 {
			return
		}

		return ""
	}
}
