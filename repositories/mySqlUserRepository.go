package repositories

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alfredomagalhaes/fiber-user-api/types"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ErrUserEmailRequired error = errors.New("user email is required")
var ErrUserDisplayNameRequired error = errors.New("users display name is required")

type MySqlRepository struct {
	Db *gorm.DB
}

type MySqlRepoConfig struct {
	Host   string
	Port   string
	User   string
	Pwd    string
	DbName string
}

func NewMySqlRepository(config MySqlRepoConfig) (*MySqlRepository, error) {
	var repo MySqlRepository
	var err error

	connectionString :=
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.User, config.Pwd, config.Host, config.Port, config.DbName)

	repo.Db, err = gorm.Open(mysql.Open(connectionString))

	if err == nil {
		//If there is no error connecting to the database
		//initialize all tables
		migrateTables(&repo)
	}

	return &repo, err
}

// migrateTables executes the migration method from gorm lib
// to create schemas in the database
func migrateTables(repo *MySqlRepository) {
	repo.Db.AutoMigrate(&types.User{})
}

// Save saves a new user in the database.
// There is no need to initialize the ID field, it will be
// auto generated	 when saving in the database
func (r *MySqlRepository) Save(u *types.User) error {

	if strings.TrimSpace(u.Email) == "" {
		return ErrUserEmailRequired
	}

	if u.DisplayName == "" {
		return ErrUserDisplayNameRequired
	}

	result := r.Db.Create(u)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Find search a user with the given UUID
func (r *MySqlRepository) Find(uuid.UUID) (types.User, error) {
	return types.User{}, nil
}

// Update updates user information in the database
func (r *MySqlRepository) Update() error {
	return nil
}

// Delete deletes the user from the database
func (r *MySqlRepository) Delete() error {
	return nil
}

func (r *MySqlRepository) ListAll() ([]types.User, error) {
	var users []types.User
	result := r.Db.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
