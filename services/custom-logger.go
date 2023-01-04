package services

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CustomLogger struct {
	log *zap.Logger
}

/*
The custom logger save logs to the fileDir.
The prefix using for the logs name.
The lifetime using to set the frequency of creating a new log file.
*/
func NewCustomLogger(prefix string, fileDir string,
	lifetime time.Duration) (*CustomLogger, error) {

	// set by default, if needed
	if prefix == "" {
		prefix = "logs"
	}

	logPath := path.Join(fileDir, prefix)

	rotator, err := rotatelogs.New(
		fmt.Sprintf("%s-%s-%s", logPath, "%Y-%m-%d-%H-%M-%S", ".log"),
		rotatelogs.WithRotationTime(lifetime),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create rotation: %s", err)
	}

	// initialize the JSON encoding config
	encoderConfig := map[string]string{
		"levelEncoder": "capital",
		"timeKey":      "date",
		"timeEncoder":  "iso8601",
	}
	data, _ := json.Marshal(encoderConfig)
	var encCfg zapcore.EncoderConfig
	if err := json.Unmarshal(data, &encCfg); err != nil {
		return nil, fmt.Errorf("failed to convert log file %s", err)
	}

	// add the encoder config and rotator to create a new zap logger
	w := zapcore.AddSync(rotator)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		w,
		zap.InfoLevel)
	logger := zap.New(core)

	return &CustomLogger{
		log: logger,
	}, nil
}

func (l CustomLogger) Info(message string, v ...interface{}) {
	for _, val := range v {
		b, _ := json.Marshal(val)
		l.log.Info("", zap.Any(message, json.RawMessage(b)))
	}
}

func (l CustomLogger) Error(message string, v ...interface{}) {
	for _, val := range v {
		b, _ := json.Marshal(val)
		l.log.Error("", zap.Any(message, json.RawMessage(b)))
	}
}

func (l CustomLogger) Warn(message string, v ...interface{}) {
	l.log.Sugar().Warnf(message, v...)
}

func (l CustomLogger) Debug(message string, v ...interface{}) {
	l.log.Sugar().Debugf(message, v...)
}
