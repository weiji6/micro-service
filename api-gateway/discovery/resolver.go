package discovery

import (
	"context"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"time"
)

const schema = "etcd"

type Resolver struct {
	schema      string
	EtcdAddress []string
	DialTimeout int

	closeCh        chan struct{}
	watchCh        clientv3.WatchChan
	cli            *clientv3.Client
	KeyPrifix      string
	srvAddressList []resolver.Address

	cc     resolver.ClientConn
	logger *logrus.Logger
}

func NewResolver(etcdAddress []string, logger *logrus.Logger) *Resolver {
	return &Resolver{
		schema:      schema,
		EtcdAddress: etcdAddress,
		DialTimeout: 3,
		logger:      logger,
	}
}

// Scheme returns the scheme supported by this resolver
func (r *Resolver) Scheme() string {
	return r.schema
}

// Build creates a new resolver.Resolver for the given target
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc

	r.KeyPrifix = BuildPrefix(Server{Name: target.Endpoint(), Version: target.URL.Host})
	if _, err := r.start(); err != nil {
		return nil, err
	}

	return r, nil
}

// ResolveNow resolver.Resolver interface
func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {}

// Close resolver.Resolver interface
func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}

// start
func (r *Resolver) start() (chan<- struct{}, error) {
	var err error
	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddress,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	resolver.Register(r)

	r.closeCh = make(chan struct{})

	if err = r.sync(); err != nil {
		return nil, err
	}

	go r.watch()

	return r.closeCh, nil
}

// watch update events
func (r *Resolver) watch() {
	ticker := time.NewTicker(time.Second)
	r.watchCh = r.cli.Watch(context.Background(), r.KeyPrifix, clientv3.WithPrefix())

	for {
		select {
		case <-r.closeCh:
			return
		case res, ok := <-r.watchCh:
			if !ok {
				r.update(res.Events)
			}
		case <-ticker.C:
			if err := r.sync(); err != nil {
				r.logger.Error("sync failed:", err)
			}
		}
	}
}

// update
func (r *Resolver) update(events []*clientv3.Event) {
	for _, event := range events {
		var info Server
		var err error

		switch event.Type {
		case clientv3.EventTypePut:
			info, err = ParseValue(event.Kv.Value)
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Address, Metadata: info.Weight}
			if !Exist(r.srvAddressList, addr) {
				r.srvAddressList = append(r.srvAddressList, addr)
				r.cc.UpdateState(resolver.State{Addresses: r.srvAddressList})
			}
		case clientv3.EventTypeDelete:
			info, err := SplitPath(string(event.Kv.Key))
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Address}
			if s, ok := Remove(r.srvAddressList, addr); ok {
				r.srvAddressList = s
				r.cc.UpdateState(resolver.State{Addresses: r.srvAddressList})
			}
		}
	}
}

// sync 同步获取所有地址信息
func (r *Resolver) sync() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := r.cli.Get(ctx, r.KeyPrifix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	r.srvAddressList = []resolver.Address{}

	for _, kv := range res.Kvs {
		info, err := ParseValue(kv.Value)
		if err != nil {
			continue
		}
		addr := resolver.Address{Addr: info.Address, Metadata: info.Weight}
		r.srvAddressList = append(r.srvAddressList, addr)
	}
	r.cc.UpdateState(resolver.State{Addresses: r.srvAddressList})
	return nil
}
