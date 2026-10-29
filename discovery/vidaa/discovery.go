package vidaa

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sidalikharef/goremote/discovery"
)

// Provider implements evidence-backed VIDAA discovery hints.
type Provider struct {
	Timeout time.Duration
}

func New() *Provider { return &Provider{Timeout: 3 * time.Second} }

func (p *Provider) Name() string { return "vidaa" }

func (p *Provider) Discover(ctx context.Context, opts discovery.Options) ([]discovery.Result, error) {
	if opts.Timeout > 0 {
		p.Timeout = opts.Timeout
	}
	// Try common VIDAA advertisement patterns defensively
	results := map[string]discovery.Result{}
	// SSDP-like search for vidaa/upnp patterns
	results = mergeResults(results, p.searchSSDP(ctx))
	// Also attempt to detect common VIDAA service ports by probing local network info not done
	// but we also add nothing automatic; this is passive
	out := make([]discovery.Result, 0, len(results))
	for _, v := range results {
		out = append(out, v)
		if opts.MaxResults > 0 && len(out) >= opts.MaxResults {
			break
		}
	}
	return out, nil
}

func (p *Provider) searchSSDP(ctx context.Context) map[string]discovery.Result {
	results := map[string]discovery.Result{}
	timeout := p.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	addr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		return results
	}
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return results
	}
	defer conn.Close()
	deadline := time.Now().Add(timeout)
	_ = conn.SetReadDeadline(deadline)
	// Try common STs for VIDAA/UPnP variants
	msgs := []string{
		"M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 2\r\nST: urn:schemas-upnp-org:device:MediaRenderer:1\r\n\r\n",
		"M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 2\r\nST: upnp:rootdevice\r\n\r\n",
		"M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 2\r\nST: ssdp:all\r\n\r\n",
	}
	for _, m := range msgs {
		select {
		case <-c.Done():
			return results
		default:
		}
		_, _ = conn.WriteToUDP([]byte(m), addr)
	}
	buf := make([]byte, 8192)
	for {
		select {
		case <-c.Done():
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
		var loc, usn, st, server string
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
			} else if k == "server" {
				server = v
			}
		}
		if loc == "" {
			continue
		}
		host := extractHost(loc)
		if host == "" || net.ParseIP(host) == nil {
			continue
		}
		s := strings.ToLower(st + " " + usn + " " + server)
		if !strings.Contains(s, "hisense") && !strings.Contains(s, "vidaa") && !strings.Contains(s, "hitachi") {
			// Be permissive but mark as probable only with strong hints; skip if clearly other brands
			// But to be safe, don't auto-classify unknown as VIDAA
			continue
		}
		id := usn
		if id == "" {
			id = "vidaa-" + host
		}
		name := "Hisense VIDAA"
		if st != "" {
			name = st
		}
		r := discovery.Result{
			ID:               id,
			Name:             friendlyName(name, usn, host),
			Address:          host,
			Location:         loc,
			USN:              usn,
			Source:           "ssdp",
			DiscoverySource:  "ssdp",
			DriverHint:       "vidaa",
			Driver:           "vidaa",
			Manufacturer:     "Hisense",
			Confidence:       "probable",
			Experimental:     true,
			RequiresPairing:  true,
		}
		if discovery.Validate(r) {
			results[id] = r
		}
	}
done:
	return results
}

func mergeResults(a, b map[string]discovery.Result) map[string]discovery.Result {
	for k, v := range b {
		a[k] = v
	}
	return a
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

func friendlyName(name, usn, host string) string {
	s := strings.TrimSpace(name)
	if s == "" {
		s = strings.TrimSpace(usn)
	}
	if s == "" {
		return "Unidentified TV candidate (" + host + ")"
	}
	if strings.Contains(s, "uuid:") && len(s) > 40 {
		return "Possible VIDAA device (" + host + ")"
	}
	return s
}

// ProbeCandidate can do a safe non-invasive probe to enrich metadata if needed
func ProbeCandidate(addr string) map[string]interface{} {
	res := map[string]interface{}{}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // read-only metadata only, not for auth bypass
	}
	client := &http.Client{Transport: tr, Timeout: 2 * time.Second}
	// Try common endpoints defensively
	for _, path := range []string{"/", "/device.xml", "/description.xml"} {
		req, _ := http.NewRequest("GET", "https://"+addr+":36669"+path, nil)
		if req == nil {
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		res["probe_path"] = path
		res["status"] = resp.StatusCode
		resp.Body.Close()
		break
	}
	// Also note open ports observed
	_, _ = json.Marshal(res)
	return res
}
