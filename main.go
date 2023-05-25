package main

import (
	"log"
	"os"

	"iamkeychecker/api"
	"iamkeychecker/config"
	"iamkeychecker/pkg"
)

func main() {
	var awsProfile string
	const envVariable = "HOURS"
	keyUsageDurationHours := os.Getenv(envVariable)

	if keyUsageDurationHours == "" {
		log.Fatal("Environment variable HOURS is not set.")
	}

	args := os.Args
	if len(args) > 2 {
		awsProfile = args[2]
	} else {
		log.Println("Not enough arguments provided.")
		os.Exit(0)
	}

	iamClient := config.SetupIAMClient(awsProfile)
	log.Println("AWS Profile name provided:", awsProfile)
	if os.Args[1] == "server" {
		api.Run(iamClient, keyUsageDurationHours)

	} else if os.Args[1] == "cli" {
		if len(os.Args) > 3 {
			userName := os.Args[3]
			// Single User
			expired, err := pkg.ProcessSingleUser(os.Args[3], iamClient, keyUsageDurationHours)
			if err != nil {
				log.Println(err)
			} else if expired {
				log.Println("User", userName, "is using an expired key.")
			} else {
				log.Println("User", userName, "is not using an expired key.")
			}
		} else {
			// All Users
			expiredUsers, numberExpired, listErrs, err := pkg.ProcessingMultipleUsers(iamClient, keyUsageDurationHours)
			if err != nil {
				log.Println(err)
			}
			if numberExpired > 0 {
				log.Println("List of all the IAM Users using an expired key:", expiredUsers)
			} else {
				log.Println("No IAM Users are using an expired Access Key.")
			}

			if len(listErrs) > 0 {
				log.Println("Number of users skipped:", len(listErrs))
			}
		}
	} else {
		log.Println("Correct first arguments values are : server or cli.")
		os.Exit(1)
	}
}
