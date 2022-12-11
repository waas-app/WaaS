package ip

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/miekg/dns"
	"github.com/waas-app/WaaS/util"
	"go.uber.org/zap"
)

type DNSServer struct {
	server   *dns.Server
	client   *dns.Client
	cache    *bigcache.BigCache
	upstream []string
}

func New(ctx context.Context, upstream []string) (*DNSServer, error) {
	if len(upstream) == 0 {
		upstream = append(upstream, "1.1.1.1")
	}

	localDNSAddr := "0.0.0.0:53"
	util.Logger(ctx).Info("Starting DNS server on", zap.String("address", localDNSAddr), zap.Strings(" with upstream", upstream))

	cache, err := bigcache.New(ctx, bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		util.Logger(ctx).Error("Failed to create cache", zap.Error(err))
		return nil, err
	}

	dnsServer := &DNSServer{
		server: &dns.Server{
			Addr: localDNSAddr,
			Net:  "udp",
		},
		client: &dns.Client{
			SingleInflight: true,
			Timeout:        5 * time.Second,
		},
		cache:    cache,
		upstream: upstream,
	}

	dnsServer.server.Handler = dnsServer

	go func() {
		if err := dnsServer.server.ListenAndServe(); err != nil {
			util.Logger(ctx).Error("Failed to start DNS server", zap.Error(err))
		}
	}()

	return dnsServer, nil
}

func makekey(m *dns.Msg) string {
	q := m.Question[0]
	return fmt.Sprintf("%s:%d:%d", q.Name, q.Qtype, q.Qclass)
}

func serialize(value interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	gob.Register(value)

	err := enc.Encode(&value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func deserialize(valueBytes []byte) (interface{}, error) {
	var value interface{}
	buf := bytes.NewBuffer(valueBytes)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(&value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (d *DNSServer) Lookup(m *dns.Msg) (*dns.Msg, error) {
	key := makekey(m)

	// check the cache first
	if item, found := d.cache.Get(key); found == nil {
		util.Logger(context.Background()).Debug("Found cached dns response", zap.Any("for", m))
		msg, err := deserialize(item)
		if err != nil {
			return nil, err
		}

		acMsg, ok := msg.(*dns.Msg)
		if !ok {
			return nil, fmt.Errorf("failed to cast dns message")
		}
		return acMsg, nil
	}

	// fallback to upstream exchange
	response, _, err := d.client.Exchange(m, net.JoinHostPort(d.upstream[0], "53"))
	if err != nil {
		return nil, err
	}

	if len(response.Answer) > 0 {
		util.Logger(context.Background()).Debug("Caching dns response", zap.Any("for", m))
		b, err := serialize(response)
		if err != nil {
			return nil, err
		}
		d.cache.Set(key, b)
	}

	return response, nil
}

func (d *DNSServer) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	defer func() {
		if err := recover(); err != nil {
			util.Logger(context.Background()).Error("Failed to handle DNS request", zap.Error(err.(error)))
			dns.HandleFailed(w, r)
		}
	}()

	util.Logger(context.Background()).Debug("Received DNS request", zap.String("from", w.RemoteAddr().String()), zap.String("with", r.Question[0].String()))

	switch r.Opcode {
	case dns.OpcodeQuery:
		m, err := d.Lookup(r)
		if err != nil {
			util.Logger(context.Background()).Error("Failed to lookup DNS request", zap.Error(err))
			dns.HandleFailed(w, r)
			return
		}
		m.SetReply(r)
		w.WriteMsg(m)
	default:
		m := &dns.Msg{}
		m.SetReply(r)
		w.WriteMsg(m)
	}
}

func (d *DNSServer) Close() error {
	return d.server.Shutdown()
}
