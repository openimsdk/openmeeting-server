package etcd_discovery

import (
	"log"
	"sync"
	"time"
)

// Registrar registers service instance liveness information to etcd.
type Registrar struct {
	client      EtcdInterface
	serviceName string // e.g. "/service/foobar/1.2.3.4:8080"
	serviceAddr string // eg."http://1.2.3.4:8080"
	lease       int64
	logger      log.Logger

	quitmtx sync.Mutex
	quit    chan struct{}
}

// TTLOption allow setting a key with a TTL. This option will be used by a loop
// goroutine which regularly refreshes the lease of the key.
type TTLOption struct {
	heartbeat time.Duration // e.g. time.Second * 3
	ttl       time.Duration // e.g. time.Second * 10
}

// NewRegistrar returns a etcd Registrar acting on the provided catalog
// registration (service).
func NewRegistrar(client EtcdInterface, serviceName, serviceAddr string, lease int64) *Registrar {
	return &Registrar{
		client:      client,
		serviceName: serviceName,
		serviceAddr: serviceAddr,
		lease:       lease,
	}
}

// Register implements the sd.Registrar interface. Call it when you want your
// service to be registered in etcd, typically at startup.
func (r *Registrar) Register() {
	if err := r.client.Register(r.serviceName, r.serviceAddr, r.lease); err != nil {
		r.logger.Printf("err:%v", err)
		return
	}
}

// Deregister implements the sd.Registrar interface. Call it when you want your
// service to be deregistered from etcd, typically just prior to shutdown.
func (r *Registrar) Deregister() {
	if err := r.client.Deregister(r.serviceName); err != nil {
		r.logger.Printf("Deregister error:%v", err)
	} else {
		r.logger.Printf("service %s deregister success", r.serviceName)
	}

	r.quitmtx.Lock()
	defer r.quitmtx.Unlock()
	if r.quit != nil {
		close(r.quit)
		r.quit = nil
	}
}
