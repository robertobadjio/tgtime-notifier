package api_pb

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pbapiv1 "github.com/robertobadjio/tgtime-api/api/v1/pb/api"
	"github.com/robertobadjio/tgtime-notifier/internal/config"
)

// Client ???
type Client interface {
	GetUserByTelegramID(ctx context.Context, telegramID int64) (
		*pbapiv1.GetUserByTelegramIdResponse,
		error,
	)
	GetUserByMacAddress(
		ctx context.Context,
		macAddress string,
	) (*pbapiv1.GetUserByMacAddressResponse, error)
}

// Client GRPC-клиент для получения пользователя из API-микросервиса
type client struct {
	cfg config.TgTimeAPIConfig
}

// NewClient Конструктор GRPC-клиента для получения пользователя из API-микросервиса
func NewClient(cfg config.TgTimeAPIConfig) Client {
	return &client{
		cfg: cfg,
	}
}

// GetUserByTelegramID Получение пользователя по telegram ID
func (tc *client) GetUserByTelegramID(
	ctx context.Context,
	telegramID int64,
) (*pbapiv1.GetUserByTelegramIdResponse, error) {
	conn, _ := grpc.NewClient(
		tc.cfg.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	defer func() {
		_ = conn.Close()
	}()

	client := pbapiv1.NewApiClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	user, err := client.GetUserByTelegramId(
		ctx,
		&pbapiv1.GetUserByTelegramIdRequest{TelegramId: telegramID},
	)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return user, nil
}

// GetUserByMacAddress Получение пользователя по MAC-адресу
func (tc *client) GetUserByMacAddress(
	ctx context.Context,
	macAddress string,
) (*pbapiv1.GetUserByMacAddressResponse, error) {
	conn, _ := grpc.NewClient(
		tc.cfg.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer func() {
		_ = conn.Close()
	}()

	client := pbapiv1.NewApiClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	user, err := client.GetUserByMacAddress(
		ctx,
		&pbapiv1.GetUserByMacAddressRequest{MacAddress: macAddress},
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
