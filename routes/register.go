package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	event, err := GetEventById(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not find event."})
		return
	}

	err = event.CheckIfUserIsRegistered(userId)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusConflict, gin.H{"message": "User already registered for event."})
		return
	}

	err = event.Register(userId)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register user for event."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User Registered!"})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	event, err := GetEventById(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not find event."})
		return
	}

	err = event.CheckIfUserIsRegistered(userId)
	if err == nil {
		fmt.Println(err)
		context.JSON(http.StatusConflict, gin.H{"message": "User isn't registered for event."})
		return
	}

	err = event.CancelRegistration(userId)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not cancel register user for event."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User Registeration Canceled!"})
}

func getRegisteredUsers(context *gin.Context) {
	event, err := GetEventById(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not find event."})
		return
	}

	registerUsers, err := event.CheckRegistrations()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusConflict, gin.H{"message": "User already registered for event."})
		return
	}

	if registerUsers == nil {
		context.JSON(http.StatusOK, gin.H{"message": "No registerations"})
		return
	}

	context.JSON(http.StatusOK, registerUsers)
}
