package main

import (
	"awesomeProject/1_channels/forecast"
	"awesomeProject/1_channels/predict_api"
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	stopCh := make(chan struct{})

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	go printer(stopCh)

	select {
	case <-ctx.Done():
		fmt.Println("stopping server...")
		for i := 0; i < 3; i++ {
			stopCh <- struct{}{}
		}
		fmt.Println("server is stopped")
	}
}

func printer(stopCh <-chan struct{}) {
	requestsChan := forecast.RequestRandomGenerator(stopCh)
	output := composer(requestsChan, stopCh)
	for out := range output {
		fmt.Println(out)
	}
}

func composer(requestsChan <-chan forecast.ForecastRequest, stopChan <-chan struct{}) (response chan forecast.ForecastPrediction) {
	model1 := predict_api.NewModel1(stopChan)
	model2 := predict_api.NewModel2(stopChan)

	go func(model1, model2 *predict_api.Model) {
		var model1Response, model2Response forecast.ForecastPrediction
		wg := sync.WaitGroup{}

		for request := range requestsChan {
			wg.Add(2)
			go func() {
				defer wg.Done()
				model1Response = model1.Predict(request)
			}()
			go func() {
				defer wg.Done()
				model2Response = model2.Predict(request)
			}()
			wg.Wait()

			if model1Response.ProbabilityPercent > model2Response.ProbabilityPercent {
				fmt.Println(model1Response)
			} else {
				fmt.Println(model2Response)
			}
		}
		close(response)
	}(model1, model2)
	return
}
