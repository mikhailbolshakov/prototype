package domain

import (
	"errors"
	"fmt"
)

type TaskConfigService interface {
	// get whole configuration
	Get(t *Type) (*Config, error)
	// retrieves a list of available statuses which might be set for the task
	NextTransitions(t *Type, currentStatus *Status) ([]*Transition, error)
	// get initial transition
	InitialTransition(t *Type) (*Transition, error)
}

func NewTaskConfigService() TaskConfigService {
	return &TaskConfigServiceImpl{}
}

type TaskConfigServiceImpl struct{}

func (c *TaskConfigServiceImpl) Get(t *Type) (*Config, error) {

	for _, c := range mockConfigs {
		if c.Type.Type == t.Type && c.Type.SubType == t.SubType {
			return c, nil
		}
	}

	// load configuration from repository
	return nil, errors.New(fmt.Sprintf("config not found for %v", t))
}

func (c *TaskConfigServiceImpl) NextTransitions(t *Type, currentStatus *Status) ([]*Transition, error) {

	cfg, err := c.Get(t)
	if err != nil {
		return nil, err
	}

	var tr []*Transition
	for _, c := range cfg.StatusModel.Transitions {
		if c.From.Status == currentStatus.Status && c.From.SubStatus == currentStatus.SubStatus {
			tr = append(tr, c)
		}
	}

	return tr, nil
}

func (c *TaskConfigServiceImpl) InitialTransition(t *Type) (*Transition, error) {

	cfg, err := c.Get(t)
	if err != nil {
		return nil, err
	}

	for _, t := range cfg.StatusModel.Transitions {
		if t.Initial {
			return t, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("cfg-error.initial transition not found for type %v", t))

}

func (tr *Transition) checkGroup(group string) bool {
	for _, g := range tr.AllowAssignGroups {
		if g == group {
			return true
		}
	}
	return false
}
