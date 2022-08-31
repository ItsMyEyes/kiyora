package utils

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

func GetSpanId(span opentracing.Span) string {
	uid := fmt.Sprintf("%v", uuid.New())
	id := strings.Split(fmt.Sprintf("%s", span.Context()), ":")
	if len(id) > 0 {
		return id[0]
	}
	return uid
}

func NewGetSpanID() string {
	uid := fmt.Sprintf("%v", uuid.New())
	return strings.ReplaceAll(uid, "-", "")
}
