package etcd_discovery

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"openmeeting-server/pkg/common/config"
	"time"
)

type EtcdDiscovery struct {
	registrar *Registrar
	instancer *Instancer

	serviceName string
}

func NewEtcdDiscovery(ctx context.Context,
	serviceName, ip string, port int, conf *config.EtcdConf, opts ...grpc.DialOption) (*EtcdDiscovery, error) {
	serviceAddress := fmt.Sprintf("%s:%d", ip, port)
	client, err := NewClient(ctx, conf.Address, ClientOptions{
		DialTimeout:   time.Duration(*conf.Timeout) * time.Second,
		DialKeepAlive: time.Duration(*conf.Ttl) * time.Second,
	})
	if err != nil {
		return nil, nil
	}

	return &EtcdDiscovery{
		registrar: NewRegistrar(client, serviceName, serviceAddress, int64(*conf.Lease)),
		instancer: NewInstancer(ctx, client, serviceName, opts...),
	}, nil
}

func (s *EtcdDiscovery) GetConns(ctx context.Context,
	serviceName string,
	opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {

	return s.instancer.GetLocalConnsByKey(serviceName)
}

func (s *EtcdDiscovery) CloseConn(conn *grpc.ClientConn) {
	if err := conn.Close(); err != nil {
		log.Printf("close conn failed")
	}
}

// Close etcd connection
func (s *EtcdDiscovery) Close() {

}

// client dial to server options
func (s *EtcdDiscovery) AddOption(opts ...grpc.DialOption) {
	s.instancer.AddOptions(opts...)
}

func (s *EtcdDiscovery) GetClientLocalConns() map[string][]*grpc.ClientConn {
	return s.instancer.GetLocalConns()
}

func (s *EtcdDiscovery) Register(rpcRegisterName, host string, port int, opts ...grpc.DialOption) error {
	s.registrar.Register()
	return nil
}

func (s *EtcdDiscovery) UnRegister() error {
	s.registrar.Deregister()
	return nil
}
