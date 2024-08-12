package main

import (
	"awesomeProject/1_channels/forecast"
	"awesomeProject/1_channels/predict_models"
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
		stopCh <- struct{}{}
		fmt.Println("http server is stopped")
	}
}

func printer(stopCh <-chan struct{}) {
	requestsChan := forecast.RequestRandomGenerator(stopCh)
	output := composer(requestsChan)
	for out := range output {
		fmt.Println(out)
	}
}

func composer(requestsChan <-chan forecast.ForecastRequest) (response chan forecast.ForecastPrediction) {
	model1 := predict_models.NewModel1()
	model2 := predict_models.NewModel2()

	go func(model1, model2 *predict_models.Model) {
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
