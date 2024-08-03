package forecast

import (
	"math/rand"
	"time"
)

type ForecastRequest struct {
	Location string
	Time     time.Time
}

type ForecastPrediction struct {
	Location           string
	Time               time.Time
	TemperatureCelsius int
	HumidityPercent    int
	ProbabilityPercent int
}

type ForecastPredictor interface {
	Predict(req ForecastRequest) ForecastPrediction
}

func RequestRandomGenerator(stopCh <-chan struct{}) <-chan ForecastRequest {
	c := make(chan ForecastRequest)
	locations := []string{"Moscow", "SaintPetersburg", "Kazan", "Nizhniy Novgorod", "Novosibirsk", "Samara"}
	go func() {
		for {
			select {
			case <-stopCh:
				close(c)
				return
			default:
				time.Sleep(100 * time.Millisecond)
				c <- ForecastRequest{
					Location: locations[rand.Intn(len(locations))],
					Time:     randate(),
				}
			}
		}
	}()
	return c
}

func randate() time.Time {
	minDate := time.Date(2025, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	maxDate := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := maxDate - minDate

	sec := rand.Int63n(delta) + minDate
	return time.Unix(sec, 0)
}
