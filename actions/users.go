package actions

import (
	"errors"
	"net/http"

	"github.com/daalvand/go-api/models"
	"github.com/daalvand/go-api/repositories"
	"github.com/daalvand/go-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserActions struct {
	repo repositories.UserRepository
}

func NewUserActions(repo repositories.UserRepository) *UserActions {
	return &UserActions{repo: repo}
}

func (ua *UserActions) Signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = ua.repo.Create(&user)

	if err != nil {
		if errors.Is(err, repositories.ErrDuplicateEmail) {
			context.JSON(http.StatusConflict, gin.H{"message": "Email already exists."})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user."})
		logrus.Error("Error on creating user", err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (ua *UserActions) Login(context *gin.Context) {
	var user *models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	user, err = ua.repo.ValidateCredentials(user.Email, user.Password)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Could not authenticate user."})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
		logrus.Error("Error on generating token", err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful!", "token": token})
}
