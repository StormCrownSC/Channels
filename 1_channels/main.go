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

	go filter(stopCh)

	select {
	case <-ctx.Done():
		fmt.Println("stopping server...")
		stopCh <- struct{}{}
		fmt.Println("http server is stopped")
	}
}

func filter(stopCh <-chan struct{}) {
	requestsChan := forecast.RequestRandomGenerator(stopCh)

	var (
		response [2]chan forecast.ForecastPrediction = [2]chan forecast.ForecastPrediction{
			make(chan forecast.ForecastPrediction),
			make(chan forecast.ForecastPrediction),
		}
		change   bool
		previous forecast.ForecastPrediction
	)

	go composer(requestsChan, response)
	output := aggregator(response)

	for out := range output {
		if change {
			if previous.ProbabilityPercent > out.ProbabilityPercent {
				fmt.Println(previous)
			} else {
				fmt.Println(out)
			}
			change = switchChange(change)
		} else {
			previous = out
			change = switchChange(change)
		}
	}

}

func aggregator(inputs [2]chan forecast.ForecastPrediction) <-chan forecast.ForecastPrediction {
	output := make(chan forecast.ForecastPrediction)
	var wg sync.WaitGroup
	for _, in := range inputs {
		wg.Add(1)
		go func(int <-chan forecast.ForecastPrediction) {
			defer wg.Done()
			for x := range in {
				output <- x
			}
		}(in)
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}

func composer(requestsChan <-chan forecast.ForecastRequest, response [2]chan forecast.ForecastPrediction) {
	model1 := predict_models.NewModel1()
	model2 := predict_models.NewModel2()
	wg := sync.WaitGroup{}

	for request := range requestsChan {
		wg.Add(2)
		go func() {
			defer wg.Done()
			response[0] <- model1.Predict(request)
		}()
		go func() {
			defer wg.Done()
			response[1] <- model2.Predict(request)
		}()
		wg.Wait()
	}
}

func switchChange(change bool) bool {
	return !change
}
