package shortener

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/Alena-Kurushkina/shortener/internal/authenticator"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/core"
	gr "github.com/Alena-Kurushkina/shortener/internal/grpc/api"
	pb "github.com/Alena-Kurushkina/shortener/internal/grpc/proto"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
)

// A ServerGRPC aggregates handler and config.
type ServerGRPC struct {
	Config          *config.Config
	IdleConnsClosed chan struct{}
	Server          *grpc.Server
}

// ChainUnaryInterceptors объединяет несколько унарных перехватчиков в один
func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Создаем цепочку вызовов
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = func(next grpc.UnaryHandler, interceptor grpc.UnaryServerInterceptor) grpc.UnaryHandler {
				return func(ctx context.Context, req interface{}) (interface{}, error) {
					return interceptor(ctx, req, info, next)
				}
			}(chain, interceptors[i])
		}
		return chain(ctx, req)
	}
}

func NewServerGRPC(core *core.ShortenerCore, cfg *config.Config, idleConnsChan chan struct{}) *ServerGRPC {
	chain := ChainUnaryInterceptors(
		authenticator.GRPCAuthInterceptor,
		logger.GRPCLogInterceptor,
	)
	// создаём gRPC-сервер без зарегистрированной службы
	server := grpc.NewServer(
		grpc.UnaryInterceptor(chain),
	)
	// регистрируем сервис
	pb.RegisterShortenerServer(server, &gr.ShortenerGRPC{
		Core: core,
	})

	return &ServerGRPC{
		Config:          cfg,
		IdleConnsClosed: idleConnsChan,
		Server:          server,
	}
}

func (s *ServerGRPC) Run() {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", s.Config.GRPCServerAddress)
	if err != nil {
		logger.Log.Fatalf("gRPC server Listen: %v", err)
	}

	logger.Log.Infof("Server gRPC is listening on %s", s.Config.GRPCServerAddress)
	// получаем запрос gRPC
	if err := s.Server.Serve(listen); err != nil {
		logger.Log.Fatal(err)
	}
}
