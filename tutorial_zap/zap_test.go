package tutorial_zap

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestZap_Sugar(t *testing.T) {
	/*
		In contexts where performance is nice, but not critical, use the SugaredLogger. It's 4-10x faster than other structured logging packages and includes both structured and printf-style APIs.
		在性能不错但不关键的上下文中，请使用 SugaredLogger .它比其他结构化日志记录包快 4-10 倍，并且包括结构化 API 和 printf 样式 API。
	*/
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	sugar := logger.Sugar()

	sugar.Infof("role %d enter map %s|%d", 10001, "主城", 666)
}

func TestZap_Logger(t *testing.T) {
	/*
		When performance and type safety are critical, use the Logger. It's even faster than the SugaredLogger and allocates far less, but it only supports structured logging.
		当性能和类型安全至关重要时，请使用 Logger .它甚至比 更快 SugaredLogger ，分配的也少得多，但它只支持结构化日志记录。
	*/
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	url := "https://baidu.com"
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	name := "张三"
	logger.Info("create role fail", zap.Error(fmt.Errorf("duplicate name %s", name)))
}

func TestZap_LoggerFile(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	// 输出日志数据到文件中
	config.OutputPaths = append(config.OutputPaths, "./zap.log")
	logger, _ := config.Build()

	for i := 0; i < 2; i++ {
		logger.Info("server start")
		time.Sleep(time.Second)
		logger.Info("server stop")
	}
}
