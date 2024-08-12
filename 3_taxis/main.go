package main

import (
	"awesomeProject/3_taxis/mock"
	"awesomeProject/3_taxis/taxis"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sort"
	"sync"
	"time"
)

type TaxiService interface {
	Name() string
	Calc(a, b taxis.Point) (*taxis.Ride, error)
}

func calcRide(a, b taxis.Point, taxiServices []TaxiService) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var answers []taxis.Ride
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, t := range taxiServices {
		wg.Add(1)
		go func(service TaxiService) {
			defer wg.Done()

			// Канал для получения результата или ошибки
			resultChan := make(chan *taxis.Ride)
			errorChan := make(chan error)

			go func() {
				taxiPrice, err := service.Calc(a, b)
				if err != nil {
					errorChan <- err
					return
				}
				resultChan <- taxiPrice
			}()

			select {
			case <-ctx.Done():
				fmt.Println("The waiting time has expired")
			case taxiPrice := <-resultChan:
				mu.Lock()
				answers = append(answers, *taxiPrice)
				mu.Unlock()
			case err := <-errorChan:
				fmt.Println("Error from service:", err)
			}
			close(resultChan)
			close(errorChan)
		}(t)
	}
	wg.Wait()

	// Сортировка результатов по цене
	sort.SliceStable(answers, func(i, j int) bool {
		return answers[i].Price < answers[j].Price
	})

	// Вывод результатов
	for _, answer := range answers {
		fmt.Println(answer)
	}
}

func calcRideWithErrGroup(a, b taxis.Point, taxiServices []TaxiService) {
	var answers []taxis.Ride
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	mu := &sync.Mutex{}

	for _, t := range taxiServices {
		service := t
		go eg.Go(func() error {
			taxiPrice, err := service.Calc(a, b)
			if err != nil {
				return err
			}
			mu.Lock()
			answers = append(answers, *taxiPrice)
			mu.Unlock()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	}

	sort.SliceStable(answers, func(i, j int) bool {
		return answers[i].Price < answers[j].Price
	})

	for _, answer := range answers {
		fmt.Println(answer)
	}
}

func main() {
	p1 := taxis.Point{
		Lat: 55.664398,
		Lon: 37.518636,
	}
	p2 := taxis.Point{
		Lat: 55.650810,
		Lon: 37.496776,
	}
	taxiServices := []TaxiService{
		mock.NewMockBoltTaxi(),
		mock.NewMockMaximTaxi(),
		mock.NewMockGettTaxi(),
		mock.NewMockYandexTaxi(),
	}
	calcRide(p1, p2, taxiServices)
	calcRideWithErrGroup(p1, p2, taxiServices)
}
