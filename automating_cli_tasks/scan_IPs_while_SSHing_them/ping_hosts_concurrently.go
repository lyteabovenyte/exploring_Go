package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"sync"
	"time"

	"inet.af/netaddr"
)

// a type to store the data that we gather
type record struct {
	Host      net.IP
	Reachable bool
	LoginSSH  bool
	Uname     string
}

func main() {
	_, errping := exec.LookPath(ping)
	if errping != nil {
		log.Fatal("cannot find ping in our PATH")
	}

	_, errssh := exec.LookPath(ssh)
	if errssh != nil {
		log.Fatal("cannot find ssh in our PATH")
	}

	if len(os.Args) != 2 {
		log.Fatal("error: only one argument allowed, the network CIDR to scan")
	}
	ipCh, err := hosts(os.Args[1])
	if err != nil {
		log.Fatalf("error: CIDR address did not parse: %s,", err)
	}

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	scanResults := scanPrefix(ipCh)
	unameResults := unamePrefixes(u.Username, scanResults)

	for rec := range unameResults {
		b, _ := json.Marshal(rec)
		fmt.Printf("%s\n", b)
	}
}

func hosts(cidr string) (chan net.IP, error) {
	ch := make(chan net.IP, 1)

	prefix, err := netaddr.ParseIPPrefix(cidr)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(ch)

		var last net.IP
		for ip := prefix.IP().Next(); prefix.Contains(ip); ip = ip.Next() {
			// Prevents sending the broadcast address.
			if len(last) != 0 {
				//log.Printf("sending: %s, contained: %v", last, prefix.Contains(ip))
				ch <- last
			}
			last = ip.IPAddr().IP
		}
	}()
	return ch, nil
}

func hostAlive(ctx context.Context, host net.IP) bool {
	cmd := exec.CommandContext(ctx, ping, "-c", "1", "-t", "2", host.String())

	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func runUname(ctx context.Context, host net.IP, user string) (string, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}
	login := fmt.Sprintf("%s@%s", user, host)
	cmd := exec.CommandContext(
		ctx,
		ssh,
		"-o StrictHostKeyChecking=no",
		"-o BatchMode=yes",
		login,
		"uname -a",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// take a channel contain IPs and scan them concurrently.
func scanPrefix(ipCh chan net.IP) chan record {
	ch := make(chan record, 1)
	go func() {
		defer close(ch)
		limit := make(chan struct{}, 100) // limit 100 pings concurrently
		wg := sync.WaitGroup{}
		for ip := range ipCh {
			limit <- struct{}{}
			wg.Add(1)
			go func(ip net.IP) {
				defer func() { <-limit }()
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				rec := record{Host: ip}
				if hostAlive(ctx, ip) {
					rec.Reachable = true
				}
				ch <- rec
			}(ip)
		}
		wg.Wait()
	}()
	return ch
}

// unamePrefixes takes a channel of net.IP and runs "uname -a" on them via the ssh binary.
func unamePrefixes(user string, recs chan record) chan record {
	ch := make(chan record, 1)
	go func() {
		defer close(ch)

		limit := make(chan struct{}, 100)
		wg := sync.WaitGroup{}
		for rec := range recs {
			if !rec.Reachable {
				ch <- rec
				continue
			}

			limit <- struct{}{}
			wg.Add(1)
			go func(rec record) {
				defer func() { <-limit }()
				defer wg.Done()

				text, err := runUname(context.Background(), rec.Host, user)
				if err != nil {
					ch <- rec
					return
				}
				rec.LoginSSH = true
				rec.Uname = text
				ch <- rec
			}(rec)
		}
		wg.Wait()
	}()
	return ch
}
