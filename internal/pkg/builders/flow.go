package builders

import (
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	desc "github.com/fidesy-pay/workflow-manager/pkg/workflow-manager"
)

func BuildFlowResponse(flow *model.Flow) *desc.FlowResponse {
	return &desc.FlowResponse{
		Id:    flow.ID.String(),
		Type:  string(flow.Type),
		State: string(flow.State),
	}
}

func BuildFlowsResponse(flows []*model.Flow) []*desc.FlowResponse {
	result := make([]*desc.FlowResponse, len(flows))
	for i := range flows {
		result[i] = BuildFlowResponse(flows[i])
	}

	return result
}
