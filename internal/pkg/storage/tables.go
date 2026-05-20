package storage

import (
	"fmt"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	"reflect"
	"strings"
)

var (
	flowsTable = (&model.Flow{}).TableName()
	flowFields = modelColumns(&model.Flow{})
)

type Model interface {
	TableName() string
}

func modelColumns(m Model) string {
	dbTags := make([]string, 0)

	dataType := reflect.TypeOf(m)

	if dataType.Kind() != reflect.Pointer {
		return ""
	}

	data := dataType.Elem()

	for i := 0; i < data.NumField(); i++ {
		field := data.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		dbTags = append(dbTags, fmt.Sprintf("%s.%s", m.TableName(), dbTag))
	}

	return strings.Join(dbTags, ",")
}
