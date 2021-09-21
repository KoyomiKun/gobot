package logger

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type ContextKeyType string

const (
	RequestID ContextKeyType = "RequestID"
)

func GetUUID() string {
	return uuid.NewV4().String()
}

func CtxSprintf(ctx context.Context, format string, v ...interface{}) string {
	if ctx != nil {
		var fmtBaseString string
		if requestID := ctx.Value(RequestID); requestID != nil {
			fmtBaseString = fmt.Sprintf("[RequestID: %s]", requestID.(string))
		}
		return fmt.Sprintf("%s %s", fmtBaseString, fmt.Sprintf(format, v...))
	}
	return fmt.Sprintf(format, v...)
}
