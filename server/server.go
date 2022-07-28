package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sort"
	"time"

	"calculator/proto"
	"google.golang.org/grpc"
)
var sieve [1000]bool

func getprimes (sieve *[1000]bool) {
	for i:=2; i<1000;i++{
		if sieve[i]==true{
			continue
		}else{

			for j:=i*i;j<1000;j+=i{ 
				sieve[j]=true
			}
		}
	}
}
type server struct {
	proto.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *proto.SumRequest) (resp *proto.SumResponse, err error) {
	fmt.Println("Starting Unary Streaming...")

	num1 := req.Num1
	num2 := req.Num2

	res := num1 + num2
	resp = &proto.SumResponse{
		Res: res,
	}
	return resp, nil
}

func (*server) PrimeNumbers(req *proto.PrimeRequest, resp proto.CalculatorService_PrimeNumbersServer) error {
	fmt.Println("Starting Server Side Streaming...")
	num := req.Num

	/*
	sieve := [1000]bool{}
	func(sieve *[1000]bool) {
		for i:=2; i<1000;i++{
			if sieve[i]==true{
				continue
			}else{

				for j:=i*i;j<1000;j+=i{ 
					sieve[j]=true
				}
			}
		}
	}(&sieve)
	*/
	for i := 2; i <= int(num); i++ {
		if !sieve[i] {
			res := proto.PrimeResponse{
				Res: int32(i),
			}
			time.Sleep(300 * time.Millisecond)
			resp.Send(&res)
		}
	}
	return nil

}

func (*server) ComputeAverage(stream proto.CalculatorService_ComputeAverageServer) error {
	fmt.Println("Client Side Streaming Detected in the server...")
	var sum int32
	var n int32
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.AverageResponse{
				Res: float32(sum) / float32(n),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream : %v", err)
		}
		sum += msg.Num
		n++

	}

}

func (*server) FindMaxNumber(stream proto.CalculatorService_FindMaxNumberServer) error {
	fmt.Println("Bidirectional Service Started..Computing...")
	sl := []int{}
	maxNow := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while receiving data from Bidirectional client : %v", err)
			return err
		}
		num := req.Num
		sl = append(sl, int(num))
		sort.Sort(sort.IntSlice(sl))
		if sl[len(sl)-1] > maxNow {
			maxNow = sl[len(sl)-1]
			sendErr := stream.Send(&proto.MaxnumberResponse{
				Res: int32(maxNow),
			})
			if sendErr != nil {
				log.Fatalf("Error while sending response to the Bidirectional Client : %v", err)
				return err
			}
		}

	}
}

func main() {
	for i:=0;i<100;i++{
		sieve[i]=false
	}
	getprimes(&sieve)
	fmt.Println("Initialising Server...")
	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to Listen: %v", err)
	}

	fmt.Println("Starting server...")

	s := grpc.NewServer()
	proto.RegisterCalculatorServiceServer(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatal("Server Failed !!! : %v", err)
	}
}
