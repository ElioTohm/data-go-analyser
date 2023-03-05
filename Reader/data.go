package Reader

import (
	"strconv"

	"github.com/ztrue/tracerr"
)

type Customer struct {
	ID                int
	FirstName         string
	LastName          string
	CustomerReference string
	Status            string
	Orders            map[string]*Order
}

type Order struct {
	ID                int
	CustomerReference string
	OrderStatus       string
	OrderReference    string
	OrderTimestamp    int
	Items             []*Item
}

type Item struct {
	ID             int
	OrderReference string
	ItemName       string
	Quantity       int
	TotalPrice     float32
}

func ParseCustomer(line []string) *Customer {
	id, err := strconv.Atoi(line[0])
	if err != nil {
		tracerr.PrintSourceColor(err)
	}
	return &Customer{
		ID:                id,
		FirstName:         line[1],
		LastName:          line[2],
		CustomerReference: line[3],
		Status:            line[4],
	}
}

func ParseOrders(line []string) *Order {
	id, err := strconv.Atoi(line[0])
	if err != nil {
		tracerr.PrintSourceColor(err)
	}
	order_timestamp, err := strconv.Atoi(line[4])
	if err != nil {
		tracerr.PrintSourceColor(err)
	}
	return &Order{
		ID:                id,
		CustomerReference: line[1],
		OrderStatus:       line[2],
		OrderReference:    line[3],
		OrderTimestamp:    order_timestamp,
	}
}

func ParseItems(line []string) *Item {
	id, err := strconv.Atoi(line[0])
	if err != nil {
		tracerr.PrintSourceColor(err)
	}
	quantity, err := strconv.Atoi(line[3])
	if err != nil {
		tracerr.PrintSourceColor(err)
	}
	total_price, err := strconv.ParseFloat(line[4], 32)
	if err != nil {
		tracerr.PrintSourceColor(err)
	}
	return &Item{
		ID:             id,
		OrderReference: line[1],
		ItemName:       line[2],
		Quantity:       quantity,
		TotalPrice:     float32(total_price),
	}
}
