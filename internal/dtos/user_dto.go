package dtos

import (
	"errors"
	"regexp"
	"time"
	"todo-app-mongo/internal/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserLoginDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponseDTO struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type UserRequestDTO struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

type UserResponseDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *UserRequestDTO) ToUserModel() (*entity.User, error) {

	if !u.validateName() {
		return nil, errors.New("name is required")
	}

	if !u.validateEmail() {
		return nil, errors.New("email is required")
	}

	if !u.validatePassword() {
		return nil, errors.New("password is required and must be at least 6 characters long")
	}

	hashedPassword, err := u.hashPassword()
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:             primitive.NewObjectID(),
		Name:           u.Name,
		Email:          u.Email,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (u *UserRequestDTO) validateName() bool {

	return u.Name != ""
}

func (u *UserRequestDTO) validateEmail() bool {

	if u.Email == "" {
		return false
	}

	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, u.Email)
	return match
}

func (u *UserRequestDTO) validatePassword() bool {

	if u.Password == "" {
		return false
	}

	if u.Password != u.ConfirmPassword {
		return false
	}

	if len(u.Password) < 6 {
		return false
	}

	return true
}

func (u *UserRequestDTO) hashPassword() (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
