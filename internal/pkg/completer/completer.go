package completer

import (
	"context"
	"fmt"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/dto"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/enums"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	"github.com/fidesy/sdk/common/logger"
	"github.com/fidesy/sdk/common/postgres"
	"sync"
	"time"
)

type (
	Service struct {
		storage Storage
		runners map[enums.FlowType]Runner
	}

	Storage interface {
		CreateFlow(ctx context.Context, flow *model.Flow) (*model.Flow, error)
		ListFlows(ctx context.Context, params dto.ListFlowsParams) ([]*model.Flow, error)
	}

	Runner interface {
		Complete(ctx context.Context, flow *model.Flow, options model.Options) (*model.Flow, error)
		Create(ctx context.Context, options model.Options) (*model.Flow, error)
		Type() enums.FlowType
	}
)

func New(
	storage Storage,
	runners ...Runner,
) *Service {

	runnersMap := make(map[enums.FlowType]Runner, len(runners))
	for _, runner := range runners {
		runnersMap[runner.Type()] = runner
	}

	return &Service{
		storage: storage,
		runners: runnersMap,
	}
}

func (s *Service) Create(ctx context.Context, flowType enums.FlowType, options model.Options) (*model.Flow, error) {
	runner := s.runners[flowType]

	flow, err := runner.Create(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("runner.Create: %w", err)
	}

	flow.Type = flowType
	flow.State = enums.StateNew

	flow, err = s.storage.CreateFlow(ctx, flow)
	if err != nil {
		return nil, fmt.Errorf("storage.CreateFlow: %w", err)
	}

	return flow, nil
}

func (s *Service) Complete(ctx context.Context, flowID string, options model.Options) (*model.Flow, error) {
	flows, err := s.storage.ListFlows(ctx, dto.ListFlowsParams{
		Filter: dto.ListFlowsFilter{
			IDIn:       []string{flowID},
			StateNotIn: []enums.State{enums.StateCompleted, enums.StateFailed},
		},
		Pagination: postgres.NewPagination(1, 1),
	})
	if err != nil {
		return nil, fmt.Errorf("storage.ListFlows: %w", err)
	}

	if len(flows) == 0 {
		return nil, postgres.ErrNotFound
	}

	flow := flows[0]

	return s.runners[flow.Type].Complete(ctx, flow, options)
}

func (s *Service) CompleteWorker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(5 * time.Second):
			s.complete(ctx)
		}
	}
}

func (s *Service) complete(ctx context.Context) {
	flows, err := s.storage.ListFlows(ctx, dto.ListFlowsParams{
		Filter: dto.ListFlowsFilter{
			StateNotIn: []enums.State{enums.StateCompleted, enums.StateFailed},
		},
		Pagination: postgres.NewPagination(1, 100),
	})
	if err != nil {
		logger.Errorf("storage.ListFlows: %v", err)
		return
	}

	var wg sync.WaitGroup

	for _, flow := range flows {
		flow := flow

		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err = s.runners[flow.Type].Complete(ctx, flow, nil)
			if err != nil {
				logger.Errorf("runner.Complete: %v", err)
				return
			}
		}()
	}

	wg.Wait()
}
