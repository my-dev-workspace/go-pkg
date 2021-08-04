package logger

// Formatter:
// https://github.com/TV4/logrus-stackdriver-formatter
//
// Можно переделать
// https://github.com/mumoshu/logrus-bunyan-formatter/blob/master/bunyan_formatter.go
// https://github.com/JustDaile/jd_logrus_formatter/blob/master/formatter.go
//
// Оно
// https://github.com/SivWatt/formatter
//
//https://github.com/godep-migrator/logrus-formatters/blob/master/l2met/formatter.go

// TODO create general interface with generic fields

/*
func Debug(msg ...interface{}) {
	logrus.Debug(msg...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Info(msg ...interface{}) {
	logrus.Info(msg...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Warn(msg ...interface{}) {
	logrus.Warn(msg...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Error(msg ...interface{}) {
	logrus.Error(msg...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}
*/
