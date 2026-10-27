package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sidalikharef/goremote/config"
	"github.com/sidalikharef/goremote/drivers/registry"
	"github.com/sidalikharef/goremote/drivers/roku"
	vidaa "github.com/sidalikharef/goremote/drivers/vidaa"
	"github.com/sidalikharef/goremote/httpapi"
	"github.com/sidalikharef/goremote/internal/auth"
	"github.com/sidalikharef/goremote/internal/validate"
)

func main() {
	var (
		tvIP     = flag.String("tv", os.Getenv("GOREMOTE_TV_IP"), "IP of TV")
		driver   = flag.String("driver", getEnv("GOREMOTE_DRIVER", "roku"), "Driver")
		addr     = flag.String("addr", getEnv("GOREMOTE_ADDR", "127.0.0.1:8080"), "Listen address")
		timeout  = flag.String("timeout", getEnv("GOREMOTE_TIMEOUT", "5s"), "Timeout")
		confPath = flag.String("config", config.DefaultPath(), "Config file path")
	)
	flag.Parse()

	// Allow zero-TV startup for multi-device support

	dur, err := time.ParseDuration(*timeout)
	if err != nil {
		dur = 5 * time.Second
	}

	reg := registry.New()
	reg.RegisterFactory("vidaa", vidaa.Factory)
	if cfg, err := config.Load(*confPath); err == nil {
		for _, d := range cfg.Devices {
			if _, err := reg.RegisterDevice(d.ID, d.Name, d.Driver, d.Address); err == nil {
				if d.Driver == "roku" && d.Address != "" {
					if c, err := roku.New(d.Address, dur); err == nil {
						if dev, ok := reg.GetDevice(d.ID); ok {
							dev.SetDriver(c)
						}
					}
				}
			}
		}
	}

	if *tvIP != "" {
		if validate.IP(*tvIP) == nil && *driver == "roku" {
			if c, err := roku.New(*tvIP, dur); err == nil {
				if _, ok := reg.GetDevice(*tvIP); !ok {
					if d, err := reg.RegisterDevice(*tvIP, "Roku", "roku", *tvIP); err == nil {
						d.SetDriver(c)
					}
				} else if d, ok := reg.GetDevice(*tvIP); ok {
					d.SetDriver(c)
				}
			}
		}
	}

	authStore := auth.New(*confPath + ".auth.json")
	srv := httpapi.NewWithAuth(reg, authStore)
	h := httpapi.SecurityHeaders(srv.Routes())
	h = httpapi.MaxBody(1024 * 64)(h)
	allowedHosts := map[string]bool{}
	allowedHosts["localhost:8080"] = true
	allowedHosts["127.0.0.1:8080"] = true
	allowedHosts["[::1]:8080"] = true
	allowedHosts["localhost"] = true
	allowedHosts["127.0.0.1"] = true
	allowedHosts["[::1]"] = true
	if *addr != "" {
		allowedHosts[*addr] = true
		if strings.HasPrefix(*addr, ":") {
			allowedHosts["127.0.0.1"+*addr] = true
			allowedHosts["localhost"+*addr] = true
		}
	}
	h = httpapi.CheckHostSafe(allowedHosts)(h)
	h = httpapi.CheckOriginSafe(allowedHosts)(h)
	h = httpapi.RequireAuth(authStore)(h)

	server := &http.Server{
		Addr:         *addr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	fmt.Printf("GoRemote listening on %s (LAN binding requires explicit config)\n", *addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// experimental VIDAA disabled by default
