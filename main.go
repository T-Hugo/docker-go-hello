// Simple program that prints the hostname of your web server
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	portPtr = flag.Int("p", 80, "The listening port")

	// IPAddr : is public IP address of the host
	IPAddr   string
	hostname string
)

func main() {
	flag.Parse()

	portStr := strconv.Itoa(*portPtr)

	var err error
	// Return external IP Address
	IPAddr, err = externalIP()
	if err != nil {
		log.Fatal(err)
	}

	// Returns the host name reported by the kernel.
	hostname, err = os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	// Create Server and Route Handlers
	http.HandleFunc("/", handler)

	srv := &http.Server{
		Addr:         ":" + portStr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		fmt.Printf("Server listening at %s:%s /...\n", IPAddr, portStr)
		fmt.Println("Hit CTRL-C to stop the server")
		log.Fatal(srv.ListenAndServe())
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func handler(w http.ResponseWriter, r *http.Request) {
	date := time.Now().Format(time.RFC3339)

	fmt.Printf("%s - - [%s] \"%s %s %s\" %d\n", r.RemoteAddr, date, r.Method, r.RequestURI, r.Proto, http.StatusOK)

	fmt.Fprintf(w, "Hello from %s\n\n", hostname)
	fmt.Fprintf(w, "Server address: %s:%d\n", IPAddr, *portPtr)
	fmt.Fprintf(w, "Server name: %s\n", hostname)
	fmt.Fprintf(w, "Date: %s\n", date)
	fmt.Fprintf(w, "URI: %s\n", r.RequestURI)
	//fmt.Fprintf(w, "UserAgent: %s\n", r.UserAgent())
}

// ExternalIP retrieves the public IP of the host
// See https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("could not determine IP Address")
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	fmt.Println("Server is shutting down...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}

	fmt.Println("Server stopped")
	os.Exit(0)
}
