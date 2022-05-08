package repository_test

import (
	"database/sql"
	"database/sql/driver"
	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/internal/role/repository"
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

var dummyRole = []domain.Role{
	domain.Role{
		ID:   1,
		Name: "admin",
	},
	domain.Role{
		ID:   2,
		Name: "user",
	},
}

func TestRoleRepository_FindByName(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `roles` WHERE name = ? ORDER BY `roles`.`id` LIMIT 1").
		WithArgs(dummyRole[0].Name).
		WillReturnRows(sqlMock.NewRows([]string{"id", "name"}).
			AddRow(dummyRole[0].ID, dummyRole[0].Name))

	roleRepository := repository.NewRoleRepository(db)
	role, err := roleRepository.FindByName(dummyRole[0].Name)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, role)
}
