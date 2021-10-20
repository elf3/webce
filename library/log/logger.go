package log

import (
	//rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

type LogConfig struct {
	Level   string `yaml:"level"`
	Path    string `yaml:"path"`
	MaxSize uint   `yaml:"save"`
}

var Log *zap.SugaredLogger

func InitLogger(logConfig LogConfig) {
	encoder := getEncoder()

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	infoWriter := getLogWriter(logConfig.Path, "Info", logConfig.MaxSize)

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	errorWriter := getLogWriter(logConfig.Path, "Error", logConfig.MaxSize)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)

	logger := zap.New(core, zap.AddCaller())
	Log = logger.Sugar()
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02 15:04:05]"))
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

//func getRoTateLogWriter(logPath, level string, maxSize uint) io.Writer {
//	logFullPath := path.Join(logPath, level)
//	hook, err := rotatelogs.New(
//		logFullPath+".%Y%m%d%H.log",                 // 没有使用go风格反人类的format格式
//		rotatelogs.WithLinkName(logFullPath+".log"), // 生成软链，指向最新日志文件
//		rotatelogs.WithRotationCount(maxSize),       // 文件最大保存份数
//		rotatelogs.WithRotationTime(24*time.Hour),   // 日志切割时间间隔
//	)
//	if err != nil {
//		panic(err)
//	}
//	return hook
//}

func getLogWriter(logPath, level string, maxSize uint) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path.Join(logPath, level) + ".log",
		MaxSize:    int(maxSize),
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
