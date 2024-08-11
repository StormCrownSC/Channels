// файл не менять!
package mock

import (
	"awesomeProject/3_taxis/taxis"
	"errors"
	"time"
)

type MockTaxiService struct {
	sleepTime time.Duration
	err       error
	price     int
	name      string
}

func (m *MockTaxiService) Name() string {
	return m.name
}

func (m *MockTaxiService) Calc(a, b taxis.Point) (*taxis.Ride, error) {
	if m.err == nil {
		return &taxis.Ride{Price: m.price}, nil
	}
	return nil, m.err
}

func NewMockYandexTaxi() *MockTaxiService {
	return &MockTaxiService{
		sleepTime: 0,
		err:       nil,
		price:     100,
		name:      "Yandex",
	}
}

func NewMockGettTaxi() *MockTaxiService {
	return &MockTaxiService{
		sleepTime: time.Millisecond * 100,
		err:       nil,
		price:     200,
		name:      "Gett",
	}
}

func NewMockBoltTaxi() *MockTaxiService {
	return &MockTaxiService{
		sleepTime: time.Second * 5,
		err:       nil,
		price:     200,
		name:      "Bolt",
	}
}

func NewMockMaximTaxi() *MockTaxiService {
	return &MockTaxiService{
		sleepTime: 0,
		err:       errors.New("Maxim taxi 500 error"),
		price:     200,
		name:      "Maxim",
	}
}
