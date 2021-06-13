package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/izanagi1995/reverse-proxy-docker/config"
	"github.com/izanagi1995/reverse-proxy-docker/proxy"
)

func main() {

	conf, err := config.ReadFromFile()
	if err != nil {
		panic(err)
	}
	fmt.Println("Config")
	rp, err := proxy.NewReverseProxyManager(conf)
	if err != nil {
		panic(err)
	}
	fmt.Println("RP OK")
	// We catch ctrl-c to clean up
	signal.Notify(rp.CloseChan, os.Interrupt, syscall.SIGTERM)

	rp.Run()
}
