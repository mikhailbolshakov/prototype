package domain

import (
	"fmt"
)

type ConfigService interface {
	// get whole configuration
	Get(t *Type) (*Config, error)
	// if status is final
	IsFinalStatus(t *Type, s *Status) bool
	// retrieves a list of available statuses which might be set for the task
	NextTransitions(t *Type, currentStatus *Status) ([]*Transition, error)
	// get initial transition
	InitialTransition(t *Type) (*Transition, error)
	// get transition by source/target statuses
	FindTransition(t *Type, current, target *Status) (*Transition, error)
}

func NewTaskConfigService() ConfigService {
	return &taskConfigServiceImpl{}
}

type taskConfigServiceImpl struct{}

func (c *taskConfigServiceImpl) Get(t *Type) (*Config, error) {

	for _, c := range mockConfigs {
		if c.Type.equals(t) {
			return c, nil
		}
	}

	// load configuration from repository
	return nil, fmt.Errorf("config not found for %v", t)
}

func (c *taskConfigServiceImpl) NextTransitions(t *Type, currentStatus *Status) ([]*Transition, error) {

	cfg, err := c.Get(t)
	if err != nil {
		return nil, err
	}

	var tr []*Transition
	for _, c := range cfg.StatusModel.Transitions {
		if c.From.equals(currentStatus) {
			tr = append(tr, c)
		}
	}

	return tr, nil
}

func (c *taskConfigServiceImpl) IsFinalStatus(t *Type, s *Status) bool {
	tr, _ := c.NextTransitions(t, s)
	return len(tr) == 0
}

func (c *taskConfigServiceImpl) FindTransition(t *Type, current, target *Status) (*Transition, error) {

	cfg, err := c.Get(t)
	if err != nil {
		return nil, err
	}

	for _, c := range cfg.StatusModel.Transitions {
		if c.From.equals(current) && c.To.equals(target) {
			return c, nil
		}
	}

	return nil, fmt.Errorf("transition not found current %v, target %v", current, target)

}

func (c *taskConfigServiceImpl) InitialTransition(t *Type) (*Transition, error) {

	cfg, err := c.Get(t)
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

func (tr *Transition) checkGroup(group string) bool {
	for _, g := range tr.AllowAssignGroups {
		if g == group {
			return true
		}
	}
	return false
}
