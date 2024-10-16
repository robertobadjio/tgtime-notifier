package api_pb

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/go-kit/kit/log"
	pbapiv1 "github.com/robertobadjio/tgtime-api/api/v1/pb/api"
	"github.com/robertobadjio/tgtime-notifier/internal/config"
)

// Client GRPC-клиент для получения пользователя из API-микросервиса
type Client struct {
	cfg    *config.Config
	logger log.Logger
	host   string
	port   string
}

// NewClient Конструктор GRPC-клиента для получения пользователя из API-микросервиса
func NewClient(cfg config.Config, logger log.Logger) *Client {
	return &Client{
		cfg:    &cfg,
		logger: logger,
		host:   cfg.TgTimeAPIHost,
		port:   cfg.TgTimeAPIPort,
	}
}

// GetUserByTelegramID Получение пользователя по telegram ID
func (tc Client) GetUserByTelegramID(
	ctx context.Context,
	telegramID int64,
) (*pbapiv1.GetUserByTelegramIdResponse, error) {
	conn, _ := grpc.NewClient(
		buildAddress(tc.port, tc.host),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	defer func() {
		_ = conn.Close()
	}()

	client := pbapiv1.NewApiClient(conn)

	//ctx, cancel := context.WithTimeout(ctx, time.Second)
	//defer cancel()

	fmt.Println(buildAddress(tc.port, tc.host))
	fmt.Println(telegramID)
	user, err := client.GetUserByTelegramId(
		ctx,
		&pbapiv1.GetUserByTelegramIdRequest{TelegramId: telegramID},
	)
	fmt.Println(user)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return user, nil
}

// GetUserByMacAddress Получение пользователя по MAC-адресу
func (tc Client) GetUserByMacAddress(
	ctx context.Context,
	macAddress string,
) (*pbapiv1.GetUserByMacAddressResponse, error) {
	conn, _ := grpc.NewClient(
		buildAddress(tc.port, tc.host),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer func() {
		_ = conn.Close()
	}()

	client := pbapiv1.NewApiClient(conn)

	//ctx, cancel := context.WithTimeout(ctx, time.Second)
	//defer cancel()

	fmt.Println(buildAddress(tc.port, tc.host))
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

func buildAddress(host, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}
