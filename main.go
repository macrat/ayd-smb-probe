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

func Check(t *url.URL) (stime time.Time, latency time.Duration, err error) {
	password, _ := t.User.Password()
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     t.User.Username(),
			Password: password,
		},
	}

	stime = time.Now()

	conn, err := net.Dial("tcp", t.Host)
	if err != nil {
		return stime, time.Now().Sub(stime), err
	}

	s, err := d.Dial(conn)
	if err != nil {
		return stime, time.Now().Sub(stime), err
	}
	s.Logoff()

	return stime, time.Now().Sub(stime), nil
}

func PrintLog(stime time.Time, status string, latency time.Duration, message string) {
	fmt.Printf("%s\t%s\t%.3f\t%s\t%s", stime.Format(time.RFC3339), status, float64(latency.Microseconds())/1000, os.Args[1], message)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "$ ayd-smb-alert TARGET_URI")
		os.Exit(2)
	}

	target, err := ParseTarget(os.Args[1])
	if err != nil {
		PrintLog(time.Now(), "UNKNOWN", 0, "invalid URL format: "+err.Error())
		return
	}

	if stime, latency, err := Check(target); err != nil {
		PrintLog(stime, "FAILURE", latency, err.Error())
	} else {
		PrintLog(stime, "HEALTHY", latency, "OK")
	}
}
