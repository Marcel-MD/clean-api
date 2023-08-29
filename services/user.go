package services

import (
	"errors"

	"github.com/Marcel-MD/clean-api/auth"
	"github.com/Marcel-MD/clean-api/config"
	"github.com/Marcel-MD/clean-api/data/repositories"
	"github.com/Marcel-MD/clean-api/models"
	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindAll(query models.PaginationQuery) ([]models.User, error)
	FindById(id string) (models.User, error)
	Register(user models.RegisterUser) (models.Token, error)
	Login(user models.LoginUser) (models.Token, error)
	RefreshToken(refreshToken string) (models.Token, error)
	Delete(id string) error
	AssignRole(id, role string) error
	RemoveRole(id, role string) error
}

func NewUserService(repository repositories.UserRepository, cfg config.Config) UserService {
	log.Info().Msg("Creating new user service")

	return &userService{
		repository: repository,
		cfg:        cfg,
	}
}

type userService struct {
	repository repositories.UserRepository
	cfg        config.Config
}

func (s *userService) FindAll(query models.PaginationQuery) ([]models.User, error) {
	return s.repository.FindAll(query)
}

func (s *userService) FindById(id string) (models.User, error) {
	return s.repository.FindById(id)
}

func (s *userService) Register(user models.RegisterUser) (models.Token, error) {
	var token models.Token

	_, err := s.repository.FindByEmail(user.Email)
	if err == nil {
		return token, errors.New("user already exists")
	}

	if user.Password == "" {
		user.Password = uuid.New().String()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return token, err
	}

	newUser := models.User{
		Email:    user.Email,
		Name:     user.Name,
		Password: string(hashedPassword),
		Roles:    []string{models.UserRole},
	}

	err = s.repository.Create(&newUser)
	if err != nil {
		return token, err
	}

	accessToken, refreshToken, err := auth.GenerateTokenPair(newUser.ID, newUser.Roles, s.cfg.AccessTokenLifespan, s.cfg.AccessTokenSecret, s.cfg.RefreshTokenSecret)
	if err != nil {
		return token, err
	}

	token.Token = accessToken
	token.RefreshToken = refreshToken
	token.User = newUser

	return token, nil
}

func (s *userService) Login(user models.LoginUser) (models.Token, error) {
	var token models.Token

	existingUser, err := s.repository.FindByEmail(user.Email)
	if err != nil {
		return token, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return token, err
	}

	accessToken, refreshToken, err := auth.GenerateTokenPair(existingUser.ID, existingUser.Roles, s.cfg.AccessTokenLifespan, s.cfg.AccessTokenSecret, s.cfg.RefreshTokenSecret)
	if err != nil {
		return token, err
	}

	token.Token = accessToken
	token.RefreshToken = refreshToken
	token.User = existingUser

	return token, nil
}

func (s *userService) RefreshToken(refreshToken string) (models.Token, error) {
	var token models.Token

	t, err := auth.Validate(refreshToken, s.cfg.RefreshTokenSecret)
	if err != nil {
		return token, err
	}

	uid, err := auth.ExtractId(t)
	if err != nil {
		return token, err
	}

	user, err := s.repository.FindById(uid)
	if err != nil {
		return token, err
	}

	accessToken, refreshToken, err := auth.GenerateTokenPair(user.ID, user.Roles, s.cfg.AccessTokenLifespan, s.cfg.AccessTokenSecret, s.cfg.RefreshTokenSecret)
	if err != nil {
		return token, err
	}

	token.Token = accessToken
	token.RefreshToken = refreshToken
	token.User = user

	return token, nil
}

func (s *userService) Delete(id string) error {
	user, err := s.repository.FindById(id)
	if err != nil {
		return err
	}

	return s.repository.Delete(&user)
}

func (s *userService) AssignRole(id, role string) error {
	user, err := s.repository.FindById(id)
	if err != nil {
		return err
	}

	for _, r := range user.Roles {
		if r == role {
			return nil
		}
	}

	user.Roles = append(user.Roles, role)

	return s.repository.Update(&user)
}

func (s *userService) RemoveRole(id, role string) error {
	user, err := s.repository.FindById(id)
	if err != nil {
		return err
	}

	for i, r := range user.Roles {
		if r == role {
			user.Roles = append(user.Roles[:i], user.Roles[i+1:]...)
			break
		}
	}

	return s.repository.Update(&user)
}
