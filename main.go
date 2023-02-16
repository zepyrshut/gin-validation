package main

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Login struct {
	User     string `json:"username" binding:"required,gt=0,email"`
	Password string `json:"password" binding:"required,gt=0"`
}

func ValidationWrapper(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)
	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs[f.Field()] = err
	}
	return errs
}

func ValidationWrapperResponse(err error) map[string]string {
	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		return ValidationWrapper(verr)
	}
	return nil
}

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.POST("/loginJSON", func(c *gin.Context) {
		var user Login

		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ValidationWrapperResponse(err)})
			return
		}

		fmt.Println(user.User)
		fmt.Println(user.Password)

		if user.User != "manu" || user.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	router.Run(":8080")
}
