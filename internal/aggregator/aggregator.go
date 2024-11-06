package aggregator

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
	"github.com/robertobadjio/tgtime-notifier/internal/config"
)

// Client ???
type Client interface {
	GetTimeSummary(
		ctx context.Context,
		macAddress, date string,
	) (*pb.GetSummaryResponse, error)
}

// Client GRPC-клиент для подключения к сервису Агрегатор
type client struct {
	client pb.TimeV1Client
}

// NewClient Конструктор gRPC-клиента для подключения к сервису Агрегатор
func NewClient(cfg config.TgTimeAPIConfig) Client {
	conn, _ := grpc.NewClient(
		cfg.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	return &client{client: pb.NewTimeV1Client(conn)}
}

// GetTimeSummary Получение времени сотрудника
func (tc *client) GetTimeSummary(
	ctx context.Context,
	macAddress, date string,
) (*pb.GetSummaryResponse, error) {
	filters := make([]*pb.Filter, 0, 2)
	if macAddress != "" {
		filters = append(filters, &pb.Filter{Key: "mac_address", Value: macAddress})
	}

	if date != "" {
		filters = append(filters, &pb.Filter{Key: "date", Value: date})
	}

	timeSummary, err := tc.client.GetSummary(
		ctx,
		&pb.GetSummaryRequest{Filters: filters},
	)

	if err != nil {
		return nil, handleError(ctx, err)
	}

	return timeSummary, nil
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
