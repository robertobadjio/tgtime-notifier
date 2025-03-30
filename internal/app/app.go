package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/oklog/oklog/pkg/group"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/robertobadjio/tgtime-notifier/internal/logger"
	"github.com/robertobadjio/tgtime-notifier/internal/metric"
	"github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram"
)

// App ???
type App struct {
	serviceProvider *serviceProvider
	gGroup          group.Group
}

// NewApp ???
func NewApp(ctx context.Context) (*App, error) {
	a := &App{
		gGroup: group.Group{},
	}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initCancelInterrupt,
		a.initAPIGateway,
		a.initTGUpdateHandle,
		a.initCheckInOfficeConsumer,
		a.initCheckPreviousDayInfo,
		//a.initPrometheus,
		//a.initPyroscope,
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

// nolint
func (a *App) initPyroscope(_ context.Context) error {
	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "notify.app",
		ServerAddress:   "http://" + a.serviceProvider.PyroscopeConfig().Address(),
		Logger:          pyroscope.StandardLogger,
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initAPIGateway(ctx context.Context) error {
	httpListener, err := net.Listen("tcp", a.serviceProvider.HTTPConfig().Address())
	if err != nil {
		return err
	}

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      a.serviceProvider.HTTPServiceHandler(ctx),
	}

	logger.Log("transport", "HTTP", "addr", a.serviceProvider.HTTPConfig().Address())

	a.gGroup.Add(func() error {
		return srv.Serve(httpListener)
	}, func(err error) {
		logger.Info("component", "APIGateway", "err", err.Error())
		_ = srv.Shutdown(ctx)
	})

	return nil
}

func (a *App) initCheckInOfficeConsumer(ctx context.Context) error {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	a.gGroup.Add(
		func() error {
			err := a.serviceProvider.Kafka().ConsumeInOffice(ctxWithCancel)
			if err != nil {
				return fmt.Errorf("error on kafka consumer in office message %w", err)
			}
			return nil
		}, func(err error) {
			logger.Error("transport", "GRPC", "component", "API", "during", "Listen", "err", err.Error())
			cancel()
		},
	)
	return nil
}

func (a *App) initCheckPreviousDayInfo(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Minute)

	a.gGroup.Add(
		func() error {
			for {
				select {
				case <-ticker.C:
					h, m, _ := time.Now().Clock()
					if h == 12 && m == 0 {
						err := a.serviceProvider.PreviousDayInfo().Run(ctx)
						if err != nil {
							logger.Error("telegram", "updates", "type", "previous day info", "msg", err.Error())
						}
					}
				case <-ctx.Done():
					return fmt.Errorf("context canceled: %s", ctx.Err())
				}
			}
			/*t := time.Now()
			n := time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, t.Location())
			d := n.Sub(t)
			if d < 0 {
				n = n.Add(24 * time.Hour)
				d = n.Sub(t)
			}
			for {
				time.Sleep(d)
				d = 24 * time.Hour

				err := a.serviceProvider.PreviousDayInfo().Run(ctx)
				if err != nil {
					logger.Error("telegram", "updates", "type", "previous day info", "msg", err.Error())
				}
			}*/
		}, func(err error) {
			logger.Error("component", "previous day info", "err", err.Error())
			ticker.Stop()
		})

	return nil
}

func (a *App) initCancelInterrupt(_ context.Context) error {
	cancelInterrupt := make(chan struct{})
	a.gGroup.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		logger.Info("component", "cancel interrupt", "err", "context canceled")
		close(cancelInterrupt)
	})

	return nil
}

func (a *App) initTGUpdateHandle(ctx context.Context) error {
	/*srv := &http.Server{
		Addr:         a.serviceProvider.HTTPConfig().Address(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	a.gGroup.Add(
		func() error {
			logger.Log("transport", "HTTP", "addr", a.serviceProvider.HTTPConfig().Address())
			err := srv.ListenAndServe()
			if err != nil {
				return fmt.Errorf("error starting server handle telegram updates %w", err)
			}

			return nil
		}, func(err error) {
			logger.Error("component", "server handle telegram updates", "err", err.Error())
			_ = srv.Shutdown(ctx)
		},
	)*/

	logger.Log("notifier", "telegram", "name", a.serviceProvider.TgBot().Self.UserName, "msg", "authorized on account")
	logger.Log(
		"notifier", "telegram",
		"name", a.serviceProvider.TgBot().Self.UserName,
		"msg", "setting webhook",
		"url", a.serviceProvider.tgConfig.WebhookLink(),
	)

	updates := a.serviceProvider.TgBot().ListenForWebhook("/" + a.serviceProvider.TelegramConfig().WebhookPath())
	a.gGroup.Add(
		func() error {
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
			logger.Error("component", "tg bot updates", "err", err.Error())
			a.serviceProvider.TgBot().StopReceivingUpdates()
		},
	)
	return nil
}

// nolint
func (a *App) initPrometheus(ctx context.Context) error {
	err := metric.Init(ctx)
	if err != nil {
		return err
	}

	httpListener, err := net.Listen("tcp", a.serviceProvider.PromConfig().Address())
	if err != nil {
		return err
	}

	a.gGroup.Add(func() error {
		logger.Info(
			"transport", "HTTP",
			"component", "prometheus",
			"addr", a.serviceProvider.PromConfig().Address(),
		)

		sm := http.NewServeMux()
		sm.Handle(a.serviceProvider.PromConfig().Path(), promhttp.Handler())

		srv := &http.Server{
			Handler:      sm,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		return srv.Serve(httpListener)
	}, func(err error) {
		logger.Error("transport", "HTTP", "component", "prometheus", "during", "listen", "err", err.Error())
		_ = httpListener.Close()
	})

	return nil
}

// Run ???
func (a *App) Run() error {
	return a.gGroup.Run()
}
