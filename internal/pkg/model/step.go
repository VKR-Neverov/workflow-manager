package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/dto"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/enums"
	"github.com/opentracing/opentracing-go"
	"time"
)

type Storage interface {
	ListFlows(ctx context.Context, params dto.ListFlowsParams) ([]*Flow, error)
	CreateFlow(ctx context.Context, flow *Flow) (*Flow, error)
	UpdateFlow(ctx context.Context, flow *Flow) (*Flow, error)
}

type StepFunc[Data any] func(ctx context.Context, input Step[Data]) (Step[Data], error)

type Step[Data any] struct {
	Options Options
	Flow    *Flow
	Data    Data
}

func (s Step[Data]) EmptyStep() Step[Data] {
	return s
}

func (s Step[Data]) Next(state enums.State) Step[Data] {
	s.Flow.State = state
	return s
}

type Stepper[Data any] struct {
	storage Storage
	steps   map[enums.State]StepFunc[Data]
}

func NewStepper[Data any](
	storage Storage,
) *Stepper[Data] {
	return &Stepper[Data]{
		storage: storage,
		steps:   make(map[enums.State]StepFunc[Data]),
	}
}

func (s *Stepper[Data]) Add(currentState enums.State, step StepFunc[Data]) {
	s.steps[currentState] = step
}

func (s *Stepper[Data]) Complete(ctx context.Context, flow *Flow, options Options) (*Flow, error) {
	for {
		currentState := flow.State

		flow, err := s.complete(ctx, flow, options)
		if err != nil {
			return nil, fmt.Errorf("s.complete: %w", err)
		}

		if flow.State == currentState {
			return flow, nil
		}
	}
}

func (s *Stepper[Data]) complete(ctx context.Context, flow *Flow, options Options) (*Flow, error) {
	var (
		err error
	)

	currentState := flow.State

	span, ctx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("complete: %s", currentState))
	defer span.Finish()

	if currentState == enums.StateCompleted {
		t := time.Now()
		flow.CompletedAt = &t

		flow, err = s.storage.UpdateFlow(ctx, flow)
		if err != nil {
			return nil, fmt.Errorf("storage.UpdateFlow: %w", err)
		}

		return flow, nil
	}

	currentData, err := getCurrentData[Data](flow)
	if err != nil {
		return nil, err
	}

	completeFunc := s.steps[currentState]
	if completeFunc == nil {
		return flow, fmt.Errorf("unknown state: %v", flow.State)
	}

	if options == nil {
		options = NewOptions()
	}

	var step = Step[Data]{
		Flow:    flow,
		Data:    currentData,
		Options: options,
	}
	step, err = completeFunc(ctx, step)
	if err != nil {
		return flow, err
	}

	flow = step.Flow
	flow.Data = step.Data

	flow, err = s.storage.UpdateFlow(ctx, flow)
	if err != nil {
		return nil, fmt.Errorf("storage.UpdateFlow: %w", err)
	}

	return flow, nil
}

func getCurrentData[Data any](flow *Flow) (Data, error) {
	var data Data

	bytes, _ := json.Marshal(flow.Data)

	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
