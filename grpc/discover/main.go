package main

import (
	"context"
	"github.com/classtorch/prpc"
	"github.com/classtorch/prpc-examples/http_grpc/api"
	"github.com/classtorch/prpc/balancer/roundrobin"
	"github.com/classtorch/prpc/grpc"
	"github.com/classtorch/prpc/logger"
	"github.com/classtorch/prpc/resolver/consul"
	rawGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func main() {
	log := logger.NewDefaultLogger(log.New(os.Stderr, logger.DefaultLogPrefix, log.LstdFlags))
	conn := prpc.NewClientConn()
	consulResolver := consul.NewResolverBuilder()
	// consul resolver,roundrobin balancer,interceptor
	err := conn.NewGrpcClientConn(context.Background(), "consul://127.0.0.1:8000/account", grpc.WithResolver(consulResolver), grpc.WithBalancerName(roundrobin.Name), grpc.WithOptions(rawGrpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *rawGrpc.ClientConn, invoker rawGrpc.UnaryInvoker, opts ...rawGrpc.CallOption) error {
		log.Infof("grpc invoke start, method:%s req:%+v", method, req, reply)
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			log.Error("grpc invoke error,err:%+v method:%s req:%+v", err, method, req, reply)
		}
		log.Infof("grpc invoke end, method:%s req:%+v res:%+v", method, req, reply)
		return err
	}), rawGrpc.WithTransportCredentials(insecure.NewCredentials())))
	if err != nil {
		log.Error(err)
		return
	}
	userClient := api.NewUserClient(conn)
	_, err = userClient.AddUser(context.Background(), &api.AddUserReq{User: &api.UserInfo{Name: "张三", Age: 18}})
	if err != nil {
		log.Error(err)
		return
	}
}
