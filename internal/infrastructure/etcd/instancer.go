package etcd_discovery

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"sync"
)

// Instancer yields instances stored in a certain etcd keyspace. Any kind of
// change in that keyspace is watched and will update the Instancer's Instancers.
type Instancer struct {
	ctx        context.Context
	client     EtcdInterface
	prefix     string
	options    []grpc.DialOption
	logger     log.Logger
	quitc      chan struct{}
	localConns map[string][]*grpc.ClientConn
	lock       sync.Locker
}

// NewInstancer returns an etcd instancer. It will start watching the given
// prefix for changes, and update the subscribers.
func NewInstancer(ctx context.Context, c EtcdInterface, prefix string, opts ...grpc.DialOption) *Instancer {
	s := &Instancer{
		ctx:     ctx,
		options: opts,
		client:  c,
		prefix:  prefix,
		quitc:   make(chan struct{}),
	}
	return s
}

func (s *Instancer) loop() {
	ch := make(chan struct{})
	go s.client.WatchPrefix(s.prefix, ch)

	for {
		select {
		case <-ch:
			instances, err := s.client.GetEntries(s.prefix)
			if err != nil {
				s.logger.Printf("failed to retrieve entries")
				continue
			}
			if err = s.doConnect(instances); err != nil {
				s.logger.Printf("get connection failed %v", err)
			}
		case <-s.quitc:
			return
		}
	}
}

// Stop terminates the Instancer.
func (s *Instancer) Stop() {
	close(s.quitc)
}

func (s *Instancer) initialize() error {
	conns := s.localConns[s.prefix]
	if len(conns) == 0 {
		entries, err := s.client.GetEntries(s.prefix)
		if err != nil {
			return fmt.Errorf("get register entries failed for service %s from etcd", s.prefix)
		}
		if len(entries) == 0 {
			return fmt.Errorf("get no register entries for service %s from etcd", s.prefix)
		}
		if err := s.doConnect(entries); err != nil {
			return err
		}
	}
	go s.loop()
	return nil
}

func (s *Instancer) doConnect(entries []string) error {
	if len(entries) == 0 {
		// no need to update
		return nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	var conns []*grpc.ClientConn
	for _, addr := range entries {
		cc, err := grpc.DialContext(s.ctx, addr, s.options...)
		if err != nil {
			//log.ZError(context.Background(), "dialContext failed", err, "addr", addr, "opts", s.options)
			return err
		}
		conns = append(conns, cc)
	}
	delete(s.localConns, s.prefix)
	s.localConns[s.prefix] = conns
	return nil
}

func (s *Instancer) GetLocalConns() map[string][]*grpc.ClientConn {
	return s.localConns
}

func (s *Instancer) GetLocalConnsByKey(prefix string) ([]*grpc.ClientConn, error) {
	return s.localConns[prefix], nil
}

func (s *Instancer) AddOptions(opts ...grpc.DialOption) {
	s.options = append(s.options, opts...)
}
