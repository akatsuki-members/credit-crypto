package log

import (
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

//This are the log levels
const (
	Trace Level = iota - 1
	Debug
	Info
	Warn
	Error
	Fatal
	Panic
)

// func to validate valid log level
func IsValidLogLevel(l Level) bool {
	return l >= Trace && l <= Panic
}

// Logger is a simple type used to manage
// the backing zerolog.Logger
type Logger struct {
	logger *zerolog.Logger
}

// Level defines log levels.
type Level int8

// New creates and returns a new logger
// using an zerolog.Logger
func New(l Level, out io.Writer) (*Logger, error) {
	if !IsValidLogLevel(l) {
		return nil, errors.New("invalid log level")
	}
	//format time like  "YYYY-MM-DDT24H:MM:SS.000Z"
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.Level(l))
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	// wrap the zerolog instance
	zlogger := zerolog.New(out).With().Stack().Logger()
	return &Logger{
		&zlogger,
	}, nil
}

// Trace logs a trace event.
func (l *Logger) Trace(msg string, args ...LogParam) {
	event := l.logger.Trace()
	eventsTypeAssertion(args, event)
	event.Msgf(msg)
}

// Debug logs a debug event.
func (l *Logger) Debug(msg string, args ...LogParam) {
	event := l.logger.Debug()
	eventsTypeAssertion(args, event)
	event.Msgf(msg)
}

// Info logs an info event.
func (l *Logger) Info(msg string, args ...LogParam) {
	event := l.logger.Info()
	eventsTypeAssertion(args, event)
	event.Msgf(msg)
}

// Warn logs a warn event.
func (l *Logger) Warn(msg string, args ...LogParam) {
	event := l.logger.Warn()
	eventsTypeAssertion(args, event)
	event.Msgf(msg)
}

// Error logs an error event.
// write the error stack trace
func (l *Logger) Error(err error, msg string, args ...LogParam) {
	err = errors.Wrap(err, "err context")
	event := l.logger.Error()
	eventsTypeAssertion(args, event)
	event.Stack().Err(err).Msgf(msg)
}

// Fatal logs a fatal event and exits the process.
func (l *Logger) Fatal(err error, msg string, args ...LogParam) {
	err = errors.Wrap(err, "err context")
	event := l.logger.Fatal()
	eventsTypeAssertion(args, event)
	event.Stack().Err(err).Msgf(msg)
}

// Panic logs a panic event and panics.
func (l *Logger) Panic(err error, msg string, args ...LogParam) {
	err = errors.Wrap(err, "err context")
	event := l.logger.Panic()
	eventsTypeAssertion(args, event)
	event.Stack().Err(err).Msgf(msg)
}

// eventsTypeAssertion is a helper function to assert the type of the event
func eventsTypeAssertion(args []LogParam, event *zerolog.Event) {
	for _, v := range args {
		switch x := v.Value.(type) {
		case string:
			event.Str(v.Key, x)
		case int:
			event.Int(v.Key, x)
		case int8:
			event.Int8(v.Key, x)
		case int16:
			event.Int16(v.Key, x)
		case int32:
			event.Int32(v.Key, x)
		case int64:
			event.Int64(v.Key, x)
		case uint:
			event.Uint(v.Key, x)
		case uint8:
			event.Uint8(v.Key, x)
		case uint16:
			event.Uint16(v.Key, x)
		case uint32:
			event.Uint32(v.Key, x)
		case uint64:
			event.Uint64(v.Key, x)
		case float32:
			event.Float32(v.Key, x)
		case float64:
			event.Float64(v.Key, x)
		case bool:
			event.Bool(v.Key, x)
		case time.Time:
			event.Time(v.Key, x)
		default:
			event.Interface(v.Key, x)
		}
	}
}
