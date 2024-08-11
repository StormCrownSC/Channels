package main

import (
	"awesomeProject/3_taxis/mock"
	"awesomeProject/3_taxis/taxis"
	"fmt"
	"sort"
	"sync"
)

type TaxiService interface {
	Name() string
	Calc(a, b taxis.Point) (*taxis.Ride, error)
}

func calcRide(a, b taxis.Point, taxiServices []TaxiService) {
	var wg sync.WaitGroup
	var answers = []taxis.Ride{}
	for _, t := range taxiServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			taxiPrice, err := t.Calc(a, b)
			if err != nil {
				return
			}
			answers = append(answers, *taxiPrice)
		}()
	}
	wg.Wait()
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
}
