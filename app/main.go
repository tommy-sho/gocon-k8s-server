package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Home         string   `env:"HOME"`
	Port         int      `env:"PORT" envDefault:"3000"`
	IsProduction bool     `env:"PRODUCTION"`
	Hosts        []string `env:"HOSTS" envSeparator:":"`
	SecretKey    string   `env:"SECRET_KEY,required"`
}

// stackdriverなどのエラーレベルに合わせたLogLevelを用意する
var zapCoreLevel = map[zapcore.Level]string{
	zapcore.DebugLevel: "DEBUG",
	zapcore.WarnLevel:  "WARN",
	zapcore.InfoLevel:  "INFO",
	zapcore.ErrorLevel: "ERROR",
}

func main() {
	//  タイムゾーンの名前とUTCとの差分となる秒数を引数で渡す
	fmt.Println(time.Now().Format("2006-01-02 MST"))

	cfg := zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "severity",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    encodeLevel,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.NanosDurationEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	server := echo.New()
	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	server.GET("/healthz", HealthCheck)

	go func() {
		logger.Info("server start serving")
		if err := server.Start(":8080"); err != http.ErrServerClosed {
			logger.Fatal("Server closed with error:", zap.Error(err))

		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	<-stopChan
	// shutdown処理

	logger.Info("start server shutdown")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Shutdown(ctx)
	}()

	select {
	case <-errChan:
		logger.Info("server shutdown")
	case <-ctx.Done():
		logger.Info("shutdown time exceed, force killed")
	}

	conf := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", conf)

}

func HealthCheck(c echo.Context) error {

	return c.String(http.StatusOK, "OK!")
}

func encodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(zapCoreLevel[l])
}
