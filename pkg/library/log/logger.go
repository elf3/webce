package log

import (
	"fmt"
	"path"

	//rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type LogConfig struct {
	Level   string `yaml:"level"`
	Path    string `yaml:"path"`
	MaxSize uint   `yaml:"save"`
}

var Log *zap.SugaredLogger

func InitLogger() {
	encoder := GetEncoder()

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	infoWriter := getLogWriter("info")

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	errorWriter := getLogWriter("error")

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)

	logger := zap.New(core, zap.AddCaller())
	Log = logger.Sugar()
}

func InitGormLogger() *zap.Logger {
	encoder := GetEncoder()

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	infoWriter := getLogWriter("mysql")

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), errorLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), zap.InfoLevel),
	)

	return zap.New(core, zap.AddCaller())
}
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02 15:04:05]"))
}

func GetEncoder() zapcore.Encoder {
	zapType := viper.GetString("runmode")
	var encoderConfig zapcore.EncoderConfig
	switch zapType {
	case "prod":
		encoderConfig = zap.NewProductionEncoderConfig()
	case "debug":
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	encoderConfig.EncodeTime = customTimeEncoder            // ??????
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // ???????????????????????????????????????????????????zapcore.CapitalLevelEncoder????????????
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder  // ????????????????????????
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(logType string) zapcore.WriteSyncer {
	filename := path.Join(viper.GetString("zap.path"), logType)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s.log", filename), // ??????????????????????????????????????????????????????????????????
		MaxSize:    viper.GetInt("zap.maxSize"),     // ??????????????????,??????MB
		MaxBackups: viper.GetInt("zap.maxBackups"),  // ??????????????????????????????
		MaxAge:     viper.GetInt("zap.maxAge"),      // ????????????????????????
		Compress:   viper.GetBool("zap.compress"),   // ??????????????????
	}
	return zapcore.AddSync(lumberJackLogger)
}
