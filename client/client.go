package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"
	"google.golang.org/grpc"
	"calculator/proto"
	
)

func Sum(c proto.CalculatorServiceClient) {
	fmt.Println("Starting Unary Sum CAll..\n")
	req := proto.SumRequest{
		Num1:  45,
		Num2: 56,
	}
	fmt.Println("Sent number 45 and 56 for addition\n")
	resp, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("error occured while calling Sum grpc unary call: %v\n", err)
	}

	log.Printf("Result of Summing two numbers : %v\n", resp.Res)

}

func PrimeNumbers(c proto.CalculatorServiceClient) {
	fmt.Println("Starting Server Side Prime number Streaming..\n")
	req := proto.PrimeRequest{
		Num: 105,
	}
	fmt.Println("Sent number 105 for getting prime numbers\n")
	resStream, err := c.PrimeNumbers(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling Prime server-side streaming grpc : %v", err)
	}
	fmt.Println("Response Achieved From Server ")
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//we have reached to the end of the file
			break
		}

		if err != nil {
			log.Fatalf("error while receving server side stream : %v", err)
		}
		
		fmt.Println("The next Prime Number is:", msg.Res)
	}

}

func ComputeAverage(c proto.CalculatorServiceClient) {
	fmt.Println("Starting Client Side Streaming service for computing Average...\n")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client-side streaming : %v", err)
	}
	req := []*proto.AverageRequest{
		&proto.AverageRequest{Num: 2},
		&proto.AverageRequest{Num: 13},
		&proto.AverageRequest{Num: 45},
		&proto.AverageRequest{Num: 34},
		&proto.AverageRequest{Num: 78},
		&proto.AverageRequest{Num: 44},
		&proto.AverageRequest{Num: 67},
		&proto.AverageRequest{Num: 56},
		&proto.AverageRequest{Num: 34},
		&proto.AverageRequest{Num: 67},
		&proto.AverageRequest{Num: 87},
		&proto.AverageRequest{Num: 16},
	}
	for _, val := range req {
		fmt.Println("Sending number to the server for computation.... : ", val)
		stream.Send(val)
		time.Sleep(500 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from server : %v", err)
	}
	fmt.Println("\n Average of the Numbers is : ", resp.Res)

}

func FindMaxNumber(c proto.CalculatorServiceClient) {
	fmt.Println("Starting Bidirectional Streaming service")
	req := []*proto.MaxnumberRequest{
		&proto.MaxnumberRequest{Num: 1},
		&proto.MaxnumberRequest{Num: 3},
		&proto.MaxnumberRequest{Num: 7},
		&proto.MaxnumberRequest{Num: 9},
		&proto.MaxnumberRequest{Num: 2},
		&proto.MaxnumberRequest{Num: 5},
		&proto.MaxnumberRequest{Num: 22},
		&proto.MaxnumberRequest{Num: 15},
		&proto.MaxnumberRequest{Num: 21},
		&proto.MaxnumberRequest{Num: 19},
	}
	stream, err := c.FindMaxNumber(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client side streaming : %v", err)
	}

	waitchan := make(chan int32)

	go func(req []*proto.MaxnumberRequest) {
		for _, val := range req {
			fmt.Println("\nSending number... : ", val.Num)
			err := stream.Send(val)
			if err != nil {
				log.Fatalf("error while sending request to Maxnumber service : %v", err)
			}
			time.Sleep(1000 * time.Millisecond)

		}
		stream.CloseSend()
	}(req)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitchan)
				return
			}

			if err != nil {
				log.Fatalf("Error receiving response from server : %v", err)
			}

			fmt.Printf("\nResponse From Server : %v", resp.Res)
		}

	}()
	<-waitchan

}

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error occured while connecting: %v", err)
	}
	defer cc.Close()

	c := proto.NewCalculatorServiceClient(cc)
	//Sum(c)
	PrimeNumbers(c)
	//ComputeAverage(c)
	//FindMaxNumber(c)
}
