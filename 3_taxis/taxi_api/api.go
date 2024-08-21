package taxi_api

import (
	"awesomeProject/3_taxis/mock"
	"awesomeProject/3_taxis/taxis"
	"context"
	"fmt"
)

type TaxiInfo struct {
	Name  string
	Price int
	Error error
}

type TaxiService interface {
	Name() string
	Calc(a, b taxis.Point) (*taxis.Ride, error)
}

type ApiTaxiService interface {
	Name() string
	Calc(ctx context.Context, a, b taxis.Point) (*taxis.Ride, error)
}

type ApiTaxi struct {
	taxi *mock.MockTaxiService
}

func (api *ApiTaxi) Name() string {
	return api.taxi.Name()
}

func (api *ApiTaxi) Calc(ctx context.Context, a, b taxis.Point) (*taxis.Ride, error) {
	var resultChan chan TaxiInfo = make(chan TaxiInfo)
	go func() {
		taxiPrices, err := api.taxi.Calc(a, b)
		if err != nil {
			taxiPrices = &taxis.Ride{Price: 0}
		}
		resultChan <- TaxiInfo{
			Price: taxiPrices.Price,
			Error: err,
		}
	}()

	select {
	case res := <-resultChan:
		return &taxis.Ride{Price: res.Price}, res.Error
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled")
	}
}

func NewMockYandexTaxi() *ApiTaxi {
	yandex := mock.NewMockYandexTaxi()
	return &ApiTaxi{
		taxi: yandex,
	}
}

func NewMockGettTaxi() *ApiTaxi {
	gett := mock.NewMockGettTaxi()
	return &ApiTaxi{
		taxi: gett,
	}
}

func NewMockBoltTaxi() *ApiTaxi {
	bolt := mock.NewMockBoltTaxi()
	return &ApiTaxi{
		taxi: bolt,
	}
}

func NewMockMaximTaxi() *ApiTaxi {
	maxim := mock.NewMockMaximTaxi()
	return &ApiTaxi{
		taxi: maxim,
	}
}
