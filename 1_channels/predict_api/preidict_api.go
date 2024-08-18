package predict_api

import (
	"awesomeProject/1_channels/forecast"
	"awesomeProject/1_channels/predict_models"
	"time"
)

var _ forecast.ForecastPredictor = (*Model)(nil)

type Model struct {
	model  *predict_models.Model
	reqMax int
	quotas <-chan struct{}
}

const RateLimitPeriod = time.Minute

func createQuotas(stopChan <-chan struct{}, rateLimit int) <-chan struct{} {
	quotas := make(chan struct{}, rateLimit)

	for i := 0; i < rateLimit; i++ {
		quotas <- struct{}{}
	}

	go func() {
		tick := time.NewTicker(RateLimitPeriod / time.Duration(rateLimit))
		defer tick.Stop()
		for {
			select {
			case <-stopChan:
				close(quotas)
				return
			case _ = <-tick.C:
				quotas <- struct{}{}
			}
		}
	}()
	return quotas
}

func (m *Model) Predict(req forecast.ForecastRequest) forecast.ForecastPrediction {
	_ = <-m.quotas
	return m.model.Predict(req)
}

func NewModel1(stopChan <-chan struct{}) *Model {
	rateLimit := 60
	quotas := createQuotas(stopChan, rateLimit)
	model := predict_models.NewModel1()
	return &Model{
		model:  model,
		reqMax: rateLimit,
		quotas: quotas,
	}
}
func NewModel2(stopChan <-chan struct{}) *Model {
	rateLimit := 90
	quotas := createQuotas(stopChan, rateLimit)
	model := predict_models.NewModel2()
	return &Model{
		model:  model,
		reqMax: rateLimit,
		quotas: quotas,
	}
}
