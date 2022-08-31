package models

import (
	"context"
	"myself_framwork/library/logger/v2"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type GeneralModel struct {
	ParentSpan opentracing.Span
	Zaplog     *zap.Logger
	SpanID     string
	Context    context.Context
	Logging    logger.LoggingInterface
}

type DatabaseModel struct {
	ID        uint `gorm:"primarykey;column:id" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// CreatedBy string `gorm:"column:created_by" json:"created_by"`
	// UpdatedBy string `gorm:"column:updated_by" json:"updated_by"`
}
