package main

import (
	"context"
	"fmt"
	"grpc_go_course/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Unary RPC...")
	req := &calculatorpb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 40,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Calculator RPC: %v\n", err)
	}

	log.Printf("Response from Sum: %v", res.SumResult)

}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting PrimeNumberDecompositionRequest Server RPC...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 1200,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition RPC: %v\n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}

}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a ComputeAverage Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	numbers := []int32{3, 6, 7, 8, 10, 59}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while recieving response of ComputeAverage: %v", err)
	}
	fmt.Printf("The average is: %v", res.GetAverage())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a FindMaximun BiDi Streaming RPC...")

	stream, err := c.FindMaximun(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximun: %v", err)
	}

	waitc := make(chan struct{})

	//send go routine
	go func() {
		numbers := []int32{4, 7, 2, 1, 18, 32, 8, 3}
		for _, number := range numbers {
			fmt.Printf("Sending number: %v\n", number)
			stream.Send(&calculatorpb.FindMaximunRequest{
				Number: number,
			})
			time.Sleep(100 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// recieve go routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while recieving FindMaximun response stream: %v", err)
				break
			}

			maximun := res.GetMaximun()
			fmt.Printf("Recieved a new maximun of: %v\n", maximun)
		}
		close(waitc)
	}()

	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")

	doErrorCall(c, 10) // without error
	doErrorCall(c, -2) // with error

}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: n,
	})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Println("Error message from server:", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
				return
			}

		} else {
			log.Fatalf("Big error calling SquareRoot: %v\n", err)
			return
		}

	}

	fmt.Printf("Result of square root of %v is %v\n", n, res.GetNumberRoot())
}
