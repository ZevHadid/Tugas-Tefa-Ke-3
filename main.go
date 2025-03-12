package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type Item struct {
	Name  string
	Price float64
}

type Customer struct {
	ID    int
	Name  string
	Items []Item
}

var items = []Item{
	{"Apple", 1.99},
	{"Banana", 0.99},
	{"Orange", 1.49},
	{"Milk", 2.49},
	{"Bread", 1.89},
}

func handler(w http.ResponseWriter, r *http.Request) {
	localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	numCustomers := 100
	numCashiers := 5

	customers := make(chan Customer, numCustomers)
	var wg sync.WaitGroup

	for i := 1; i <= numCustomers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			customer := Customer{
				ID:   id,
				Name: gofakeit.Name(),
				Items: func() []Item {
					numItems := localRand.Intn(3) + 1
					var customerItems []Item
					for j := 0; j < numItems; j++ {
						customerItems = append(customerItems, items[localRand.Intn(len(items))])
					}
					return customerItems
				}(),
			}
			customers <- customer
		}(i)
	}

	for i := 1; i <= numCashiers; i++ {
		wg.Add(1)
		go func(cashierID int) {
			defer wg.Done()
			for customer := range customers {
				processTime := time.Duration(localRand.Intn(5)+1) * time.Second
				time.Sleep(processTime)
				total := 0.0
				for _, item := range customer.Items {
					total += item.Price
				}
				fmt.Fprintf(w, "Proses kasir %d {\nnama pelanggan: %d\nid pelanggan: %s\nharga total $%.2f\nwaktu: %v\n}\n\n", cashierID, customer.ID, customer.Name, total, processTime)
			}
		}(i)
	}

	wg.Wait()
	close(customers)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
