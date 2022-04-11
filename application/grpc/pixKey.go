package grpc

import (
	"context"

	"github.com/guiPython/codepix/application/grpc/pb"
	"github.com/guiPython/codepix/domain/usecase"
)

type PixKeyGrpcService struct {
	pixKeyUseCase *usecase.PixKeyService
	pb.UnimplementedPixKeyServiceServer
}

func NewPixGrpcService(pixKeyUseCase *usecase.PixKeyService) *PixKeyGrpcService {
	grpcService := PixKeyGrpcService{pixKeyUseCase: pixKeyUseCase}
	return &grpcService
}

func (c *PixKeyGrpcService) RegisterPixKey(ctx context.Context, in *pb.PixKeyRegistration) (*pb.PixKeyCreatedResult, error) {
	kind := pb.Kind_value[in.Kind.Enum().String()]
	key, err := c.pixKeyUseCase.RegisterKey(kind, in.Key, in.AccountId)
	if err != nil {
		return &pb.PixKeyCreatedResult{
			Status: "not created",
			Error:  err.Error(),
		}, err
	}

	return &pb.PixKeyCreatedResult{
		Id:     key.ID,
		Status: "created",
	}, nil
}

func (c *PixKeyGrpcService) Find(ctx context.Context, in *pb.PixKey) (*pb.PixKeyInfo, error) {
	kind := pb.Kind_value[in.Kind.Enum().String()]
	pixKey, err := c.pixKeyUseCase.FindKey(kind, in.Key)
	if err != nil {
		return &pb.PixKeyInfo{}, err
	}
	return &pb.PixKeyInfo{
		Id:   pixKey.ID,
		Kind: in.Kind,
		Key:  pixKey.Key,
		Account: &pb.Account{
			AccountId:     pixKey.AccountID,
			AccountNumber: pixKey.Account.Number,
			BankId:        pixKey.Account.Bank.ID,
			BankName:      pixKey.Account.Bank.Name,
			OwnerName:     pixKey.Account.OwnerName,
			CreatedAt:     pixKey.Account.CreatedAt.String(),
		},
		CreatedAt: pixKey.CreatedAt.String(),
	}, nil
}
