package main

import (
	"fmt"
	"go-changeip-tor/config"
	"go-changeip-tor/modules"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net/http"
	"time"
)

func newHttpProxy(addr string) (*http.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	httpTransport := &http.Transport{
		Dial: dialer.Dial,
	}
	httpClient := &http.Client{
		Transport: httpTransport,
	}
	return httpClient, nil
}

func getData(url string) {
	httpClient, err := newHttpProxy(config.Cfg.TORSERVER_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}
	res, err := httpClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}

func main() {

	config.LoadConfig()

	tor := modules.Tor{
		ControlAddress: config.Cfg.TORCONTROL_ADDRESS,
	}

	tor.Init()
	for {
		tor.ChangeIP()
		getData("https://api.ipify.org")
		time.Sleep(3 * time.Second)
	}
}
