package startrpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/tools/mw"
	"github.com/OpenIMSDK/tools/network"
	"github.com/OpenIMSDK/tools/utils"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	config "openmeeting-server/dto"
	etcd_discovery "openmeeting-server/internal/infrastructure/etcd"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func Start(rpcPort int,
	rpcRegisterName string,
	rpcFn func(server *grpc.Server) error,
	options ...grpc.ServerOption) error {

	fmt.Printf("start %s server, port: %d\n", rpcRegisterName, rpcPort)
	listener, err := net.Listen(
		"tcp",
		net.JoinHostPort(network.GetListenIP(config.Config.RPC.ListenIP), strconv.Itoa(rpcPort)),
	)
	if err != nil {
		return err
	}

	defer listener.Close()
	registerIP, err := network.GetRpcRegisterIP(config.Config.RPC.RegisterIP)
	if err != nil {
		return err
	}

	client, err := etcd_discovery.NewEtcdDiscovery(
		context.Background(), rpcRegisterName, registerIP, rpcPort, &config.Config.Etcd)
	if err != nil {
		return utils.Wrap1(err)
	}
	defer client.Close()

	//client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	options = append(options, mw.GrpcServer())
	srv := grpc.NewServer(options...)
	once := sync.Once{}
	defer func() {
		once.Do(srv.GracefulStop)
	}()

	err = rpcFn(srv)
	if err != nil {
		return utils.Wrap1(err)
	}
	err = client.Register(
		rpcRegisterName,
		registerIP,
		rpcPort,
	)
	if err != nil {
		return utils.Wrap1(err)
	}

	var wg errgroup.Group

	wg.Go(func() error {
		return utils.Wrap1(srv.Serve(listener))
	})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigs

	var (
		done = make(chan struct{}, 1)
		gerr error
	)

	go func() {
		once.Do(srv.GracefulStop)
		gerr = wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return gerr

	case <-time.After(15 * time.Second):
		return utils.Wrap1(errors.New("timeout exit"))
	}

}
