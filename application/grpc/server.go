package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/guiPython/codepix/application/grpc/pb"
	"github.com/guiPython/codepix/domain/usecase"
	repository "github.com/guiPython/codepix/infrastructure/repository/pixKey"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(database *gorm.DB, port int) {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	//pixKeyRepository := repository.NewPixKeyRepositoryInMemory()
	pixKeyRepository := repository.NewPixKeyRepository(database)
	pixKeyUseCase := usecase.NewPixKeyService(pixKeyRepository)

	pixKeyGrpcService := NewPixGrpcService(pixKeyUseCase)
	pb.RegisterPixKeyServiceServer(grpcServer, pixKeyGrpcService)

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("cannot create listener for grpc server ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server", err)
	}

	log.Printf("gRPC server has been started on port %d", port)
}
