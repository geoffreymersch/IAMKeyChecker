package api

import (
	"iamkeychecker/pkg"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/gin-gonic/gin"
)

func Run(client *iam.Client, keyUsageDurationHours string) {
	r := gin.Default()

	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		expired, err := pkg.ProcessSingleUser(name, client, keyUsageDurationHours)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"iam_user_name":        name,
				"is_using_expired_key": expired,
			})
		}
	})

	r.GET("/users/", func(c *gin.Context) {
		expiredUsers, numberExpired, listErrs, err := pkg.ProcessingMultipleUsers(client, keyUsageDurationHours)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else if len(listErrs) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": listErrs,
			})
		} else if numberExpired > 0 {
			c.JSON(http.StatusOK, gin.H{
				"number_users_expired_key": numberExpired,
				"list_users_expired_key":   expiredUsers,
			})
		}

	})

	err := r.Run()
	if err != nil {
		log.Println(err)
	}

}
