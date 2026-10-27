package transport

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	Broker   string
	ClientID string
	Username string
	Password string
	TLS      *tls.Config
	KeepAlive time.Duration
	Timeout  time.Duration
}

type Client struct {
	mu     sync.RWMutex
	c      mqtt.Client
	opts   *mqtt.ClientOptions
	state  string
}

func New(cfg Config) (*Client, error) {
	if cfg.Broker == "" {
		return nil, errors.New("broker required")
	}
	if cfg.ClientID == "" {
		return nil, errors.New("clientid required")
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
	if cfg.KeepAlive == 0 {
		cfg.KeepAlive = 30
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientID)
	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}
	if cfg.TLS != nil {
		opts.SetTLSConfig(cfg.TLS)
	}
	opts.SetAutoReconnect(false)
	opts.SetConnectTimeout(cfg.Timeout)
	opts.SetKeepAlive(cfg.KeepAlive)

	c := mqtt.NewClient(opts)
	return &Client{c: c, opts: opts, state: "disconnected"}, nil
}

func (cl *Client) Connect(ctx context.Context) error {
	if cl == nil || cl.c == nil {
		return errors.New("nil client")
	}
	t := tokenWithContext(ctx)
	cl.c.Connect().Wait()
	if !t.WaitTimeout(cl.timeout()) {
		return errors.New("connect timeout")
	}
	if err := t.Error(); err != nil {
		return err
	}
	cl.mu.Lock()
	cl.state = "connected"
	cl.mu.Unlock()
	return nil
}

func (cl *Client) Disconnect(quiesce time.Duration) {
	if cl == nil || cl.c == nil {
		return
	}
	cl.c.Disconnect(uint(quiesce.Milliseconds()))
	cl.mu.Lock()
	cl.state = "disconnected"
	cl.mu.Unlock()
}

func (cl *Client) Publish(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error {
	if cl == nil || cl.c == nil {
		return errors.New("nil client")
	}
	if topic == "" {
		return errors.New("topic required")
	}
	t := tokenWithContext(ctx)
	cl.c.Publish(topic, qos, retained, payload).Wait()
	if !t.WaitTimeout(cl.timeout()) {
		return errors.New("publish timeout")
	}
	return t.Error()
}

func (cl *Client) Subscribe(ctx context.Context, topic string, qos byte, callback mqtt.MessageHandler) error {
	if cl == nil || cl.c == nil {
		return errors.New("nil client")
	}
	if topic == "" {
		return errors.New("topic required")
	}
	t := tokenWithContext(ctx)
	cl.c.Subscribe(topic, qos, callback).Wait()
	if !t.WaitTimeout(cl.timeout()) {
		return errors.New("subscribe timeout")
	}
	return t.Error()
}

func (cl *Client) Unsubscribe(ctx context.Context, topics ...string) error {
	if cl == nil || cl.c == nil {
		return errors.New("nil client")
	}
	t := tokenWithContext(ctx)
	cl.c.Unsubscribe(topics...).Wait()
	if !t.WaitTimeout(cl.timeout()) {
		return errors.New("unsubscribe timeout")
	}
	return t.Error()
}

func (cl *Client) IsConnected() bool {
	if cl == nil || cl.c == nil {
		return false
	}
	return cl.c.IsConnected()
}

func (cl *Client) State() string {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return cl.state
}

func (cl *Client) timeout() time.Duration {
	if cl.opts == nil {
		return 10 * time.Second
	}
	// use opts timeout if set
	return 10 * time.Second
}

type token interface {
	Wait() bool
	WaitTimeout(time.Duration) bool
	Error() error
}

func tokenWithContext(ctx context.Context) token {
	// simple wrapper
	type t struct{ mqtt.Token }
	return t{}
}

func NewTLSConfig(caCert []byte, clientCert tls.Certificate, insecureSkipVerify bool) (*tls.Config, error) {
	if insecureSkipVerify {
		return nil, errors.New("insecure skip verify not allowed")
	}
	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if len(caCert) > 0 {
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to load ca cert")
		}
		cfg.RootCAs = pool
	}
	if len(clientCert.Certificate) > 0 {
		cfg.Certificates = []tls.Certificate{clientCert}
	}
	return cfg, nil
}
