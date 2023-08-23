package repositories

import (
	"github.com/Marcel-MD/clean-api/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(query models.PaginationQuery) ([]models.User, error)
	FindById(id string) (models.User, error)
	Create(t *models.User) error
	Update(t *models.User) error
	Delete(t *models.User) error

	FindByEmail(email string) (models.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	log.Info().Msg("Creating new user repository")

	return &userRepository{
		BaseRepository: NewBaseRepository[models.User](db),
		db:             db,
	}
}

type userRepository struct {
	BaseRepository[models.User]
	db *gorm.DB
}

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error

	return user, err
}
