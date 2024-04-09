package etcd_discovery

import (
	"context"
	"errors"
	"github.com/OpenIMSDK/tools/log"
	"google.golang.org/grpc"
)

func (s *EtcdDiscovery) RegisterConf2Registry(key string, conf []byte) error {
	log.ZWarn(context.Background(), "not implement", errors.New("etcd client not implement RegisterConf2Registry method"))
	return nil
}

func (s *EtcdDiscovery) GetConfFromRegistry(key string) ([]byte, error) {
	log.ZWarn(context.Background(), "not implement", errors.New("etcd client not implement GetConfFromRegistry method"))
	return nil, nil
}

func (s *EtcdDiscovery) CreateRpcRootNodes(serviceNames []string) error {
	log.ZWarn(context.Background(), "not implement", errors.New("etcd client not implement CreateRpcRootNodes method"))
	return nil
}

func (s *EtcdDiscovery) GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) {
	log.ZWarn(ctx, "not implement", errors.New("etcd client not implement GetUserIdHashGatewayHost method"))
	return "", nil
}

func (s *EtcdDiscovery) GetSelfConnTarget() string {
	log.ZWarn(context.Background(), "not implement", errors.New("etcd client not implement GetSelfConnTarget method"))
	return ""
}

func (s *EtcdDiscovery) GetConn(ctx context.Context,
	serviceName string,
	opts ...grpc.DialOption) (*grpc.ClientConn, error) {

	// depends on the strategy, not implement yet
	log.ZWarn(ctx, "not implement", errors.New("etcd client not implement GetConn method"))
	return nil, errors.New("etcd client not implement GetConn method")
}
