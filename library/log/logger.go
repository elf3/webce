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

	encoderConfig.EncodeTime = customTimeEncoder            // 时间
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder  // 显示完整文件路径
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(logType string) zapcore.WriteSyncer {
	filename := path.Join(viper.GetString("zap.path"), logType)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s.log", filename), // 日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    viper.GetInt("zap.maxSize"),     // 文件大小限制,单位MB
		MaxBackups: viper.GetInt("zap.maxBackups"),  // 最大保留日志文件数量
		MaxAge:     viper.GetInt("zap.maxAge"),      // 日志文件保留天数
		Compress:   viper.GetBool("zap.compress"),   // 是否压缩处理
	}
	return zapcore.AddSync(lumberJackLogger)
}
