package pkg

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(logFilePath string) (*zap.Logger, error) {
	// Настройка ротации логов
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath, // Например, "logs/app.log"
		MaxSize:    10,          // Мб перед ротацией
		MaxBackups: 5,           // Максимум 5 резервных файлов
		MaxAge:     30,          // Дней хранения
		Compress:   true,        // Сжимать старые логи
	}

	// Настройка Zap для записи в файл и консоль
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(lumberjackLogger),
			zapcore.InfoLevel,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}
