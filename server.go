package main

import (
	proto "go-stream-files/proto"
	"io"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	proto.UnimplementedUploadServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	svr := grpc.NewServer()
	proto.RegisterUploadServiceServer(svr, &Server{})
	reflection.Register(svr)

	if err := svr.Serve(listener); err != nil {
		panic(err)
	}
}

func (s *Server) Upload(stream proto.UploadService_UploadServer) error {
	var fileBytes []byte
	var fileSize int64 = 0

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// fileName = req.GetFilePath()
				break
			}
		}
		chunks := req.GetChunk()
		fileBytes = append(fileBytes, chunks...)
		fileSize += int64(len(chunks))
	}
	f, err := os.Create("./abc.bin")
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(fileBytes)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&proto.UploadResponse{
		Size:    fileSize,
		Message: "File written successfully!",
	})
}
