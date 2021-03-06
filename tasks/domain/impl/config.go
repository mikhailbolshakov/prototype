package impl

import (
	"context"
	"fmt"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
)


func NewTaskConfigService() domain.ConfigService {
	return &taskConfigServiceImpl{}
}

type taskConfigServiceImpl struct{}

func (c *taskConfigServiceImpl) GetAll(ctx context.Context) []*domain.Config {
	return mockConfigs
}

func (c *taskConfigServiceImpl) Get(ctx context.Context, t *domain.Type) (*domain.Config, error) {

	for _, c := range mockConfigs {
		if c.Type.Equals(t) {
			return c, nil
		}
	}

	// load configuration from repository
	return nil, fmt.Errorf("config not found for %v", t)
}

func (c *taskConfigServiceImpl) NextTransitions(ctx context.Context, t *domain.Type, currentStatus *domain.Status) ([]*domain.Transition, error) {

	cfg, err := c.Get(ctx, t)
	if err != nil {
		return nil, err
	}

	var tr []*domain.Transition
	for _, c := range cfg.StatusModel.Transitions {
		if c.From.Equals(currentStatus) {
			tr = append(tr, c)
		}
	}

	return tr, nil
}

func (c *taskConfigServiceImpl) IsFinalStatus(ctx context.Context, t *domain.Type, s *domain.Status) bool {
	tr, _ := c.NextTransitions(ctx, t, s)
	return len(tr) == 0
}

func (c *taskConfigServiceImpl) FindTransition(ctx context.Context, t *domain.Type, current, target *domain.Status) (*domain.Transition, error) {

	cfg, err := c.Get(ctx, t)
	if err != nil {
		return nil, err
	}

	for _, c := range cfg.StatusModel.Transitions {
		if c.From.Equals(current) && c.To.Equals(target) {
			return c, nil
		}
	}

	return nil, fmt.Errorf("transition not found current %v, target %v", current, target)

}

func (c *taskConfigServiceImpl) InitialTransition(ctx context.Context, t *domain.Type) (*domain.Transition, error) {

	cfg, err := c.Get(ctx, t)
	if err != nil {
		return nil, err
	}

	for _, t := range cfg.StatusModel.Transitions {
		if t.Initial {
			return t, nil
		}
	}

	return nil, fmt.Errorf("cfg-error.initial transition not found for type %v", t)

}

