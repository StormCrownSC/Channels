package main

import (
	"awesomeProject/3_taxis/taxi_api"
	"awesomeProject/3_taxis/taxis"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sort"
	"sync"
	"time"
)

func calcRide(a, b taxis.Point, taxiServices []taxi_api.ApiTaxiService) {
	var (
		resultChan = make(chan taxi_api.TaxiInfo, len(taxiServices))
		wg         sync.WaitGroup
		answers    []taxi_api.TaxiInfo
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, t := range taxiServices {
		wg.Add(1)
		go func(service taxi_api.ApiTaxiService) {
			defer wg.Done()
			taxiPrice, err := service.Calc(ctx, a, b)
			name := t.Name()
			if err != nil {
				taxiPrice = &taxis.Ride{Price: 0}
			}
			resultChan <- taxi_api.TaxiInfo{Name: name, Price: taxiPrice.Price, Error: err}
		}(t)
	}

	wg.Wait()
	close(resultChan)

	for result := range resultChan {
		if result.Error != nil {
			fmt.Println("Error: ", result.Error)
			continue
		}
		answers = append(answers, result)
	}

	// Сортировка результатов по цене
	sort.SliceStable(answers, func(i, j int) bool {
		return answers[i].Price < answers[j].Price
	})

	// Вывод результатов
	for _, answer := range answers {
		fmt.Println(answer)
	}
}

func calcRideWithErrGroup(a, b taxis.Point, taxiServices []taxi_api.ApiTaxiService) {
	var (
		resultChan = make(chan taxi_api.TaxiInfo, len(taxiServices))
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(len(taxiServices))

	for _, t := range taxiServices {
		t := t
		eg.Go(func() error {
			taxiPrice, err := t.Calc(ctx, a, b)
			name := t.Name()
			if err != nil {
				return err
			}
			resultChan <- taxi_api.TaxiInfo{Name: name, Price: taxiPrice.Price}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Println("wg err: ", err)
	} else {
		fmt.Println("All service completed successfully")
	}
	close(resultChan)
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
	taxiServices := []taxi_api.ApiTaxiService{
		taxi_api.NewMockBoltTaxi(),
		taxi_api.NewMockMaximTaxi(),
		taxi_api.NewMockGettTaxi(),
		taxi_api.NewMockYandexTaxi(),
	}
	calcRide(p1, p2, taxiServices)
	fmt.Printf("\n\nWait group:\n")
	calcRideWithErrGroup(p1, p2, taxiServices)
}
