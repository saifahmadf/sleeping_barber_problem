package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numChairs = 5
)

// Customer represents a customer
type Customer struct {
	id int
}

// BarberShop represents the barber shop
type BarberShop struct {
	customers    chan Customer
	waitingRoom  chan Customer
	barberChair  chan Customer
	barberAsleep bool
	mutex        sync.Mutex
}

// NewBarberShop creates a new BarberShop instance
func NewBarberShop() *BarberShop {
	return &BarberShop{
		customers:    make(chan Customer),
		waitingRoom:  make(chan Customer, numChairs),
		barberChair:  make(chan Customer),
		barberAsleep: true,
	}
}

// Barber represents the barber
func (b *BarberShop) Barber() {
	for {
		select {
		case customer := <-b.customers:
			b.mutex.Lock()
			b.barberAsleep = false
			b.mutex.Unlock()

			fmt.Printf("Barber is cutting hair of Customer %d\n", customer.id)
			time.Sleep(time.Second) // Simulating haircut

			fmt.Printf("Barber finished cutting hair of Customer %d\n", customer.id)
			b.barberChair <- customer

			b.mutex.Lock()
			b.barberAsleep = true
			b.mutex.Unlock()

		default:
			b.mutex.Lock()
			if b.barberAsleep {
				fmt.Println("Barber is sleeping.")
			}
			b.mutex.Unlock()
		}
	}
}

// CustomerGenerator represents the customer generator
func (b *BarberShop) CustomerGenerator() {
	for i := 1; ; i++ {
		customer := Customer{id: i}

		b.mutex.Lock()
		if len(b.waitingRoom) < numChairs {
			b.waitingRoom <- customer
			fmt.Printf("Customer %d arrived and is waiting\n", customer.id)
		} else {
			fmt.Printf("Customer %d arrived but the waiting room is full. Leaving.\n", customer.id)
		}
		b.mutex.Unlock()

		time.Sleep(time.Second * 2) // Simulating time between customer arrivals
	}
}

// BarberShopSimulation simulates the barber shop
func BarberShopSimulation() {
	barberShop := NewBarberShop()

	go barberShop.Barber()
	go barberShop.CustomerGenerator()

	// Infinite loop to simulate the passing of time
	for {
		select {
		case customer := <-barberShop.waitingRoom:
			if barberShop.barberAsleep {
				fmt.Printf("Barber woke up by Customer %d\n", customer.id)
			}
			barberShop.customers <- customer
		case customer := <-barberShop.barberChair:
			fmt.Printf("Customer %d leaving with a new haircut\n", customer.id)
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func main() {
	BarberShopSimulation()
}
