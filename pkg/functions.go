package pkg

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go/aws"
)

type ListAccessKeysAPI interface {
	ListAccessKeys(ctx context.Context, params *iam.ListAccessKeysInput, optFns ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error)
}

func ListAccessKeys(ctx context.Context, api ListAccessKeysAPI, input *iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error) {
	return api.ListAccessKeys(ctx, input)
}

func ListAllUsers(ctx context.Context, client *iam.Client) (list []string, err error) {
	var t []string
	iamInput := iam.ListUsersInput{}

	paginator := iam.NewListUsersPaginator(client, &iamInput, func(o *iam.ListUsersPaginatorOptions) {})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		for _, value := range output.Users {
			t = append(t, *value.UserName)
		}
	}

	return t, nil
}

type User struct {
	UserName   string
	AccessKeys map[string]bool
}

func MutipleUsersInfo(ctx context.Context, client *iam.Client, userName string, resultChan chan<- User, keyUsageDurationHours string) error {
	iamInput := iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	}

	user := User{
		UserName:   userName,
		AccessKeys: make(map[string]bool),
	}

	accessKeys, err := ListAccessKeys(context.Background(), client, &iamInput)
	if err != nil {
		log.Println(userName, err)
		return err
	}

	hourInt, _ := strconv.ParseFloat(keyUsageDurationHours, 64)

	for i := range accessKeys.AccessKeyMetadata {
		keyEpiration, _ := isExpired(accessKeys.AccessKeyMetadata[i].CreateDate, hourInt)
		key := *accessKeys.AccessKeyMetadata[i].AccessKeyId
		if keyEpiration {
			user.AccessKeys[key] = true
		} else {
			user.AccessKeys[key] = false
		}
	}

	resultChan <- user
	return nil
}

func SingleUserKeyInfo(ctx context.Context, client *iam.Client, userName string, keyUsageDurationHours string) (bool, error) {
	iamInput := iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	}

	user := User{
		UserName:   userName,
		AccessKeys: make(map[string]bool),
	}

	accessKeys, err := ListAccessKeys(context.Background(), client, &iamInput)
	if err != nil {
		return false, err
	}

	hourInt, _ := strconv.ParseFloat(keyUsageDurationHours, 64)

	for i := range accessKeys.AccessKeyMetadata {
		keyEpiration, _ := isExpired(accessKeys.AccessKeyMetadata[i].CreateDate, hourInt)
		key := *accessKeys.AccessKeyMetadata[i].AccessKeyId
		if keyEpiration {
			user.AccessKeys[key] = true
		} else {
			user.AccessKeys[key] = false
		}
	}

	for _, expired := range user.AccessKeys {
		if expired {
			return true, nil
		}
	}

	return false, nil
}

func isExpired(createAt *time.Time, keyUsageDurationHours float64) (bool, error) {
	currentTime := time.Now()

	creationDate := currentTime.Sub(*createAt)
	if creationDate.Hours() > keyUsageDurationHours {
		return true, nil
	}

	return false, nil
}
