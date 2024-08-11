package predict_models

import (
	"awesomeProject/1_channels/forecast"
	"math/rand"
	"sync"
	"time"
)

var _ forecast.ForecastPredictor = (*Model)(nil)

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

type Model struct {
	reqMax   int
	reqCount int
	mu       sync.Mutex
}

func (m *Model) Predict(req forecast.ForecastRequest) forecast.ForecastPrediction {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reqCount++
	if m.reqCount > m.reqMax+5 {
		time.Sleep(time.Minute)
		m.reqCount = 0
	}
	return forecast.ForecastPrediction{
		Location:           req.Location,
		Time:               req.Time,
		TemperatureCelsius: randRange(5, 30),
		HumidityPercent:    randRange(30, 95),
		ProbabilityPercent: randRange(20, 100),
	}
}

func NewModel1() *Model {
	return &Model{
		reqMax: 60,
		mu:     sync.Mutex{},
	}
}
func NewModel2() *Model {
	return &Model{
		reqMax: 90,
		mu:     sync.Mutex{},
	}
}
