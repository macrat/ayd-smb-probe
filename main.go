package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/hirochachacha/go-smb2"
)

func ParseTarget(s string) (*url.URL, error) {
	target, err := url.Parse(s)
	if err != nil {
		return nil, errors.New("invalid target URI")
	}

	if target.Hostname() == "" {
		return nil, errors.New("invalid target URI: hostname is required")
	}

	if pass, _ := target.User.Password(); pass == "" {
		target.User = url.User("guest")
	}

	if target.Port() == "" {
		target.Host = fmt.Sprintf("%s:445", target.Hostname())
	}

	return target, nil
}

func Check(t *url.URL) (latency float64, err error) {
	password, _ := t.User.Password()
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     t.User.Username(),
			Password: password,
		},
	}

	stime := time.Now()

	conn, err := net.Dial("tcp", t.Host)
	if err != nil {
		return float64(time.Now().Sub(stime).Microseconds()) / 1000, err
	}

	s, err := d.Dial(conn)
	if err != nil {
		return float64(time.Now().Sub(stime).Microseconds()) / 1000, err
	}
	s.Logoff()

	return float64(time.Now().Sub(stime).Microseconds()) / 1000, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "$ ayd-smb-alert TARGET_URI")
		os.Exit(2)
	}

	target, err := ParseTarget(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "::status::UNKNOWN")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if latency, err := Check(target); err != nil {
		fmt.Fprintln(os.Stderr, "::status::FAILURE")
		fmt.Fprintf(os.Stderr, "::latency::%.3f\n", latency)
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Fprintln(os.Stderr, "::status::HEALTHY")
		fmt.Fprintf(os.Stderr, "::latency::%.3f\n", latency)
		fmt.Fprintln(os.Stderr, "OK")
	}
}
