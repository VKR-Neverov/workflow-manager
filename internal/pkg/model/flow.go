package model

import (
	"github.com/fidesy-pay/workflow-manager/internal/pkg/enums"
	"github.com/google/uuid"
	"time"
)

type Flow struct {
	ID          uuid.UUID      `db:"id"`
	State       enums.State    `db:"state"`
	Type        enums.FlowType `db:"type"`
	CreatedAt   time.Time      `db:"created_at"`
	CompletedAt *time.Time     `db:"completed_at"`
	Data        interface{}    `db:"data"`
}

func (f *Flow) TableName() string {
	return "flows"
}
