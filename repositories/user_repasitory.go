package repositories

import (
	"errors"

	"github.com/daalvand/go-api/models"
	"github.com/daalvand/go-api/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	Exists(email string) (bool, error)
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	ValidateCredentials(email, password string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func (ur *userRepository) Exists(email string) (bool, error) {
	var count int64
	if err := ur.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	exists, err := ur.Exists(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return ErrDuplicateEmail
	}

	return ur.db.Create(user).Error
}

func (ur *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) ValidateCredentials(email, password string) (*models.User, error) {
	var user models.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err // User not found
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, ErrIncorrectPassword
	}

	return &user, nil
}
