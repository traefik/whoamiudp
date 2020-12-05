package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	port string
	name string
)

func init() {
	flag.StringVar(&port, "port", ":8080", "give me a port number")
	flag.StringVar(&name, "name", "", "name")
}

func main() {
	flag.Parse()
	fmt.Printf("Starting up on port %s", port)

	addrL, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatalf("Invalid UDP address: %s", port)
	}

	listener, err := net.ListenUDP("udp", addrL)
	if err != nil {
		log.Fatalf("Error listening on %s: %v", port, err)
	}

	for {
		buffer := make([]byte, 1024)
		n, addrD, err := listener.ReadFrom(buffer)
		if err != nil {
			log.Println(err)
			return
		}

		temp := strings.TrimSpace(string(buffer[:n]))
		if temp == "WHO" {
			_, err := listener.WriteTo([]byte(whoAmIInfo()), addrD)
			if err != nil {
				log.Println(err)
			}
		} else {
			_, err := listener.WriteTo([]byte(fmt.Sprintf("Received: %s", buffer[:n])), addrD)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func whoAmIInfo() string {
	var out strings.Builder

	if len(name) > 0 {
		out.WriteString(fmt.Sprintf("Name: %s\n", name))
	}

	hostname, _ := os.Hostname()
	out.WriteString(fmt.Sprintf("Hostname: %s\n", hostname))

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			out.WriteString(fmt.Sprintf("IP: %s\n", ip))
		}
	}

	return out.String()
}
