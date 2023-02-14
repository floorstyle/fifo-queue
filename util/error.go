package util

import (
	"fmt"
	"regexp"
	"runtime/debug"

	"github.com/pkg/errors"
)

const (
	REDIS_ERROR_NOT_FOUND = "redis: nil"
	REDIS_PANIC_ERROR     = "panic: db-manager-msg-queue - redis error:"
	QUEUE_ERROR_IS_EMPTY  = "msg store is empty"
)

var (
	RedisErrConnectionRegexp   = regexp.MustCompile(`\bconnection\b`)
	HealthCheckNotAvailableErr = regexp.MustCompile(`not available`)
)

func Try(err error) error {
	if err != nil {
		panic(WithStack(err))
	}
	return nil
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

func LogErr(err error) error {
	if err != nil {
		fmt.Println(err)
		error, ok := errors.Cause(err).(StackTracer)
		if !ok {
			debug.PrintStack()
			return nil
		}

		frames := error.StackTrace()
		for _, stack := range frames {
			fmt.Printf("%+s:%d\n", stack, stack)
		}
	}
	return nil
}

func WithStack(err error) error {
	if !HasStack(err) {
		return errors.WithStack(err)
	}
	return err
}

// TODO: Remove
func HasStack(err error) bool {
	for {
		if err == nil {
			return false
		}

		// Hidden interface implemented by some types in "github.com/pkg/errors".
		_, ok := err.(interface{ StackTrace() errors.StackTrace })
		if ok {
			return true
		}

		cause := errors.Unwrap(err)
		if cause == err {
			return false
		}
		err = cause
	}
}

func IsErrorNotFound(err error) bool {
	return err != nil && (err.Error() == REDIS_ERROR_NOT_FOUND || err.Error() == QUEUE_ERROR_IS_EMPTY)
}

func IsErrConnection(err error) bool {
	return err != nil && RedisErrConnectionRegexp.MatchString(err.Error())
}

func LogRedisError(err error) {
	if err != nil {
		fmt.Println(REDIS_PANIC_ERROR, err)
	}
}

func AllowErrNotFound(err error) error {
	if IsErrorNotFound(err) {
		return nil
	}
	return err
}

func AddStack(err error) error {
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
