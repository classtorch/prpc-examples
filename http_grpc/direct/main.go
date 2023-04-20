package main

import (
	"context"
	"github.com/classtorch/prpc"
	"github.com/classtorch/prpc-examples/http_grpc/api"
	"github.com/classtorch/prpc/grpc"
	"github.com/classtorch/prpc/logger"
	rawGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func main() {
	log := logger.NewDefaultLogger(log.New(os.Stderr, logger.DefaultLogPrefix, log.LstdFlags))
	conn := prpc.NewClientConn()

	// grpc direct
	err := conn.NewGrpcClientConn(context.Background(), "172.25.27.106:33000", grpc.WithOptions(rawGrpc.WithTransportCredentials(insecure.NewCredentials())))
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
	// http direct
	err = conn.NewHttpClientConn(context.Background(), "127.0.0.1:8000/account")
	if err != nil {
		log.Error(err)
		return
	}
	listReply, err := userClient.GetUserList(context.Background(), &api.GetUserListReq{Age: 18})
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(listReply)

}
