package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	rand.Seed(time.Now().Unix())

	startTime := time.Now()

	usersCount, workerCount := 100, 10
	wg := &sync.WaitGroup{}

	jobs := make(chan User, usersCount)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go saveUserInfo(jobs, wg)
	}

	generateUsers(usersCount, jobs)

	wg.Wait()

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func saveUserInfo(jobs <-chan User, wg *sync.WaitGroup) {
	for user := range jobs {
		fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

		filename := fmt.Sprintf("users/uid%d.txt", user.id)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}

		file.WriteString(user.getActivityInfo())
		file.Close()
		time.Sleep(time.Second)

	}
	wg.Done()
}

func generateUsers(count int, jobs chan<- User) {

	for i := 0; i < count; i++ {
		jobs <- User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@company.com", i+1),
			logs:  generateLogs(rand.Intn(1000)),
		}
		fmt.Printf("generated users %d\n", i+1)
		time.Sleep(time.Millisecond * 100)
	}
	close(jobs)
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
