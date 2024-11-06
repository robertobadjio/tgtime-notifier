package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oklog/oklog/pkg/group"

	"github.com/robertobadjio/tgtime-notifier/internal/logger"
	"github.com/robertobadjio/tgtime-notifier/internal/notifier/telegram"
)

// App ???
type App struct {
	serviceProvider *serviceProvider
	apiGateway      group.Group
}

// NewApp ???
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initAPIGateway,
		a.initTGBot,
		a.initTGUpdateHandle,
		a.initCheckInOfficeConsumer,
		a.initCheckPreviousDayInfoTask,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			logger.Fatal(
				"init", "deps",
				"error", err.Error(),
			)
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initAPIGateway(_ context.Context) error {
	var g group.Group
	{
		srv := &http.Server{
			Addr:         a.serviceProvider.HTTPConfig().Address(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		logger.Log("transport", "HTTP", "addr", a.serviceProvider.HTTPConfig().Address())
		err := srv.ListenAndServe()
		if err != nil {
			logger.Error("telegram", "updates", "type", "serve", "msg", err)
		}
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	a.apiGateway = g

	return nil
}

func (a *App) initCheckInOfficeConsumer(ctx context.Context) error {
	a.apiGateway.Add(
		func() error {
			err := a.serviceProvider.Kafka().ConsumeInOffice(ctx)
			if err != nil {
				return fmt.Errorf("error on kafka consumer in office message %w", err)
			}
			return nil
		}, func(err error) {
			logger.Error("transport", "GRPC", "component", "API", "during", "Listen", "err", err.Error())
			// TODO: !?
		},
	)
	return nil
}

func (a *App) initCheckPreviousDayInfoTask(ctx context.Context) error {
	a.apiGateway.Add(
		func() error {
			t := time.Now()
			n := time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, t.Location())
			d := n.Sub(t)
			if d < 0 {
				n = n.Add(24 * time.Hour)
				d = n.Sub(t)
			}
			for {
				time.Sleep(d)
				d = 24 * time.Hour

				err := a.serviceProvider.PreviousDayInfoTask().Run(ctx)
				if err != nil {
					logger.Error("telegram", "updates", "type", "previous day info", "msg", err.Error())
				}
			}
		}, func(err error) {
			logger.Error("transport", "GRPC", "component", "API", "during", "Listen", "err", err.Error())
			// TODO: ?!
		})

	return nil
}

func (a *App) initTGBot(_ context.Context) error {
	a.apiGateway.Add(
		func() error {
			TGBot := a.serviceProvider.TgBot()
			logger.Log("notifier", "telegram", "name", TGBot.Self.UserName, "msg", "authorized on account")

			err := a.TGBotSetWebhook(TGBot)
			if err != nil {
				return fmt.Errorf("error setting webhook: %w", err)
			}

			logger.Log(
				"notifier", "telegram",
				"name", TGBot.Self.UserName,
				"msg", "setting webhook",
				"url", a.serviceProvider.tgConfig.GetWebhookLink(),
			)

			return nil
		}, func(err error) {
			logger.Error("transport", "GRPC", "component", "API", "during", "Listen", "err", err.Error())
			// TODO: !?
		},
	)
	return nil
}

func (a *App) initTGUpdateHandle(ctx context.Context) error {
	a.apiGateway.Add(
		func() error {
			srv := &http.Server{
				Addr:         a.serviceProvider.HTTPConfig().Address(),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}
			logger.Log("transport", "HTTP", "addr", a.serviceProvider.HTTPConfig().Address())
			err := srv.ListenAndServe()
			if err != nil {
				return fmt.Errorf("error starting server hadle telegram updates %w", err)
			}

			return nil
		}, func(err error) {
			logger.Error("transport", "GRPC", "component", "API", "during", "Listen", "err", err.Error())
			// TODO: !?
		},
	)
	a.apiGateway.Add(
		func() error {
			updates := a.serviceProvider.TgBot().ListenForWebhook("/" + a.serviceProvider.TelegramConfig().GetWebhookPath())
			for update := range updates {
				if update.Message == nil {
					continue
				}

				err := a.serviceProvider.TGNotifier().SendCommandMessage(ctx, telegram.ParamsUpdate{Update: update})
				if err != nil {
					logger.Log("telegram", "updates", "type", "info", "err", err)
				}
			}
			return nil
		}, func(err error) {
			logger.Error("transport", "GRPC", "component", "API", "during", "Listen", "err", err.Error())
			// TODO: !?
		},
	)
	return nil
}

// TGBotSetWebhook ???
func (a *App) TGBotSetWebhook(bot *TGBotAPI.BotAPI) error {
	_, err := bot.SetWebhook(TGBotAPI.NewWebhook(a.serviceProvider.tgConfig.GetWebhookLink()))
	if err != nil {
		return fmt.Errorf("set telegram webhook: %w", err)
	}

	//info, err := bot.GetWebhookInfo()
	_, err = bot.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("get telegram webhook info: %w", err)
	}
	/*if info.LastErrorDate != 0 {
		return fmt.Errorf("telegram callback failed: %s", info.LastErrorMessage)
	}*/

	return nil
}

// Run ???
func (a *App) Run() error {
	return a.apiGateway.Run()
}
