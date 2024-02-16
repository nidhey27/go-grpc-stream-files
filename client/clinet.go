package main

import (
	"context"
	"fmt"
	proto "go-stream-files/proto"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.UploadServiceClient

func main() {
	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client = proto.NewUploadServiceClient(conn)

	mb := 1024 * 1024 * 2
	uploadStreamFilePath("./1gb.bin", mb)
}

func uploadStreamFilePath(path string, batchSize int) {
	t := time.Now()
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	buff := make([]byte, batchSize)
	batchNumber := 1
	stream, err := client.Upload(context.TODO())
	if err != nil {
		panic(err)
	}

	for {
		num, err := file.Read(buff)
		if err == io.EOF {

			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		chunk := buff[:num]

		if err := stream.Send(&proto.UploadRequest{
			FilePath: path,
			Chunk:    chunk,
		}); err != nil {
			fmt.Println(err)
			return
		}

		log.Printf("Sent - batch #%v - size - %v\n", batchNumber, len(chunk))
		batchNumber += 1
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(time.Since(t))
	log.Printf("Sent - %v bytes - %s\n", res.GetSize(), res.GetMessage())
}
