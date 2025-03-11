package grpc

import (
	"context"
	"errors"

	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/generator"
	pb "github.com/Alena-Kurushkina/shortener/internal/grpc/proto"
	"github.com/Alena-Kurushkina/shortener/internal/sherr"
)

type ShortenerGRPC struct {
	pb.UnimplementedUsersServer
	repo       Storager
	config     *config.Config
}

func (s *ShortenerGRPC) CreateShortening(ctx context.Context, in *pb.CreateShorteningRequest)(response *pb.ShorteningResponse, err error) {
	err=nil
	// generate shortening
	shortStr := generator.GenerateRandomString(15)

	insertErr := s.repo.Insert(ctx, in.UserId, shortStr, in.LongUrl)

	var existError *sherr.AlreadyExistError
	if errors.As(insertErr, &existError) {
		// make response
		response.Shortening=s.config.BaseURL + existError.ExistShortStr	
		return	
	} else if insertErr != nil {
		response.Error = insertErr.Error()
		return
	}

	response.Shortening=s.config.BaseURL + shortStr
	return
}