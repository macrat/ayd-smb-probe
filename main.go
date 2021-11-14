package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/macrat/ayd/lib-ayd"
)

var (
	version = "HEAD"
	commit  = "UNKNOWN"
)

func NormalizeTarget(u *url.URL) *url.URL {
	if u.Path == "" {
		u.Path = "/"
	}

	u = &url.URL{
		Scheme: "smb",
		Host:   u.Host,
		User:   u.User,
		Path:   path.Clean(u.Path),
	}

	if pass, _ := u.User.Password(); pass == "" {
		u.User = url.User("guest")
	}

	return u
}

func SplitPath(urlPath string) (share, filePath string) {
	ss := strings.SplitN(urlPath[1:], "/", 2)
	if len(ss) == 1 {
		return ss[0], "."
	}
	return ss[0], path.Clean(ss[1])
}

func Check(t *url.URL) (msg string, stime time.Time, latency time.Duration, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

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

	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", host)
	if err != nil {
		return "", stime, time.Since(stime), err
	}

	sess, err := d.DialContext(ctx, conn)
	if err != nil {
		return "", stime, time.Since(stime), err
	}
	defer sess.Logoff()

	sess = sess.WithContext(ctx)

	shareName, path := SplitPath(t.Path)
	if shareName == "" {
		shares, err := sess.ListSharenames()
		if err != nil {
			return "", stime, time.Since(stime), err
		}
		return fmt.Sprintf("type=server shares=%d", len(shares)), stime, time.Since(stime), nil
	}

	share, err := sess.Mount(shareName)
	if err != nil {
		return "", stime, time.Since(stime), err
	}
	defer share.Umount()

	share = share.WithContext(ctx)

	stat, err := share.Stat(path)
	if err != nil {
		return "", stime, time.Since(stime), err
	}

	if stat.IsDir() {
		files, err := share.ReadDir(path)
		if err != nil {
			return "", stime, time.Since(stime), err
		}
		msg = fmt.Sprintf("type=directory files=%d", len(files))
	} else {
		msg = fmt.Sprintf("type=file size=%d", stat.Size())
	}
	return msg, stime, time.Since(stime), nil
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
		logger.Failure("invalid URL: hostname is required")
		return
	}

	if msg, stime, latency, err := Check(args.TargetURL); err != nil {
		logger.WithTime(stime, latency).Failure(err.Error())
	} else {
		logger.WithTime(stime, latency).Healthy(msg)
	}
}
