package axslogparser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var part = `"(?P<%s>(?:[^"]*(?:\\")?)*)"`

var logRe = regexp.MustCompile(
	`(?:(?P<vhost>\S+)\s)?` + // %v(The canonical ServerName/virtual host)
		`(?P<remote_addr>\S+)\s` + // %h(Remote Hostname) $remote_addr
		`\S+\s+` + // %l(Remote Logname)
		`(?P<remote_user>\S+)\s+` + // $remote_user
		`\[(?P<time_local>[^]]+)\]\s` + // $time_local
		fmt.Sprintf(part, "request") + `\s` + // $request
		`(?P<status>[0-9]{3})\s` + // $status
		`(?P<body_bytes_sent>[0-9]+)` + // $body_bytes_sent
		`(?:\s` + // combined option start
		fmt.Sprintf(part, "http_referer") + `\s` + // $http_referer
		fmt.Sprintf(part, "http_user_agent") + // $http_user_agent
		`)?`) // combined option end

type Log struct {
	VirtualHost string
	Host        string
	User        string
	Time        time.Time
	Request     string
	Status      int
	Size        uint64
	Referer     string
	UA          string
}

func Parse(line string, loc *time.Location) (*Log, error) {
	matches := logRe.FindStringSubmatch(line)
	if len(matches) < 1 {
		return nil, fmt.Errorf("not matched")
	}
	l := &Log{}
	for i, name := range logRe.SubexpNames() {
		switch name {
		case "vhost":
			l.VirtualHost = matches[i]
		case "remote_addr":
			l.Host = matches[i]
		case "remote_user":
			l.User = matches[i]
		case "time_local":
			const timeLayout = "02/Jan/2006:15:04:05 -0700"
			l.Time, _ = time.ParseInLocation(timeLayout, matches[i], loc)
		case "request":
			l.Request = dequote(matches[i])
		case "status":
			l.Status, _ = strconv.Atoi(matches[i])
		case "body_bytes_sent":
			l.Size, _ = strconv.ParseUint(matches[i], 10, 64)
		case "http_referer":
			l.Referer = dequote(matches[i])
		case "http_user_agent":
			l.UA = dequote(matches[i])
		}
	}
	// TODO parse request into method, path and proto
	return l, nil
}

func dequote(item string) string {
	// TOOD
	return item
}
