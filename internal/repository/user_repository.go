package repository

/*
	why :
	I am creating a repository layer here to abstracts the raw SQL/GORM calls.
	So that if in future I want to change teh ORM or DB, I only need to change this layer
*/
import (
	"github.com/MonalBarse/tradelog/internal/domain"
	"gorm.io/gorm"
)

// @Desc: interface -> methods to interact with userdata
type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// @Desc: returns -> new instance of UserRepository; (why? to access DB methods)
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	//verify the email and also ensure we dont get deleted users
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
