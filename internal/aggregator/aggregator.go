package aggregator

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	pb "github.com/robertobadjio/tgtime-aggregator/api/v1/pb/aggregator"
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

func (tc Client) GetTimeSummary(ctx context.Context, macAddress, date string) (*pb.GetTimeSummaryResponse, error) {
	client, err := grpc.NewClient(
		tc.buildAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}
	defer func() { _ = client.Close() }()

	timeAggregatorClient := pb.NewAggregatorClient(client)
	ctxTemp, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	filters := make([]*pb.Filter, 0, 2)
	filters = append(filters, &pb.Filter{Key: "mac_address", Value: macAddress})
	filters = append(filters, &pb.Filter{Key: "date", Value: date})
	timeSummary, err := timeAggregatorClient.GetTimeSummary(
		ctxTemp,
		&pb.GetTimeSummaryRequest{Filters: filters},
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

	return timeSummary, nil
}

func (tc Client) buildAddress() string {
	return fmt.Sprintf("%s:%s", tc.cfg.TgTimeAggregatorHost, tc.cfg.TgTimeAggregatorPort)
}
