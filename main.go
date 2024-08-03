package main

import (
	"awesomeProject/forecast"
	"awesomeProject/predict_models"
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

	go process(ctx, stopCh)

	select {
	case <-ctx.Done():
		fmt.Println("stopping server...")
		stopCh <- struct{}{}
		fmt.Println("http server is stopped")
	}
}

func process(ctx context.Context, stopCh <-chan struct{}) {
	requestsChan := forecast.RequestRandomGenerator(stopCh)

	var (
		response [2]chan forecast.ForecastPrediction = [2]chan forecast.ForecastPrediction{
			make(chan forecast.ForecastPrediction),
			make(chan forecast.ForecastPrediction),
		}
		change   bool
		previous forecast.ForecastPrediction
	)

	go processRequests(requestsChan, response)
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

func processRequests(requestsChan <-chan forecast.ForecastRequest, response [2]chan forecast.ForecastPrediction) {
	model1 := predict_models.NewModel1()
	model2 := predict_models.NewModel2()

	for request := range requestsChan {
		response[0] <- model1.Predict(request)
		response[1] <- model2.Predict(request)
	}
}

func switchChange(change bool) bool {
	return !change
}
