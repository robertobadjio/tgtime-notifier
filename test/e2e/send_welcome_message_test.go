package e2e

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	kafkaClient "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/wait"
)

const macAddress = "F0:98:9D:1C:93:F6"
const (
	inOfficeTopic = "in-office"
	partition     = 0
)
const (
	chatID int64 = -1002586394034
	// nolint : test token
	botTokenSender = "8132539422:AAHb_lH0dHEfmxcjS6RqfVOh3egX__t3lU4"
	// nolint : test token
	botTokenReader = "7597797393:AAFB-do3FwX5Rd_MPkmbJwfy7MsTLUTEpU8"
)
const (
	TGTimAPIGRPCServicePort = "4770"
	TGTimAPIHTTPServicePort = "4771"
)
const (
	pathToService = "../../tgtime-notifier"
	servicePort   = "8081"
)

func TestE2ESendWelcomeMessage(t *testing.T) {
	_, errExistsServiceBinFile := os.Stat(pathToService)
	require.NoError(t, errExistsServiceBinFile)

	ctx := context.Background()

	currPath, errPath := os.Getwd()
	require.NoError(t, errPath)

	t.Log("Starting API user container")
	APIContainer, errGenericContainer := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "tkpd/gripmock:v1.14",
			ExposedPorts: []string{TGTimAPIGRPCServicePort + "/tcp", TGTimAPIHTTPServicePort + "/tcp"},
			WaitingFor:   wait.ForExposedPort(),
			Cmd:          []string{"/proto/api.proto"},
			HostConfigModifier: func(config *container.HostConfig) {
				config.AutoRemove = true
				config.Binds = []string{currPath + "/api_proto:/proto"}
			},
		},
		Started: true,
	})
	require.NoError(t, errGenericContainer, "error starting API container")

	ip, errHost := APIContainer.Host(ctx)
	require.NoError(t, errHost, "error getting host IP")

	mappedPort, errMappedPort := APIContainer.MappedPort(ctx, "4771")
	require.NoError(t, errMappedPort, "error getting mapped port HTTP")

	mappedPortGRPC, errMappedPortGRPC := APIContainer.MappedPort(ctx, "4770")
	require.NoError(t, errMappedPortGRPC, "error getting mapped port gRPC")

	APIContainerURL := fmt.Sprintf("%s:%s", ip, mappedPort.Port())

	defer func() {
		errTerminate := testcontainers.TerminateContainer(APIContainer)
		require.NoError(t, errTerminate, "error terminating API container")
	}()

	stateAPIService, errGetStateAPIService := APIContainer.State(ctx)
	require.NoError(t, errGetStateAPIService)
	t.Log("API service state running:", stateAPIService.Running)

	t.Log("Added stub user")
	errAddStub := addStubRequestToMockedGrpcServer(APIContainerURL, currPath)
	require.NoError(t, errAddStub, "error adding stub user")

	kafkaContainer, err := kafka.Run(
		ctx,
		"confluentinc/confluent-local:7.5.0",
		kafka.WithClusterID("test-cluster"),
	)
	require.NoError(t, err)
	defer func() {
		errTerminateKafka := testcontainers.TerminateContainer(kafkaContainer)
		require.NoError(t, errTerminateKafka)
	}()

	state, errGetState := kafkaContainer.State(ctx)
	require.NoError(t, errGetState)
	t.Log("Kafka cluster name:", kafkaContainer.ClusterID)
	t.Log("Kafka cluster state running:", state.Running)

	kafkaBrokers, errGetBrokers := kafkaContainer.Brokers(ctx)
	require.NoError(t, errGetBrokers, "Error getting kafka brokers")
	require.NotEmpty(t, kafkaBrokers, "Empty kafka brokers")
	t.Log("Kafka brokers:", strings.Join(kafkaBrokers, ", "))

	time.Sleep(time.Second)

	t.Log("Start service")
	cmd := exec.Command(pathToService)
	cmd.Env = append(
		os.Environ(),
		"HTTP_PORT="+servicePort,
		"KAFKA_BROKER_1="+kafkaBrokers[0],
		"BOT_TOKEN="+botTokenSender,
		"WEBHOOK_PATH=telegram",
		"WEBHOOK_LINK=https://tgtime.ru/telegram",
		"TGTIME_AGGREGATOR_HOST=",
		"TGTIME_AGGREGATOR_PORT=1080",
		"TGTIME_API_HOST=",
		"TGTIME_API_PORT="+mappedPortGRPC.Port(),
	)
	require.NoError(t, cmd.Start())

	// output
	/*stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	require.NoError(t, cmd.Start())

	fmt.Println("produce in office message")
	err = produceInOffice(ctx, macAddress, kafkaBrokers[0])
	require.NoError(t, err)

	go func() {
		for {
			tmp := make([]byte, 1024)
			_, err = stdout.Read(tmp)
			fmt.Print(string(tmp))
			if err != nil {
				break
			}
		}
	}()*/

	timeAssert := time.Now().Unix() - 1

	time.Sleep(time.Second)

	t.Log("Produce in office message")
	errProduceInOffice := produceInOffice(ctx, macAddress, kafkaBrokers[0])
	require.NoError(t, errProduceInOffice)

	t.Log("Sleep 5 seconds")
	time.Sleep(5 * time.Second)

	t.Log("Get updates from telegram")
	bot, errNewBot := tgbotapi.NewBotAPI(botTokenReader)
	require.NoError(t, errNewBot)

	updates, errGetUpdates := bot.GetUpdates(tgbotapi.UpdateConfig{})
	require.NoError(t, errGetUpdates)

	t.Log("Time unix assert after message:", timeAssert)
	flag := false
	// TODO: get last messages
	for _, update := range updates {
		if update.ChannelPost == nil {
			continue
		}

		if update.ChannelPost.Date < int(timeAssert) {
			continue
		}

		if update.ChannelPost.Text == "Вы пришли в офис" {
			flag = true
			assert.True(t, update.ChannelPost.Chat.ID == chatID)
		}
	}
	assert.True(t, flag)

	testcontainers.CleanupContainer(t, APIContainer)
	testcontainers.CleanupContainer(t, kafkaContainer)

	t.Log("Gracefully terminate service")
	require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))

	time.Sleep(time.Second)
}

func produceInOffice(ctx context.Context, macAddress, address string) error {
	conn, errDialLeader := kafkaClient.DialLeader(
		ctx,
		"tcp",
		address,
		inOfficeTopic,
		partition,
	)
	if errDialLeader != nil {
		return fmt.Errorf("failed to dial leader: %w", errDialLeader)
	}

	errSetWriteDeadline := conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if errSetWriteDeadline != nil {
		return fmt.Errorf("failed to set deadline: %w", errSetWriteDeadline)
	}

	_, errWriteMessages := conn.WriteMessages(kafkaClient.Message{Value: []byte(macAddress)})
	if errWriteMessages != nil {
		return fmt.Errorf("failed to write messages: %w", errWriteMessages)
	}

	if errConnClose := conn.Close(); errConnClose != nil {
		return fmt.Errorf("failed to close writer: %w", errConnClose)
	}

	return nil
}

func addStubRequestToMockedGrpcServer(httpInputMockAddress, path string) error {
	httpClient := &http.Client{}

	stub, errReadFile := os.ReadFile(filepath.Clean(path + "/data/stub/add_user.json"))
	if errReadFile != nil {
		return fmt.Errorf("failed to read file: %w", errReadFile)
	}

	request, errNewRequest := http.NewRequest("POST", fmt.Sprintf(
		"http://%s/add",
		httpInputMockAddress,
	), bytes.NewReader(stub))
	if errNewRequest != nil {
		return fmt.Errorf("failed to create request: %w", errNewRequest)
	}
	resp, errDoRequest := httpClient.Do(request)
	if errDoRequest != nil {
		return fmt.Errorf("failed to add user: %w", errDoRequest)
	}

	_, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		return fmt.Errorf("failed to read body: %w", errReadAll)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return nil
}
