package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"

	pb "github.com/Alena-Kurushkina/shortener/internal/grpc/proto"
)

type claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

const tokenExp = time.Hour * 3

const secretKey = "secretkey"

// buildJWTString makes token and returns it as a string.
func buildJWTString(id uuid.UUID) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		// собственное утверждение
		UserID: id,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient("localhost:3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.NewShortenerClient(conn)

	userID := uuid.NewV4()
	jwt, err := buildJWTString(userID)
	if err != nil {
		panic(err)
	}
	ctx := metadata.AppendToOutgoingContext(context.Background(), "token", jwt)

	// test client
	shortening, err := testCreateShortening(ctx, c, userID.String(), "http://long-url.ru/long/long")
	if err == nil {
		splits := strings.Split(shortening, "/")
		shrt := splits[len(splits)-1]
		testGetFull(ctx, c, shrt)
	}
	testCreateShorteningBatch(ctx, c, userID.String())
	testGetAll(ctx, c, userID.String())
	testDelete(ctx, c)
	testGetAll(ctx, c, userID.String())
}

func testCreateShortening(ctx context.Context, c pb.ShortenerClient, userID string, lURL string) (string, error) {
	longURL := pb.CreateShorteningRequest{
		LongUrl: lURL,
		UserId:  userID,
	}
	resp, err := c.CreateShortening(ctx, &longURL)
	if err != nil {
		fmt.Println("Ошибка ", err)
		return "", err
	}
	fmt.Println(resp.Shortening)
	return resp.Shortening, nil
}

func testGetFull(ctx context.Context, c pb.ShortenerClient, shortening string) {
	resp1, err := c.GetFullString(ctx, &pb.LongURLRequest{ShortUrl: shortening})
	if err != nil {
		fmt.Println("Ошибка ", err)
		return
	}
	fmt.Println(resp1.LongUrl)
}

func testGetAll(ctx context.Context, c pb.ShortenerClient, userID string) {
	resp1, err := c.GetUserAllShortenings(ctx, &pb.UserID{UserId: userID})
	if err != nil {
		fmt.Println("Ошибка ", err)
		return
	}
	fmt.Println(resp1.UrlBatch)
}

func testCreateShorteningBatch(ctx context.Context, c pb.ShortenerClient, userID string) {
	batch := []*pb.CreateShorteningBatchRequest_BatchRequest{}
	batch = append(batch, &pb.CreateShorteningBatchRequest_BatchRequest{
		CorrelationId: "1234",
		LongUrl:       "long-url1",
	})
	batch = append(batch, &pb.CreateShorteningBatchRequest_BatchRequest{
		CorrelationId: "5678",
		LongUrl:       "long-url2",
	})
	req := pb.CreateShorteningBatchRequest{
		UserId:   userID,
		UrlBatch: batch,
	}
	resp, err := c.CreateShorteningBatch(ctx, &req)
	if err != nil {
		fmt.Println("Ошибка ", err)
		return
	}
	fmt.Println(resp.UrlBatch)
}

func testDelete(ctx context.Context, c pb.ShortenerClient) {
	_, err := c.DeleteRecord(ctx, &pb.DeleteRecordRequest{CorrelationId: []string{"5678"}})
	if err != nil {
		fmt.Println("Ошибка ", err)
		return
	}
}
