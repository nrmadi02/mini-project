package repository_test

import (
	"database/sql"
	"database/sql/driver"
	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/internal/user/repository"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func SetupDBMock(dbMock *sql.DB) *gorm.DB {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      dbMock,
		DSN:                       "sqlmock_db_0",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{PrepareStmt: false})
	if err != nil {
		panic(err)
	}
	return gormDB
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

var dummyUser = []domain.User{
	domain.User{
		ID:               uuid.FromStringOrNil("1"),
		Fullname:         "user1",
		Email:            "satu@email.com",
		Username:         "usr1",
		Password:         "12345678",
		Roles:            nil,
		Enterprises:      nil,
		RatingEnterprise: nil,
		Reviews:          nil,
		Favorite:         domain.Favorite{},
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
	},
	domain.User{
		ID:               uuid.FromStringOrNil("2"),
		Fullname:         "user3",
		Email:            "dua@email.com",
		Username:         "usr2",
		Password:         "12345678",
		Roles:            nil,
		Enterprises:      nil,
		RatingEnterprise: nil,
		Reviews:          nil,
		Favorite:         domain.Favorite{},
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
	},
}

func TestUserRepository_FindAllUsers(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `users`").
		WillReturnRows(sqlMock.NewRows([]string{"id", "fullname", "email", "username", "password", "created_at", "updated_at"}).
			AddRow(dummyUser[0].ID, dummyUser[0].Fullname, dummyUser[0].Email, dummyUser[0].Username, dummyUser[0].Password, dummyUser[0].CreatedAt, dummyUser[0].UpdatedAt).
			AddRow(dummyUser[1].ID, dummyUser[1].Fullname, dummyUser[1].Email, dummyUser[1].Username, dummyUser[1].Password, dummyUser[1].CreatedAt, dummyUser[1].UpdatedAt))

	userRepository := repository.NewUserRepository(db)
	users, err := userRepository.FindAllUsers()
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, users)
}

func TestUserRepository_FindUserById(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `users` WHERE id = ? ORDER BY `users`.`id` LIMIT 1").
		WithArgs(dummyUser[0].ID).
		WillReturnRows(sqlMock.NewRows([]string{"id", "fullname", "email", "username", "password", "created_at", "updated_at"}).
			AddRow(dummyUser[0].ID, dummyUser[0].Fullname, dummyUser[0].Email, dummyUser[0].Username, dummyUser[0].Password, dummyUser[0].CreatedAt, dummyUser[0].UpdatedAt))

	userRepository := repository.NewUserRepository(db)
	user, err := userRepository.FindUserById(dummyUser[0].ID.String())
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindUserByEmail(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `users` WHERE email = ? ORDER BY `users`.`id` LIMIT 1").
		WithArgs(dummyUser[0].Email).
		WillReturnRows(sqlMock.NewRows([]string{"id", "fullname", "email", "username", "password", "created_at", "updated_at"}).
			AddRow(dummyUser[0].ID, dummyUser[0].Fullname, dummyUser[0].Email, dummyUser[0].Username, dummyUser[0].Password, dummyUser[0].CreatedAt, dummyUser[0].UpdatedAt))

	userRepository := repository.NewUserRepository(db)
	user, err := userRepository.FindUserByEmail(dummyUser[0].Email)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_Save(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users` (`id`,`fullname`,`email`,`username`,`password`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)").
		WithArgs(dummyUser[0].ID, dummyUser[0].Fullname, dummyUser[0].Email, dummyUser[0].Username, dummyUser[0].Password, AnyTime{}, AnyTime{}).WillReturnResult(sqlMock.NewErrorResult(nil))
	mock.ExpectCommit()

	userRepository := repository.NewUserRepository(db)
	user, err := userRepository.Save(dummyUser[0])
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, user)
}
