package api_pb

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/go-kit/kit/log"
	pb "github.com/robertobadjio/tgtime-api/api/v1/pb/api"
	"github.com/robertobadjio/tgtime-notifier/internal/config"
)

// Client GRPC-клиент для получения пользователя из API-микросервиса
type Client struct {
	cfg    *config.Config
	logger log.Logger
	client pb.ApiClient
}

// NewClient Конструктор GRPC-клиента для получения пользователя из API-микросервиса
func NewClient(cfg config.Config, logger log.Logger) *Client {
	conn, _ := grpc.NewClient(
		buildAddress(cfg.TgTimeAPIHost, cfg.TgTimeAPIPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	return &Client{cfg: &cfg, logger: logger, client: pb.NewApiClient(conn)}
}

// GetUserByTelegramID Получение пользователя по telegram ID
func (tc Client) GetUserByTelegramID(
	ctx context.Context,
	telegramID int64,
) (*pb.GetUserByTelegramIdResponse, error) {
	user, err := tc.client.GetUserByTelegramId(
		ctx,
		&pb.GetUserByTelegramIdRequest{TelegramId: telegramID},
	)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return user, nil
}

// GetUserByMacAddress Получение пользователя по MAC-адресу
func (tc Client) GetUserByMacAddress(
	ctx context.Context,
	macAddress string,
) (*pb.GetUserByMacAddressResponse, error) {
	user, err := tc.client.GetUserByMacAddress(
		ctx,
		&pb.GetUserByMacAddressRequest{MacAddress: macAddress},
	)

	if err != nil {
		return nil, handleError(ctx, err)
	}

	return user, nil
}

func handleError(ctx context.Context, err error) error {
	if s, ok := status.FromError(err); ok {
		// Handle the error based on its status code
		if s.Code() == codes.NotFound {
			return fmt.Errorf("requested resource not found")
		}
		return fmt.Errorf("RPC error: %v, %v", s.Message(), ctx.Err())
	}

	return fmt.Errorf("Non-RPC error: %v", err)
}

func buildAddress(host, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}
