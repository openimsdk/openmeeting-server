package startrpc

import (
	"context"
	"fmt"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/network"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net"
	"net/http"
	etcd_discovery "openmeeting-server/internal/infrastructure/etcd"
	subConfig "openmeeting-server/pkg/common/config"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

//func Start1(rpcPort int,
//	rpcRegisterName string,
//	rpcFn func(server *grpc.Server) error,
//	options ...grpc.ServerOption) error {
//
//	fmt.Printf("start %s server, port: %d\n", rpcRegisterName, rpcPort)
//	listener, err := net.Listen(
//		"tcp",
//		net.JoinHostPort(network.GetListenIP(subConfig.Config.RPC.ListenIP), strconv.Itoa(rpcPort)),
//	)
//	if err != nil {
//		return err
//	}
//
//	defer listener.Close()
//	registerIP, err := network.GetRpcRegisterIP(subConfig.Config.RPC.RegisterIP)
//	if err != nil {
//		return err
//	}
//
//	//client, err := etcd_discovery.NewEtcdDiscovery(
//	//	context.Background(), rpcRegisterName, registerIP, rpcPort, &config.Config.Etcd)
//	//if err != nil {
//	//	return utils.Wrap1(err)
//	//}
//	//defer client.Close()
//
//	//client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
//
//	options = append(options, mw.GrpcServer())
//	srv := grpc.NewServer(options...)
//	once := sync.Once{}
//	defer func() {
//		once.Do(srv.GracefulStop)
//	}()
//	err = rpcFn(srv)
//	if err != nil {
//		return errs.Wrap(err)
//	}
//
//	cli, cerr := clientv3.NewFromURL("http://localhost:2379")
//	if cerr != nil {
//		panic(cerr)
//	}
//
//	ctx := context.Background()
//	lease, _ := cli.Grant(ctx, 10)
//	em, err := endpoints.NewManager(cli, "openmeeting/rtc-service")
//	if err != nil {
//		panic(err)
//	}
//	grpcAddress := fmt.Sprintf("%s:%d", registerIP, rpcPort)
//	err = em.AddEndpoint(ctx, "openmeeting/rtc-service/g1", endpoints.Endpoint{Addr: grpcAddress}, clientv3.WithLease(lease.ID))
//	if err != nil {
//		panic(err)
//	}
//
//	ch, err := cli.KeepAlive(ctx, lease.ID)
//
//	if err != nil {
//		return err
//	}
//
//	go func() {
//		for {
//			select {
//			case r := <-ch:
//				// avoid dead loop when channel was closed
//				if r == nil {
//					return
//				}
//			case <-ctx.Done():
//				return
//			}
//		}
//	}()
//
//	var wg errgroup.Group
//
//	wg.Go(func() error {
//		return errs.Wrap(srv.Serve(listener))
//	})
//
//	sigs := make(chan os.Signal, 1)
//	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
//	<-sigs
//
//	var (
//		done = make(chan struct{}, 1)
//		gerr error
//	)
//
//	go func() {
//		once.Do(srv.GracefulStop)
//		gerr = wg.Wait()
//		close(done)
//	}()
//
//	select {
//	case <-done:
//		return gerr
//
//	case <-time.After(15 * time.Second):
//		return errs.Wrap(errors.New("timeout exit"))
//	}
//
//}

func Start(ctx context.Context, prometheusConfig *config.Prometheus, listenIP,
	registerIP string, rpcPorts []int, index int, rpcRegisterName string, serverConfig subConfig.Config, share *config.Share,
	rpcFn func(ctx context.Context, serverConfig subConfig.Config, server *grpc.Server) error, options ...grpc.ServerOption) error {

	rpcPort, err := datautil.GetElemByIndex(rpcPorts, index)
	if err != nil {
		return err
	}
	prometheusPort, err := datautil.GetElemByIndex(prometheusConfig.Ports, index)
	if err != nil {
		return err
	}
	log.CInfo(ctx, "RPC server is initializing", "rpcRegisterName", rpcRegisterName,
		"rpcPort", rpcPort, "prometheusPort", prometheusPort)
	rpcTcpAddr := net.JoinHostPort(network.GetListenIP(listenIP), strconv.Itoa(rpcPort))
	listener, err := net.Listen(
		"tcp",
		rpcTcpAddr,
	)
	if err != nil {
		return errs.WrapMsg(err, "listen err", "rpcTcpAddr", rpcTcpAddr)
	}

	defer listener.Close()
	registerIP, err = network.GetRpcRegisterIP(registerIP)
	if err != nil {
		return err
	}

	var reg *prometheus.Registry
	var metric *grpcprometheus.ServerMetrics
	if prometheusConfig.Enable {
		cusMetrics := prommetrics.GetGrpcCusMetrics(rpcRegisterName, share)
		reg, metric, _ = prommetrics.NewGrpcPromObj(cusMetrics)
		options = append(options, mw.GrpcServer(), grpc.StreamInterceptor(metric.StreamServerInterceptor()),
			grpc.UnaryInterceptor(metric.UnaryServerInterceptor()))
	} else {
		options = append(options, mw.GrpcServer())
	}

	srv := grpc.NewServer(options...)
	once := sync.Once{}
	defer func() {
		once.Do(srv.GracefulStop)
	}()

	err = rpcFn(ctx, serverConfig, srv)
	if err != nil {
		return err
	}

	if err := etcd_discovery.InitEtcdService(ctx, &serverConfig, rpcTcpAddr, rpcPort); err != nil {
		return errs.WrapMsg(err, "init etcd failed")
	}

	var (
		netDone    = make(chan struct{}, 2)
		netErr     error
		httpServer *http.Server
	)

	go func() {
		if prometheusConfig.Enable && prometheusPort != 0 {
			metric.InitializeMetrics(srv)
			// Create a HTTP server for prometheus.
			httpServer = &http.Server{
				Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
				Addr:    fmt.Sprintf("0.0.0.0:%d", prometheusPort)}
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				netErr = errs.WrapMsg(err, "prometheus start err", httpServer.Addr)
				netDone <- struct{}{}
			}
		}
	}()

	go func() {
		err := srv.Serve(listener)
		if err != nil {
			netErr = errs.WrapMsg(err, "rpc start err: ", rpcTcpAddr)
			netDone <- struct{}{}
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	select {
	case <-sigs:
		program.SIGTERMExit()
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := gracefulStopWithCtx(ctx, srv.GracefulStop); err != nil {
			return err
		}
		ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err := httpServer.Shutdown(ctx)
		if err != nil {
			return errs.WrapMsg(err, "shutdown err")
		}
		return nil
	case <-netDone:
		close(netDone)
		return netErr
	}
}

func gracefulStopWithCtx(ctx context.Context, f func()) error {
	done := make(chan struct{}, 1)
	go func() {
		f()
		close(done)
	}()
	select {
	case <-ctx.Done():
		return errs.New("timeout, ctx graceful stop")
	case <-done:
		return nil
	}
}
