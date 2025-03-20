// Package api contains gRPC methods realizations for shortener
package api

import (
	"context"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	_ "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alena-Kurushkina/shortener/internal/core"
	pb "github.com/Alena-Kurushkina/shortener/internal/grpc/proto"
	"github.com/Alena-Kurushkina/shortener/internal/sherr"
)

// A UserID is user identificator.
type UserID string

var userIDKey UserID = "userUUID"

// A ShortenerGRPC realises grpc methods.
type ShortenerGRPC struct {
	pb.UnimplementedShortenerServer
	Core *core.ShortenerCore
}

// NewShortenerGRPC initializes structure ShortenerGRPC
func NewShortenerGRPC(core *core.ShortenerCore) *ShortenerGRPC {
	return &ShortenerGRPC{
		Core: core,
	}
}

func extractUserIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value("userUUID").(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("can't extract user id from context")
	}
	userID, err := uuid.FromString(id)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

// CreateShortening get long URL in body and retrieves base URL with shortening.
func (s *ShortenerGRPC) CreateShortening(ctx context.Context, in *pb.CreateShorteningRequest) (*pb.ShorteningResponse, error) {
	userID, err := extractUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	// generate shortening
	shortURL, err := s.Core.CreateShortening(ctx, userID, in.LongUrl)
	var existError *sherr.AlreadyExistError
	if err != nil {
		if errors.As(err, &existError) {
			shr := pb.ShorteningResponse{
				Shortening: shortURL,
			}
			return &shr, status.Errorf(codes.AlreadyExists, "Shortening record already exists")
		}
		return nil, err
	}

	return &pb.ShorteningResponse{
		Shortening: shortURL,
	}, nil
}

// CreateShorteningJSONBatch get long URLs retrieves set of shortenings.
func (s *ShortenerGRPC) CreateShorteningBatch(ctx context.Context, in *pb.CreateShorteningBatchRequest) (*pb.ShorteningBatchResponse, error) {
	var response pb.ShorteningBatchResponse

	userID, err := extractUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	batch := make([]core.BatchElement, 0, len(in.UrlBatch))
	for _, item := range in.UrlBatch {
		batch = append(batch, core.BatchElement{
			CorrelarionID: item.CorrelationId,
			OriginalURL:   item.LongUrl,
		})
	}

	batch, err = s.Core.CreateShorteningBatch(ctx, userID, batch)
	if err != nil {
		return nil, err
	}

	retBatch := make([]*pb.ShorteningBatchResponse_BatchResponse, 0, len(batch))
	for _, item := range batch {
		retBatch = append(retBatch, &pb.ShorteningBatchResponse_BatchResponse{
			CorrelationId: item.CorrelarionID,
			LongUrl:       item.OriginalURL,
			ShortUrl:      item.ShortURL,
		})
	}

	response.UrlBatch = retBatch

	return &response, nil
}

// GetFullString get shortening and returns long URL.
func (s *ShortenerGRPC) GetFullString(ctx context.Context, in *pb.LongURLRequest) (*pb.LongURLResponse, error) {
	var response pb.LongURLResponse

	longURL, err := s.Core.GetFullString(ctx, in.ShortUrl)
	if err != nil {
		if errors.Is(err, sherr.ErrDBRecordDeleted) {
			return nil, status.Errorf(codes.NotFound, "Record was deleted")
		}
		return nil, err
	}

	response.LongUrl = longURL

	return &response, nil
}

// DeleteRecord saves record's id for future deletion.
// Deletion itself is performed periodically.
func (s *ShortenerGRPC) DeleteRecord(ctx context.Context, in *pb.DeleteRecordRequest) (*pb.None, error) {
	userID, err := extractUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	s.Core.RegisterToDelete(ctx, in.CorrelationId, userID)

	return &pb.None{}, nil
}

// GetStats returns number of shorten URLs and users number.
func (s *ShortenerGRPC) GetStats(ctx context.Context, in *pb.None) (*pb.GetStatsResponse, error) {
	stats, err := s.Core.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetStatsResponse{
		UrlsCount:  int32(stats.URLs),
		UsersCount: int32(stats.Users),
	}, nil
}

// GetUserAllShortenings returns all user's shortenings.
func (s *ShortenerGRPC) GetUserAllShortenings(ctx context.Context, in *pb.UserID) (*pb.ShorteningBatchResponse, error) {
	userID, err := extractUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	batch, err := s.Core.GetAllUserShortenings(ctx, userID)
	if err != nil {
		if errors.Is(err, sherr.ErrNoShortenings) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, err
	}
	retBatch := make([]*pb.ShorteningBatchResponse_BatchResponse, 0, len(batch))
	for _, item := range batch {
		retBatch = append(retBatch, &pb.ShorteningBatchResponse_BatchResponse{
			CorrelationId: item.CorrelarionID,
			LongUrl:       item.OriginalURL,
			ShortUrl:      item.ShortURL,
		})
	}
	return &pb.ShorteningBatchResponse{
		UrlBatch: retBatch,
	}, err
}
