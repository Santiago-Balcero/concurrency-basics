package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Concurrency intro!")
	now := time.Now()
	userID := 10

	responseCh := make(chan string, 3)

	waitGroup := &sync.WaitGroup{}

	go fetchUserData(userID, responseCh, waitGroup)
	go fetchUserRecommendations(userID, responseCh, waitGroup)
	go fetchUserLikes(userID, responseCh, waitGroup)

	waitGroup.Add(3)
	waitGroup.Wait()

	close(responseCh)

	for resp := range responseCh {
		fmt.Println(resp)
	}

	fmt.Println(time.Since(now))
}

func fetchUserData(userID int, responseCh chan string, wg *sync.WaitGroup) {
	time.Sleep(80 * time.Millisecond)
	fmt.Println(userID, "calling API user data")
	responseCh <- "user data"
	wg.Done()
}

func fetchUserRecommendations(userID int, responseCh chan string, wg *sync.WaitGroup) {
	time.Sleep(120 * time.Millisecond)
	fmt.Println(userID, "calling API user recommendations")
	responseCh <- "user recommendations"
	wg.Done()
}

func fetchUserLikes(userID int, responseCh chan string, wg *sync.WaitGroup) {
	time.Sleep(50 * time.Millisecond)
	fmt.Println(userID, "calling API user likes")
	responseCh <- "user likes"
	wg.Done()
}
