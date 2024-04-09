package etcd_discovery

import (
	"context"
	"crypto/tls"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"time"
)

// Client is a wrapper around the etcd client.
type EtcdInterface interface {
	// GetEntries queries the given prefix in etcd and returns a slice
	// containing the values of all keys found, recursively, underneath that
	// prefix.
	GetEntries(prefix string) ([]string, error)

	// WatchPrefix watches the given prefix in etcd for changes. When a change
	// is detected, it will signal on the passed channel. Clients are expected
	// to call GetEntries to update themselves with the latest set of complete
	// values. WatchPrefix will always send an initial sentinel value on the
	// channel after establishing the watch, to ensure that clients always
	// receive the latest set of values. WatchPrefix will block until the
	// context passed to the NewClient constructor is terminated.
	WatchPrefix(prefix string, ch chan struct{})

	// Register a service with etcd.
	Register(key, value string, lease int64) error

	// Deregister a service with etcd.
	Deregister(key string) error

	// LeaseID returns the lease id created for this service instance
	LeaseID() int64
}

type etcdClient struct {
	client *clientv3.Client
	ctx    context.Context

	// watcher context
	wctx context.Context
	// watcher cancel func
	wcf context.CancelFunc

	hbch <-chan *clientv3.LeaseKeepAliveResponse

	// leaseID will be 0 (clientv3.NoLease) if a lease was not created
	leaseID clientv3.LeaseID
}

// ClientOptions defines options for the etcd client. All values are optional.
// If any duration is not specified, a default of 3 seconds will be used.
type ClientOptions struct {
	//Cert          string
	//CACert        string
	//Key           string
	DialTimeout   time.Duration
	DialKeepAlive time.Duration

	// DialOptions is a list of dial options for the gRPC client (e.g., for interceptors).
	// For example, pass grpc.WithBlock() to block until the underlying connection is up.
	// Without this, Dial returns immediately and connecting the server happens in background.
	DialOptions []grpc.DialOption

	Username string
	Password string
}

// NewClient returns Client with a connection to the named machines. It will
// return an error if a connection to the cluster cannot be made.
func NewClient(ctx context.Context, endpoints *[]string, options ClientOptions) (EtcdInterface, error) {
	if options.DialTimeout == 0 {
		options.DialTimeout = 3 * time.Second
	}
	if options.DialKeepAlive == 0 {
		options.DialKeepAlive = 3 * time.Second
	}

	var err error
	var tlscfg *tls.Config

	//if options.Cert != "" && options.Key != "" {
	//	tlsInfo := transport.TLSInfo{
	//		CertFile:      options.Cert,
	//		KeyFile:       options.Key,
	//		TrustedCAFile: options.CACert,
	//	}
	//	tlscfg, err = tlsInfo.ClientConfig()
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	cli, err := clientv3.New(clientv3.Config{
		Context:           ctx,
		Endpoints:         *endpoints,
		DialTimeout:       options.DialTimeout,
		DialKeepAliveTime: options.DialKeepAlive,
		DialOptions:       options.DialOptions,
		TLS:               tlscfg,
		Username:          options.Username,
		Password:          options.Password,
	})
	if err != nil {
		return nil, err
	}

	return &etcdClient{
		client: cli,
		ctx:    ctx,
	}, nil
}

func (c *etcdClient) LeaseID() int64 { return int64(c.leaseID) }

func (c *etcdClient) GetEntries(prefix string) ([]string, error) {
	resp, err := c.client.Get(c.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	entries := make([]string, len(resp.Kvs))
	for i, kv := range resp.Kvs {
		entries[i] = string(kv.Value)
	}

	return entries, nil
}

// WatchPrefix implements the etcd Client interface.
func (c *etcdClient) WatchPrefix(prefix string, ch chan struct{}) {
	c.wctx, c.wcf = context.WithCancel(c.ctx)

	wch := c.client.Watch(c.wctx, prefix, clientv3.WithPrefix(), clientv3.WithRev(0))
	ch <- struct{}{}
	for wr := range wch {
		if wr.Canceled {
			return
		}
		ch <- struct{}{}
	}
}

func (c *etcdClient) Register(key, value string, lease int64) error {
	var err error
	if key == "" || value == "" {
		return errors.New("error: no key or value")
	}

	grantResp, err := c.client.Grant(c.ctx, lease)
	if err != nil {
		return err
	}
	c.leaseID = grantResp.ID

	_, err = c.client.Put(
		c.ctx,
		key,
		value,
		clientv3.WithLease(c.leaseID),
	)
	if err != nil {
		return err
	}

	// this will keep the key alive 'forever' or until we revoke it or
	// the context is canceled
	c.hbch, err = c.client.KeepAlive(c.ctx, c.leaseID)

	if err != nil {
		return err
	}

	// discard the keepalive response, make etcd library not to complain
	// fix bug #799
	go func() {
		for {
			select {
			case r := <-c.hbch:
				// avoid dead loop when channel was closed
				if r == nil {
					return
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (c *etcdClient) Deregister(key string) error {
	defer c.Close()

	if key == "" {
		return errors.New("error: no key")
	}
	if _, err := c.client.Delete(c.ctx, key, clientv3.WithIgnoreLease()); err != nil {
		return err
	}

	return nil
}

func (c *etcdClient) Close() {
	if err := c.client.Close(); err != nil {

	}
}
