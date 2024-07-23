package main

// Producer consumer pattern
// https://medium.com/@nirajranasinghe/design-patterns-for-concurrent-programming-producer-consumer-pattern-39193cac195a

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type PaymentStatus string

const (
	pending  PaymentStatus = "PENDING"
	approved PaymentStatus = "APPROVED"
	rejected PaymentStatus = "REJECTED"
)

type PaymentMethod string

const (
	creditCard PaymentMethod = "CREDIT_CARD"
	giftCard   PaymentMethod = "GIFT_CARD"
	ach        PaymentMethod = "ACH"
)

type Payment struct {
	paymentID     int
	paymentMethod PaymentMethod
	paymentStatus PaymentStatus
	updatedAt     time.Time
}

func getRandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func getPendingPayments() []Payment {
	fmt.Println("Getting pending payments...")
	now := time.Now()
	// Preallocates a slice of random length with pending payments
	// to avoid multiple allocations
	payments := make([]Payment, 1000)
	for i := range payments {
		payment := Payment{
			paymentID:     i,
			paymentStatus: pending,
			updatedAt:     time.Now(),
		}
		method := getRandomInt(1, 3)
		if method == 1 {
			payment.paymentMethod = creditCard
		} else if method == 2 {
			payment.paymentMethod = giftCard
		} else {
			payment.paymentMethod = ach
		}
		payments[i] = payment
	}
	fmt.Printf("Got %d pending payments\n", len(payments))
	fmt.Println("Getting pending payments took:", time.Since(now))
	return payments
}

func processPayment(payment Payment, paych chan<- Payment, wg *sync.WaitGroup) {
	time.Sleep(200 * time.Millisecond)
	newStatus := getRandomInt(1, 2)
	if newStatus == 1 {
		payment.paymentStatus = approved
	} else {
		payment.paymentStatus = rejected
	}
	paych <- payment
	wg.Done()
}

func processPayments(payments []Payment, paych chan<- Payment) {
	now := time.Now()
	paywg := &sync.WaitGroup{}
	for _, payment := range payments {
		paywg.Add(1)
		go processPayment(payment, paych, paywg)
	}
	paywg.Wait()
	close(paych)
	fmt.Println("Processing payments took:", time.Since(now))
}

func sendNotification() {
	time.Sleep(200 * time.Millisecond)
}

func sendNotifications(paych <-chan Payment) {
	now := time.Now()
	counter := 0
	for range paych {
		go sendNotification()
		counter++
	}
	fmt.Printf("Sending %d notifications took: %v\n", counter, time.Since(now))
}

func main() {
	fmt.Println("Batch payments example!")
	now := time.Now()
	payments := getPendingPayments()
	// Buffered channel to allow non-blocking sends
	paych := make(chan Payment, len(payments))
	go processPayments(payments, paych)
	sendNotifications(paych)
	fmt.Println("Total time:", time.Since(now))
}
