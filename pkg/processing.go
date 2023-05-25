package pkg

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func ProcessSingleUser(userName string, client *iam.Client, keyUsageDurationHours string) (bool, error) {
	expired, err := SingleUserKeyInfo(context.Background(), client, userName, keyUsageDurationHours)

	if err != nil {
		return false, err
	} else {
		return expired, nil
	}
}

func ProcessingMultipleUsers(client *iam.Client, keyUsageDurationHours string) ([]string, int, []string, error) {
	log.Println("Fetching all the IAM Users.")

	users, err := ListAllUsers(context.Background(), client)
	if err != nil {
		log.Println(err)
	}

	log.Println("All the IAM Users fetched :", len(users))

	var wg sync.WaitGroup
	wg.Add(len(users))

	resultChan := make(chan User, len(users))

	// Keep 100 Goroutines at a time to throttle API requests
	guard := make(chan struct{}, 100)

	var listErr []string

	for _, userName := range users {
		guard <- struct{}{}

		go func(userName string) {
			defer func() {
				wg.Done()
				// To avoid 'failed to get rate limit token' error when querying a lot of users simultaneously.
				time.Sleep(1000 * time.Millisecond)
				<-guard
			}()
			err := MutipleUsersInfo(context.Background(), client, userName, resultChan, keyUsageDurationHours)
			if err != nil {
				listErr = append(listErr, err.Error())
			}
		}(userName)
	}
	wg.Wait()
	close(resultChan)

	listExpired := []string{}

	for user := range resultChan {
		for _, expired := range user.AccessKeys {
			if expired {
				listExpired = append(listExpired, user.UserName)
			}
			break
		}
	}

	return listExpired, len(listExpired), listErr, nil
}
