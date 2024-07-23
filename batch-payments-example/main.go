package main

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
	payments := []Payment{}
	for i := 0; i < getRandomInt(1000, 1005); i++ {
		payment := Payment{
			paymentID:     i,
			paymentMethod: PaymentMethod(getRandomInt(1, 3)),
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
		payments = append(payments, payment)
	}
	fmt.Printf("Got %d pending payments\n", len(payments))
	fmt.Println("Getting pending payments took:", time.Since(now))
	return payments
}

func processPayment(payment Payment, paych chan Payment, wg *sync.WaitGroup) {
	time.Sleep(time.Duration(getRandomInt(50, 300)) * time.Millisecond)
	newStatus := getRandomInt(1, 2)
	if newStatus == 1 {
		payment.paymentStatus = approved
	} else {
		payment.paymentStatus = rejected
	}
	paych <- payment
	wg.Done()
}

func processPayments(payments []Payment, paych chan Payment, wg *sync.WaitGroup) {
	now := time.Now()
	paywg := &sync.WaitGroup{}
	for _, payment := range payments {
		paywg.Add(1)
		go processPayment(payment, paych, paywg)
	}
	paywg.Wait()
	close(paych)
	fmt.Println("Processing payments took:", time.Since(now))
	wg.Done()
}

func sendNotifications(paych chan Payment, wg *sync.WaitGroup) {
	now := time.Now()
	counter := 0
	for range paych {
		time.Sleep(time.Duration(getRandomInt(50, 300)) * time.Millisecond)
		counter++
	}
	fmt.Printf("Sending %d notifications took: %v\n", counter, time.Since(now))
	wg.Done()
}

func main() {
	fmt.Println("Batch payments example!")
	now := time.Now()
	payments := getPendingPayments()
	paych := make(chan Payment)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go processPayments(payments, paych, wg)
	go sendNotifications(paych, wg)
	wg.Wait()
	fmt.Println("Total time:", time.Since(now))
}
