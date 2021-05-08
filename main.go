package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/macrat/ayd/lib-ayd"
)

func NormalizeTarget(u *url.URL) error {
	if u.Hostname() == "" {
		return errors.New("invalid target URI: hostname is required")
	}

	if pass, _ := u.User.Password(); pass == "" {
		u.User = url.User("guest")
	}

	if u.Port() == "" {
		u.Host = u.Hostname() + ":445"
	}

	return nil
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

func main() {
	args, err := ayd.ParseProbePluginArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "$ ayd-smb-alert TARGET_URI")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	logger := ayd.NewLogger(args.TargetURL)

	if err = NormalizeTarget(args.TargetURL); err != nil {
		logger.Failure("invalid URL format: " + err.Error())
		return
	}

	if stime, latency, err := Check(args.TargetURL); err != nil {
		logger.WithTime(stime, latency).Failure(err.Error())
	} else {
		logger.WithTime(stime, latency).Healthy("OK")
	}
}
