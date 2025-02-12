package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Address struct {
	Host string
	Port int
}
type Options struct {
	listenHost  string `default:"0.0.0.0"`
	listenPort  string `default:"8080"`
	remoteHost  string `default:"8.8.8.8"`
	remotePort  int    `default:53`
	dialTimeout int    `default:4`
	TrafficType string `default:"TCP"` // TCP UDP TCP/UDP
}

func (a *Address) String() string {
	return fmt.Sprintf(`%s:%d`, a.Host, a.Port)
}

func main() {
	listenHost := flag.String("lHost", "0.0.0.0", "listen host")
	listenPort := flag.String("lPort", "8080", "listen port ( can be a port range like 1024-5555 too )")
	remoteHost := flag.String("rHost", "8.8.8.8", "remote host")
	remotePort := flag.Int("rPort", 53, "remote port ( will be equal to lPort when doing range forward )")
	dialTimeout := flag.Int("timeout", 4, "dial timeout in seconds")
	help := flag.Bool("help", false, "print help")
	h := flag.Bool("h", false, "")
	timeout := 5 * time.Second
	flag.Parse()
	if *help || *h {
		flag.PrintDefaults()
		return
	}
	if *dialTimeout <= 0 || *dialTimeout > 32 {
		panic("invalid dial timeout, it should be bigger than 0 and smaller than 32")
	}
	src := Address{}
	dst := Address{
		Host: *remoteHost,
		Port: *remotePort,
	}

	if strings.Contains(*listenPort, "-") {
		pRange := strings.Split(*listenPort, "-")
		if len(pRange) > 2 {
			println("[E] invalid port range format, should be like 1024-5555")
			os.Exit(1)
		}
		start, err := strconv.Atoi(pRange[0])
		if err != nil || start < 1 || start > 65534 {
			println(`[E] invalid port range start, should be a Number between 1 and 65534`)
			os.Exit(1)
		}
		end, err := strconv.Atoi(pRange[1])
		if err != nil || end <= start || end > 65534 {
			println(`[E] invalid port range start, should be a Number between 1 and 65534 ( start should be smaller than end too )`)
			os.Exit(1)
		}
		for port := start; port <= end; port++ {
			// TCP
			println(fmt.Sprintf(`[I] tcp://%s:%d <-> tcp://%s:%d`+"\n", *listenHost, port, *remoteHost, port))
			go startForwarder(Address{
				Host: *listenHost,
				Port: port,
			}, Address{
				Host: *remoteHost,
				Port: port,
			}, *dialTimeout)
			// UDP
			println(fmt.Sprintf(`[I] udp://%s:%d <-> tcp://%s:%d`+"\n", *listenHost, port, *remoteHost, port))
			sourceUDP := *listenHost + ":" + strconv.Itoa(port)
			destUDP := *remoteHost + ":" + strconv.Itoa(port)
			destinationForwarder, err := Forward_udp(sourceUDP, destUDP, timeout)
			if err != nil {
				panic(err)
			}
			onDisconnect(destinationForwarder)
		}
	} else if strings.Contains(*listenPort, ",") {
		ports := strings.Split(*listenPort, ",")
		println(fmt.Sprintf(`%s`+"\n", ports))
		for _, portStr := range ports {
			// Trim any whitespace
			portStr = strings.TrimSpace(portStr)
			port, err := strconv.Atoi(portStr)
			if err != nil || port < 1 || port > 65534 {
				println(fmt.Sprintf("[E] invalid port number %s, should be a Number between 1 and 65534", portStr))
				os.Exit(1)
			}

			// TCP
			println(fmt.Sprintf(`[I] tcp://%s:%d <-> tcp://%s:%d`+"\n", *listenHost, port, *remoteHost, port))
			go startForwarder(Address{
				Host: *listenHost,
				Port: port,
			}, Address{
				Host: *remoteHost,
				Port: port,
			}, *dialTimeout)

			// UDP
			println(fmt.Sprintf(`[I] udp://%s:%d <-> tcp://%s:%d`+"\n", *listenHost, port, *remoteHost, port))
			sourceUDP := *listenHost + ":" + strconv.Itoa(port)
			destUDP := *remoteHost + ":" + strconv.Itoa(port)
			sourceforwarder, err := Forward_udp(sourceUDP, destUDP, timeout)
			if err != nil {
				panic(err)
			}
			onDisconnect(sourceforwarder)
		}
	} else {
		// TCP
		n, err := strconv.Atoi(*listenPort)
		if err != nil {
			println("[E] invalid listenPort, should be a Number between 1 and 65534")
			os.Exit(1)
		}
		src.Host = *listenHost
		src.Port = n
		println(fmt.Sprintf(`[I] tcp://%s <-> tcp://%s`+"\n", src.String(), dst.String()))
		go startForwarder(src, dst, *dialTimeout)

		// UDP
		println(fmt.Sprintf(`[I] udp://%s:%d <-> tcp://%s`+"\n", *listenHost, n, dst.String()))
		sourceUDP := *listenHost + ":" + strconv.Itoa(n)
		sourceforwarder, err := Forward_udp(sourceUDP, dst.String(), timeout)
		if err != nil {
			panic(err)
		}
		onDisconnect(sourceforwarder)
	}
	select {}
}

func onDisconnect(client *Forwarder_udp) {
	client.OnDisconnect(func(addr string) {
		log.Println("Client disconnected:", addr)
	})
}

func startForwarder(src, dst Address, dialTimeout int) {
	forwarder := NewForwarder(src, dst, dialTimeout)
	forwarder.Start()
}
