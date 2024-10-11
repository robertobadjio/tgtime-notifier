package api_pb

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	pb "github.com/robertobadjio/tgtime-api/api/v1/pb/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"tgtime-notifier/internal/config"
	"time"
)

type Client struct {
	cfg    *config.Config
	logger log.Logger
}

func NewClient(cfg config.Config, logger log.Logger) *Client {
	return &Client{cfg: &cfg, logger: logger}
}

func (tc Client) GetUserByTelegramId(ctx context.Context, telegramId int64) (*pb.GetUserByTelegramIdResponse, error) {
	client, err := grpc.NewClient(
		tc.buildAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}
	defer func() { _ = client.Close() }()

	apiClient := pb.NewApiClient(client)
	ctxTemp, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	user, err := apiClient.GetUserByTelegramId(
		ctxTemp,
		&pb.GetUserByTelegramIdRequest{TelegramId: telegramId},
	)

	if err != nil {
		if s, ok := status.FromError(err); ok {
			// Handle the error based on its status code
			if s.Code() == codes.NotFound {
				return nil, fmt.Errorf("requested resource not found")
			} else {
				return nil, fmt.Errorf("RPC error: %v, %v", s.Message(), ctxTemp.Err())
			}
		} else {
			// Handle non-RPC errors
			return nil, fmt.Errorf("Non-RPC error: %v", err)
		}
	}

	return user, nil
}

func (tc Client) buildAddress() string {
	return fmt.Sprintf("%s:%s", tc.cfg.TgTimeApiHost, tc.cfg.TgTimeApiPort)
}
