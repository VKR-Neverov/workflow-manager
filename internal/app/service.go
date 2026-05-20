package app

import (
	"context"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/enums"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/storage"
	desc "github.com/fidesy-pay/workflow-manager/pkg/workflow-manager"
	"google.golang.org/grpc"
)

type (
	Implementation struct {
		desc.UnimplementedWorkflowManagerServer

		storage *storage.Storage
		runner  Runner
	}

	Runner interface {
		Create(ctx context.Context, flowType enums.FlowType, options model.Options) (*model.Flow, error)
		Complete(ctx context.Context, flowID string, options model.Options) (*model.Flow, error)
	}
)

func New(
	storage *storage.Storage,
	runner Runner,
) *Implementation {
	return &Implementation{
		storage: storage,
		runner:  runner,
	}
}

func (i *Implementation) GetDescription() *grpc.ServiceDesc {
	return &desc.WorkflowManager_ServiceDesc
}
