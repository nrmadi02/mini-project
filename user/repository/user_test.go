package repository_test

import (
	"database/sql"
	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go_cicd/domain"
	"go_cicd/user/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
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

func TestUserRepository_ReadAll(t *testing.T) {
	dbMock, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mockUser := []domain.User{
		domain.User{
			ID:       1,
			Email:    "test@gmail.com",
			Password: "12345",
		},
		domain.User{
			ID:       2,
			Email:    "test2@gmail.com",
			Password: "12345",
		},
	}

	rows := sqlMock.NewRows([]string{"id", "email", "password"}).
		AddRow(mockUser[0].ID, mockUser[0].Email, mockUser[0].Password).
		AddRow(mockUser[1].ID, mockUser[1].Email, mockUser[1].Password)

	query := "SELECT"

	mock.ExpectQuery(query).WillReturnRows(rows)
	mock.ExpectCommit()

	userRepository := repository.NewUserRepository(db)
	response, err := userRepository.ReadAll()

	assert.NotEmpty(t, response)
	assert.NoError(t, err)
	assert.Len(t, *response, 2)
}

func TestUserRepository_ReadByID(t *testing.T) {
	dbMock, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	user := domain.User{
		ID:       1,
		Email:    "test@gmail.com",
		Password: "12345",
	}

	mock.ExpectQuery(
		"SELECT(.*)").
		WithArgs(1).
		WillReturnRows(
			sqlMock.NewRows([]string{"email", "password", "id"}).
				AddRow(user.Email, user.Password, user.ID))

	userRepository := repository.NewUserRepository(db)
	response, err := userRepository.ReadByID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, *response)

	_, err = userRepository.ReadByID(2)
	assert.Error(t, err)

}

func TestUserRepository_Create(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	user := domain.User{
		ID:       1,
		Email:    "test@gmail.com",
		Password: "12345",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users` (`email`,`password`,`id`) VALUES (?,?,?)").
		WithArgs(user.Email, user.Password, user.ID).
		WillReturnResult(sqlMock.NewResult(1, 1))
	mock.ExpectCommit()

	userRepository := repository.NewUserRepository(db)
	response, err := userRepository.Create(&user)

	assert.NoError(t, err)
	assert.NotNil(t, *response)

	user2 := domain.User{
		ID:       1,
		Email:    "test@gmail.com",
		Password: "12345",
	}

	_, err = userRepository.Create(&user2)

	assert.Error(t, err)
}

func TestUserRepository_CheckLogin(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	user := domain.User{
		ID:       1,
		Email:    "test@gmail.com",
		Password: "12345",
	}

	mock.ExpectQuery("SELECT * FROM `users` WHERE (email = ? AND password = ?) AND `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1").
		WithArgs(user.Email, user.Password, user.ID).
		WillReturnRows(sqlMock.
			NewRows([]string{"id", "email", "password"}).
			AddRow(user.ID, user.Email, user.Password))

	userRepository := repository.NewUserRepository(db)
	response, val, err := userRepository.CheckLogin(&user)
	assert.Equal(t, val, true)
	assert.NoError(t, err)
	assert.NotNil(t, *response)

	user2 := domain.User{
		ID:       3,
		Email:    "testttt@gmail.com",
		Password: "12345",
	}
	_, val, err = userRepository.CheckLogin(&user2)
	assert.Equal(t, val, false)
	assert.Error(t, err)

}
