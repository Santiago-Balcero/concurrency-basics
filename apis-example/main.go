package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Status string

const (
	ok   Status = "OK"
	fail Status = "FAIL"
)

type Scores struct {
	overall      float64
	status       Status
	spoofAttempt bool
}

type UserData struct {
	userID    int
	userName  string
	userEmail string
}

type Request struct {
	uuid     string
	token    string
	isFraud  bool
	scores   Scores
	userData UserData
}

func getRandomMilliseconds() int {
	min := 500
	max := 3000
	return rand.Intn(max-min+1) + min
}

func getToken(uuid string) string {
	randNum := getRandomMilliseconds()
	time.Sleep(time.Duration(randNum) * time.Millisecond)
	fmt.Printf("%s - calling API token took %d milliseconds\n", uuid, randNum)
	return "token-1234"
}

func doGovernmentValidation(uuid string, wg *sync.WaitGroup) {
	randNum := getRandomMilliseconds()
	time.Sleep(time.Duration(randNum) * time.Millisecond)
	fmt.Printf("%s - calling API government validation took %d milliseconds\n", uuid, randNum)
	wg.Done()
}

func doBankValidation(uuid string, wg *sync.WaitGroup) {
	randNum := getRandomMilliseconds()
	time.Sleep(time.Duration(randNum) * time.Millisecond)
	fmt.Printf("%s - calling API bank validation took %d milliseconds\n", uuid, randNum)
	wg.Done()
}

func fetchUserData(uuid string, userID int, wg *sync.WaitGroup, datach chan UserData) {
	randNum := getRandomMilliseconds()
	time.Sleep(time.Duration(randNum) * time.Millisecond)
	fmt.Printf("%s - calling API user data took %d milliseconds\n", uuid, randNum)
	datach <- UserData{
		userID:    userID,
		userName:  "John Doe",
		userEmail: "john@mail.com",
	}
	wg.Done()
}

func fetchUserScores(uuid string, wg *sync.WaitGroup, scoresch chan Scores) {
	randNum := getRandomMilliseconds()
	time.Sleep(time.Duration(randNum) * time.Millisecond)
	fmt.Printf("%s - calling API scores took %d milliseconds\n", uuid, randNum)
	scoresch <- Scores{
		overall:      94.8,
		status:       ok,
		spoofAttempt: false,
	}
	wg.Done()
}

func main() {
	fmt.Println("APIs with concurrency!")

	now := time.Now()

	req := Request{
		uuid:    "1234",
		isFraud: false,
		scores:  Scores{},
		userData: UserData{
			userID: 10,
		},
	}

	req.token = getToken(req.uuid)

	wg := &sync.WaitGroup{}

	go doGovernmentValidation(req.uuid, wg)
	go doBankValidation(req.uuid, wg)

	wg.Add(2)
	wg.Wait()
	fmt.Println("Government and bank validation done!")

	datach := make(chan UserData)
	scoresch := make(chan Scores)

	go fetchUserData(req.uuid, req.userData.userID, wg, datach)
	go fetchUserScores(req.uuid, wg, scoresch)
	wg.Add(2)
	req.userData = <-datach
	req.scores = <-scoresch
	wg.Wait()
	close(scoresch)
	close(datach)

	fmt.Println("User data and scores fetched!")
	fmt.Println("Request:", req)
	fmt.Println(time.Since(now))
}
