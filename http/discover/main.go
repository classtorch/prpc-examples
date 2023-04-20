package main

import (
	"context"
	"github.com/classtorch/prpc"
	"github.com/classtorch/prpc-examples/http/api"
	"github.com/classtorch/prpc/balancer/roundrobin"
	"github.com/classtorch/prpc/http"
	"github.com/classtorch/prpc/logger"
	"github.com/classtorch/prpc/resolver/consul"
	"log"
	rawHttp "net/http"
	"os"
	"strings"
)

func main() {
	log := logger.NewDefaultLogger(log.New(os.Stderr, logger.DefaultLogPrefix, log.LstdFlags))
	conn := prpc.NewClientConn()
	// consul resolver,roundrobin balancer,interceptor
	consulResolver := consul.NewResolverBuilder()
	err := conn.NewHttpClientConn(context.Background(), "consul://127.0.0.1:8000/account", http.WithResolver(consulResolver), http.WithBalancerName(roundrobin.Name), http.WithInterceptor(
		func(ctx context.Context, req interface{}, reply interface{}, httpRequest *rawHttp.Request, httpResponse *rawHttp.Response, cc *http.ClientConn, invoker http.Invoker, option ...http.CallOption) error {
			addr := strings.Split(httpRequest.Host, "//")[1]
			log.Infof("http invoke start, host:%s addr:%s req:%+v", httpRequest.Host, addr, req)
			err := invoker(ctx, req, reply, httpRequest, httpResponse, cc, option...)
			if err != nil {
				log.Error("http invoke error, err:%+v host:%s addr:%s req:%+v", err, httpRequest.Host, addr, req)
			}
			log.Infof("http invoke end, host:%s addr:%s req:%+v res:%+v", httpRequest.Host, addr, req, reply)
			return err
		}))
	if err != nil {
		log.Error(err)
		return
	}
	userClient := api.NewUserClient(conn)
	listReply, err := userClient.GetUserList(context.Background(), &api.GetUserListReq{Age: 18})
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(listReply)
}
