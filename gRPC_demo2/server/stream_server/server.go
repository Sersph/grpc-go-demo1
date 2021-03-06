package main

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	//pb "github.com/EDDYCJY/go-grpc-example/proto"

	pb "gRPC_demo2/po"
)

type StreamService struct{}

const (
	PORT    = ":9002"
	Network = "tcp"
)

func main() {
	//创建服务、注入服务、设置监听、启动监听
	server := grpc.NewServer()

	pb.RegisterStreamServiceServer(server, &StreamService{})

	lis, err := net.Listen(Network, PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}

//服务端 流模式
func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	for n := 0; n <= 8; n++ {
		//主要：stream.Send()
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		log.Println("服务端模式，接收客户端流------ Name: ", r.Pt.Name)
		log.Println("Value:", r.Pt.Value)

		if err != nil {
			return err
		}
	}
	return nil
}

//客户端 流模式
func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{Pt: &pb.StreamPoint{
				Name:  "服务端发送--> 这个请求 可以通过",
				Value: 1}})
		}
		if err != nil {
			return err
		}

		log.Println("客户端模式，")
		log.Printf("服务端发送： pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}
	return nil
}

//双向流模式
func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "双向模式: 服务端发送 --->  ",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}

		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		n++

		log.Printf("双向模式--服务端接收: pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}

	return nil
}
