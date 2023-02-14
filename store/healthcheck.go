package store

import (
	"context"

	"github.com/pkg/errors"

	"github.com/floorstyle/fifo-queue/config"
	"github.com/floorstyle/fifo-queue/util"
)

type HealthCheckReal struct {
	Retries     int64
	MaxRetries  int64
	isAvailable bool
	RestoredCh  chan bool
}

func NewHealthCheck(config *config.Configuration) *HealthCheckReal {
	return &HealthCheckReal{
		Retries:     0,
		MaxRetries:  config.HealthCheckMaxRetries,
		isAvailable: true,
		RestoredCh:  make(chan bool),
	}
}

// Note: healthCheck counts retries only in series, otherwise it is refreshed
func (healthCheck *HealthCheckReal) Run(ctx context.Context, worker func(context.Context) error) error {
	if !healthCheck.IsHealthy() {
		return errors.New(util.HealthCheckNotAvailableErr.String())
	}
	err := worker(ctx)
	if util.IsErrConnection(err) {
		healthCheck.Retries++
		if healthCheck.Retries >= healthCheck.MaxRetries {
			healthCheck.isAvailable = false
		}
		return err
	}

	healthCheck.Refresh()
	return err
}

func (healthCheck HealthCheckReal) IsHealthy() bool {
	return healthCheck.isAvailable
}

func (healthCheck *HealthCheckReal) Refresh() {
	healthCheck.isAvailable = true
	healthCheck.Retries = 0
}

func (healthCheck *HealthCheckReal) RestoreChannel() {
	healthCheck.RestoredCh <- true
}

func (healthCheck *HealthCheckReal) Restored() chan bool {
	return healthCheck.RestoredCh
}
