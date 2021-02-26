package domain

import (
	"context"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
)

type WebrtcStorage interface {
}

type IonService interface {

}

type ConfigService interface {
	Get(ctx context.Context) (*kit.Config, error)
}