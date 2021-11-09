package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/macrat/ayd/lib-ayd"
)

var (
	version = "HEAD"
	commit  = "UNKNOWN"
)

func NormalizeTarget(u *url.URL) *url.URL {
	u = &url.URL{
		Scheme: "smb",
		Host:   u.Host,
		User:   u.User,
	}

	if pass, _ := u.User.Password(); pass == "" {
		u.User = url.User("guest")
	}

	return u
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

	host := t.Host
	if t.Port() == "" {
		host += ":445"
	}

	conn, err := net.Dial("tcp", host)
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

func main() {
	flag.Usage = func() {
		fmt.Println("SMB protocol plugin for Ayd?")
		fmt.Println()
		fmt.Println("usage: ayd-smb-probe TARGET_URL")
	}
	showVersion := flag.Bool("v", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("ayd-smb-probe %s (%s)\n", version, commit)
		return
	}

	args, err := ayd.ParseProbePluginArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "$ ayd-smb-probe TARGET_URL")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	args.TargetURL = NormalizeTarget(args.TargetURL)
	logger := ayd.NewLogger(args.TargetURL)

	if args.TargetURL.Hostname() == "" {
		logger.Failure("invalid target URI: hostname is required")
		os.Exit(2)
	}

	if stime, latency, err := Check(args.TargetURL); err != nil {
		logger.WithTime(stime, latency).Failure(err.Error())
	} else {
		logger.WithTime(stime, latency).Healthy("OK")
	}
}
