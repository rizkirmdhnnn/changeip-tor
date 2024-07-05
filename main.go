package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"time"

	"golang.org/x/net/proxy"
)

const torProxy = "socks5://127.0.0.1:9050"

func main() {
	cmd := exec.Command("tor")
	err := cmd.Start()
	if err != nil {
		log.Fatal("Gagal memulai Tor:", err)
	}
	defer cmd.Process.Kill()

	// Tunggu sebentar agar Tor siap
	time.Sleep(5 * time.Second)

	for {
		// Buat koneksi Tor baru
		dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
		if err != nil {
			log.Fatal("Gagal membuat dialer:", err)
		}

		// Buat HTTP client yang menggunakan koneksi Tor
		httpClient := &http.Client{
			Transport: &http.Transport{
				Dial: dialer.Dial,
			},
		}

		// Periksa IP
		resp, err := httpClient.Get("https://api.ipify.org")
		if err != nil {
			log.Println("Gagal mendapatkan IP:", err)
		} else {
			defer resp.Body.Close()
			ip, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("Gagal membaca respons:", err)
			} else {
				fmt.Printf("IP saat ini: %s\n", string(ip))
			}
		}

		time.Sleep(3 * time.Second)

		conn, _ := net.Dial("tcp", "127.0.0.1:9051")
		fmt.Fprintf(conn, "AUTHENTICATE \"rizkirmdhn\"\r\n")
		fmt.Fprintf(conn, "SIGNAL NEWNYM\r\n")
		time.Sleep(1 * time.Second)
	}
}
