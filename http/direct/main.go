package main

import (
	"context"
	"github.com/classtorch/prpc"
	"github.com/classtorch/prpc-examples/http/api"
	"github.com/classtorch/prpc/http"
	"github.com/classtorch/prpc/logger"
	"log"
	rawHttp "net/http"
	"os"
)

func main() {
	Direct()
	DirectCustomHttpImpl()
}

// direct mode
func Direct() {
	log := logger.NewDefaultLogger(log.New(os.Stderr, logger.DefaultLogPrefix, log.LstdFlags))
	conn := prpc.NewClientConn()

	// direct
	err := conn.NewHttpClientConn(context.Background(), "127.0.0.1:8000/account")
	if err != nil {
		log.Error(err)
		return
	}
	userClient := api.NewUserClient(conn)
	listReply, err := userClient.GetUserList(context.Background(), &api.GetUserListReq{Name: "张三", Age: 18})
	if err != nil {
		log.Error(err)
	}
	log.Info(listReply)
	// api with variable param(rest api format)
	infoReply, err := userClient.GetUserInfo(context.Background(), &api.GetUserInfoReq{Name: "张三"}, http.WithUrlParams(map[string]string{"uid": "123"}))
	if err != nil {
		log.Error(err)
	}
	log.Info(infoReply)
}

// direct mode, custom http impl
func DirectCustomHttpImpl() {
	log := logger.NewDefaultLogger(log.New(os.Stderr, logger.DefaultLogPrefix, log.LstdFlags))
	conn := prpc.NewClientConn()

	// direct
	err := conn.NewHttpClientConn(context.Background(), "127.0.0.1:8000/account", http.WithCallClient(&myHttpImpl{}))
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

type myHttpImpl struct {
}

func (mockHttpImpl myHttpImpl) Get(ctx context.Context, addr string, api string, req interface{}, reply interface{}, opts ...http.CallOption) (*rawHttp.Request, *rawHttp.Response, error) {
	request, err := rawHttp.NewRequest(rawHttp.MethodGet, addr+api, nil)
	if err != nil {
		return nil, nil, err
	}
	return request, &rawHttp.Response{}, nil
}

func (mockHttpImpl myHttpImpl) Post(ctx context.Context, addr string, api string, req interface{}, reply interface{}, opts ...http.CallOption) (*rawHttp.Request, *rawHttp.Response, error) {
	return nil, nil, nil
}
func (mockHttpImpl myHttpImpl) Delete(ctx context.Context, addr string, api string, req interface{}, reply interface{}, opts ...http.CallOption) (*rawHttp.Request, *rawHttp.Response, error) {
	return nil, nil, nil
}
func (mockHttpImpl myHttpImpl) Put(ctx context.Context, addr string, api string, req interface{}, reply interface{}, opts ...http.CallOption) (*rawHttp.Request, *rawHttp.Response, error) {
	return nil, nil, nil
}
func (mockHttpImpl myHttpImpl) Default(ctx context.Context, addr string, api string, req interface{}, reply interface{}, opts ...http.CallOption) (*rawHttp.Request, *rawHttp.Response, error) {
	return nil, nil, nil
}
