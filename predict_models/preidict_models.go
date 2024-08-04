package predict_models

import (
	"awesomeProject/forecast"
	"math/rand"
	"time"
)

var _ forecast.ForecastPredictor = (*Model)(nil)

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

type Model struct {
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
	return forecast.ForecastPrediction{
		Location:           req.Location,
		Time:               req.Time,
		TemperatureCelsius: randRange(5, 30),
		HumidityPercent:    randRange(30, 95),
		ProbabilityPercent: randRange(20, 100),
	}
}

func NewModel1(stopChan <-chan struct{}) *Model {
	rateLimit := 60
	quotas := createQuotas(stopChan, rateLimit)
	return &Model{
		reqMax: rateLimit,
		quotas: quotas,
	}
}
func NewModel2(stopChan <-chan struct{}) *Model {
	rateLimit := 90
	quotas := createQuotas(stopChan, rateLimit)
	return &Model{
		reqMax: rateLimit,
		quotas: quotas,
	}
}
