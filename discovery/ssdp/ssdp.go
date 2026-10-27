package ssdp

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/sidalikharef/goremote/discovery"
)

// Discoverer implements SSDP discovery.
type Discoverer struct{}

func New() *Discoverer { return &Discoverer{} }

func (d *Discoverer) Discover(ctx context.Context, opts discovery.Options) ([]discovery.Result, error) {
	if opts.Timeout <= 0 {
		opts.Timeout = 3 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	results := map[string]discovery.Result{}
	addr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		return []discovery.Result{}, nil
	}
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return []discovery.Result{}, nil
	}
	defer conn.Close()

	deadline := time.Now().Add(opts.Timeout)
	conn.SetReadDeadline(deadline)

	msg := "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 2\r\nST: ssdp:all\r\n\r\n"
	_, _ = conn.WriteToUDP([]byte(msg), addr)

	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			goto done
		default:
		}
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				break
			}
			break
		}
		resp := string(buf[:n])
		if !strings.Contains(resp, "HTTP/") {
			continue
		}
		var loc, usn, st string
		for _, line := range strings.Split(resp, "\r\n") {
			kv := strings.SplitN(line, ":", 2)
			if len(kv) != 2 {
				continue
			}
			k := strings.TrimSpace(strings.ToLower(kv[0]))
			v := strings.TrimSpace(kv[1])
			if k == "location" {
				loc = v
			} else if k == "usn" {
				usn = v
			} else if k == "st" {
				st = v
			}
		}
		if loc == "" {
			continue
		}
		host := extractHost(loc)
		if host == "" {
			continue
		}
		id := usn
		if id == "" {
			id = host
		}
		r := discovery.Result{
			ID:         id,
			Name:       st,
			Address:    host,
			Location:   loc,
			USN:        usn,
			Source:     "ssdp",
			DriverHint: guessDriver(st, usn),
		}
		if discovery.Validate(r) {
			results[id] = r
		}
	}
done:
	out := make([]discovery.Result, 0, len(results))
	for _, v := range results {
		out = append(out, v)
		if opts.MaxResults > 0 && len(out) >= opts.MaxResults {
			break
		}
	}
	return out, nil
}

func extractHost(loc string) string {
	loc = strings.TrimSpace(loc)
	if strings.HasPrefix(loc, "http://") {
		loc = strings.TrimPrefix(loc, "http://")
	} else if strings.HasPrefix(loc, "https://") {
		loc = strings.TrimPrefix(loc, "https://")
	}
	if i := strings.IndexAny(loc, "/:?"); i >= 0 {
		loc = loc[:i]
	}
	if h, _, err := net.SplitHostPort(loc); err == nil {
		return h
	}
	return loc
}

func guessDriver(st, usn string) string {
	s := strings.ToLower(st + " " + usn)
	if strings.Contains(s, "roku") {
		return "roku"
	}
	if strings.Contains(s, "samsung") || strings.Contains(s, "tizen") {
		return "samsung"
	}
	if strings.Contains(s, "lge") || strings.Contains(s, "webos") {
		return "webos"
	}
	if strings.Contains(s, "android") || strings.Contains(s, "google") {
		return "androidtv"
	}
	return ""
}

func (d *Discoverer) DiscoverRoku(ctx context.Context, opts discovery.Options) ([]discovery.Result, error) {
	if opts.Timeout <= 0 {
		opts.Timeout = 3 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	results := map[string]discovery.Result{}
	addr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		return []discovery.Result{}, nil
	}
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return []discovery.Result{}, nil
	}
	defer conn.Close()

	deadline := time.Now().Add(opts.Timeout)
	conn.SetReadDeadline(deadline)

	msg := "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 2\r\nST: roku:ecp\r\n\r\n"
	_, _ = conn.WriteToUDP([]byte(msg), addr)

	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			goto done
		default:
		}
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				break
			}
			break
		}
		resp := string(buf[:n])
		if !strings.Contains(resp, "HTTP/") {
			continue
		}
		var loc, usn, st string
		for _, line := range strings.Split(resp, "\r\n") {
			kv := strings.SplitN(line, ":", 2)
			if len(kv) != 2 {
				continue
			}
			k := strings.TrimSpace(strings.ToLower(kv[0]))
			v := strings.TrimSpace(kv[1])
			if k == "location" {
				loc = v
			} else if k == "usn" {
				usn = v
			} else if k == "st" {
				st = v
			}
		}
		if loc == "" {
			continue
		}
		host := extractHost(loc)
		if host == "" || net.ParseIP(host) == nil {
			continue
		}
		id := usn
		if id == "" {
			id = "roku-" + host
		}
		r := discovery.Result{
			ID:         id,
			Name:       st,
			Address:    host,
			Location:   loc,
			USN:        usn,
			Source:     "ssdp",
			DriverHint: "roku",
		}
		if discovery.Validate(r) {
			results[id] = r
		}
	}
done:
	out := make([]discovery.Result, 0, len(results))
	for _, v := range results {
		out = append(out, v)
		if opts.MaxResults > 0 && len(out) >= opts.MaxResults {
			break
		}
	}
	return out, nil
}
