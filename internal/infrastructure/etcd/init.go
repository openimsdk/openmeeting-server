package etcd_discovery

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"openmeeting-server/pkg/common/config"
)

func InitEtcdService(ctx context.Context, serverConfig *config.Config, registerIP string, rpcPort int) error {
	client, cerr := clientv3.NewFromURL((*serverConfig.EtcdConfig.Address)[0])
	if cerr != nil {
		panic(cerr)
	}

	lease, _ := client.Grant(ctx, 10)
	em, err := endpoints.NewManager(client, "openmeeting/rtc-service")
	if err != nil {
		panic(err)
	}
	grpcAddress := fmt.Sprintf("%s:%d", registerIP, rpcPort)
	err = em.AddEndpoint(ctx, "openmeeting/rtc-service/g1", endpoints.Endpoint{Addr: grpcAddress}, clientv3.WithLease(lease.ID))
	if err != nil {
		panic(err)
	}

	ch, err := client.KeepAlive(ctx, lease.ID)

	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case r := <-ch:
				// avoid dead loop when channel was closed
				if r == nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}
