package repository_test

import (
	"database/sql"
	"database/sql/driver"
	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/internal/enterprise/repository"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

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

var created_at, _ = time.Parse("2022-05-07 18:11:36.681 +0800 WITA", "2022-05-07 18:11:36.681 +0800 WITA")
var updated_at, _ = time.Parse("2022-05-07 18:11:36.681 +0800 WITA", "2022-05-07 18:11:36.681 +0800 WITA")

var dummyEnterprise = []domain.Enterprise{
	domain.Enterprise{
		ID:               uuid.FromStringOrNil("2"),
		UserID:           uuid.FromStringOrNil("3"),
		Name:             "enterprise satu",
		NumberPhone:      "0012798232",
		Address:          "bjb",
		Postcode:         707722,
		Latitude:         "4235,23",
		Longitude:        "4225,233231",
		Description:      "testing1",
		Status:           0,
		Tags:             nil,
		RatingEnterprise: nil,
		Reviews:          nil,
		CreatedAt:        created_at,
		UpdatedAt:        updated_at,
	},
	domain.Enterprise{
		ID:               uuid.FromStringOrNil("3"),
		UserID:           uuid.FromStringOrNil("4"),
		Name:             "enterprise dua",
		NumberPhone:      "0012798232",
		Address:          "bjb",
		Postcode:         707722,
		Latitude:         "4235,23",
		Longitude:        "4225,233231",
		Description:      "testing1",
		Status:           0,
		Tags:             nil,
		RatingEnterprise: nil,
		Reviews:          nil,
		CreatedAt:        created_at,
		UpdatedAt:        updated_at,
	},
}

func TestEnterpriseRepository_FindAll(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `enterprises`").
		WillReturnRows(sqlMock.
			NewRows([]string{"id", "name", "user_id", "number_phone",
				"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
			AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
				dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
				dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description).
			AddRow(dummyEnterprise[1].ID, dummyEnterprise[1].Name, dummyEnterprise[1].UserID, dummyEnterprise[1].NumberPhone,
				dummyEnterprise[1].Address, dummyEnterprise[1].Status, dummyEnterprise[1].Postcode, dummyEnterprise[1].Longitude,
				dummyEnterprise[1].Latitude, dummyEnterprise[1].CreatedAt, dummyEnterprise[1].UpdatedAt, dummyEnterprise[1].Description))

	mock.ExpectQuery("SELECT * FROM `enterprises` LIMIT 2").
		WillReturnRows(sqlMock.
			NewRows([]string{"id", "name", "user_id", "number_phone",
				"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
			AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
				dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
				dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description).
			AddRow(dummyEnterprise[1].ID, dummyEnterprise[1].Name, dummyEnterprise[1].UserID, dummyEnterprise[1].NumberPhone,
				dummyEnterprise[1].Address, dummyEnterprise[1].Status, dummyEnterprise[1].Postcode, dummyEnterprise[1].Longitude,
				dummyEnterprise[1].Latitude, dummyEnterprise[1].CreatedAt, dummyEnterprise[1].UpdatedAt, dummyEnterprise[1].Description))

	enterpriseRepository := repository.NewEnterpriseRepository(db)
	enterprises, totalData, err := enterpriseRepository.FindAll("", 1, 2)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, totalData)
	assert.NotNil(t, enterprises)
}

func TestEnterpriseRepository_FindAllSearch(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)
	var page int
	page = 1
	if page == 1 {
		mock.ExpectQuery("SELECT * FROM `enterprises` WHERE name LIKE ?").
			WithArgs("%e%").
			WillReturnRows(sqlMock.
				NewRows([]string{"id", "name", "user_id", "number_phone",
					"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
				AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
					dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
					dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description).
				AddRow(dummyEnterprise[1].ID, dummyEnterprise[1].Name, dummyEnterprise[1].UserID, dummyEnterprise[1].NumberPhone,
					dummyEnterprise[1].Address, dummyEnterprise[1].Status, dummyEnterprise[1].Postcode, dummyEnterprise[1].Longitude,
					dummyEnterprise[1].Latitude, dummyEnterprise[1].CreatedAt, dummyEnterprise[1].UpdatedAt, dummyEnterprise[1].Description))

		mock.ExpectQuery("SELECT * FROM `enterprises` WHERE name LIKE ? LIMIT 1").
			WithArgs("%e%").
			WillReturnRows(sqlMock.
				NewRows([]string{"id", "name", "user_id", "number_phone",
					"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
				AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
					dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
					dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description).
				AddRow(dummyEnterprise[1].ID, dummyEnterprise[1].Name, dummyEnterprise[1].UserID, dummyEnterprise[1].NumberPhone,
					dummyEnterprise[1].Address, dummyEnterprise[1].Status, dummyEnterprise[1].Postcode, dummyEnterprise[1].Longitude,
					dummyEnterprise[1].Latitude, dummyEnterprise[1].CreatedAt, dummyEnterprise[1].UpdatedAt, dummyEnterprise[1].Description))

		enterpriseRepository := repository.NewEnterpriseRepository(db)
		enterprises, totalData, err := enterpriseRepository.FindAll("e", 1, 1)
		if err != nil {
			assert.Error(t, err)
		}
		assert.NoError(t, err)
		assert.NotNil(t, totalData)
		assert.NotNil(t, enterprises)
	}

	page = 2
	if page == 2 {
		mock.ExpectQuery("SELECT * FROM `enterprises` WHERE name LIKE ?").
			WithArgs("%e%").
			WillReturnRows(sqlMock.
				NewRows([]string{"id", "name", "user_id", "number_phone",
					"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
				AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
					dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
					dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description).
				AddRow(dummyEnterprise[1].ID, dummyEnterprise[1].Name, dummyEnterprise[1].UserID, dummyEnterprise[1].NumberPhone,
					dummyEnterprise[1].Address, dummyEnterprise[1].Status, dummyEnterprise[1].Postcode, dummyEnterprise[1].Longitude,
					dummyEnterprise[1].Latitude, dummyEnterprise[1].CreatedAt, dummyEnterprise[1].UpdatedAt, dummyEnterprise[1].Description))

		mock.ExpectQuery("SELECT * FROM `enterprises` WHERE name LIKE ? LIMIT 1 OFFSET 1").
			WithArgs("%e%").
			WillReturnRows(sqlMock.
				NewRows([]string{"id", "name", "user_id", "number_phone",
					"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
				AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
					dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
					dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description).
				AddRow(dummyEnterprise[1].ID, dummyEnterprise[1].Name, dummyEnterprise[1].UserID, dummyEnterprise[1].NumberPhone,
					dummyEnterprise[1].Address, dummyEnterprise[1].Status, dummyEnterprise[1].Postcode, dummyEnterprise[1].Longitude,
					dummyEnterprise[1].Latitude, dummyEnterprise[1].CreatedAt, dummyEnterprise[1].UpdatedAt, dummyEnterprise[1].Description))

		enterpriseRepository := repository.NewEnterpriseRepository(db)
		enterprises, totalData, err := enterpriseRepository.FindAll("e", 2, 1)
		if err != nil {
			assert.Error(t, err)
		}
		assert.NoError(t, err)
		assert.NotNil(t, totalData)
		assert.NotNil(t, enterprises)
	}
}

func TestEnterpriseRepository_FindAllPage(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)
	var page int
	page = 1
	if page == 1 {
		mock.ExpectQuery("SELECT * FROM `enterprises` WHERE name LIKE ?").
			WithArgs("%e%").
			WillReturnRows(sqlMock.
				NewRows([]string{"id", "name", "user_id", "number_phone",
					"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
				AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
					dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
					dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description).
				AddRow(dummyEnterprise[1].ID, dummyEnterprise[1].Name, dummyEnterprise[1].UserID, dummyEnterprise[1].NumberPhone,
					dummyEnterprise[1].Address, dummyEnterprise[1].Status, dummyEnterprise[1].Postcode, dummyEnterprise[1].Longitude,
					dummyEnterprise[1].Latitude, dummyEnterprise[1].CreatedAt, dummyEnterprise[1].UpdatedAt, dummyEnterprise[1].Description))

		mock.ExpectQuery("SELECT * FROM `enterprises` WHERE name LIKE ? LIMIT 1").
			WithArgs("%e%").
			WillReturnRows(sqlMock.
				NewRows([]string{"id", "name", "user_id", "number_phone",
					"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
				AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
					dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
					dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description).
				AddRow(dummyEnterprise[1].ID, dummyEnterprise[1].Name, dummyEnterprise[1].UserID, dummyEnterprise[1].NumberPhone,
					dummyEnterprise[1].Address, dummyEnterprise[1].Status, dummyEnterprise[1].Postcode, dummyEnterprise[1].Longitude,
					dummyEnterprise[1].Latitude, dummyEnterprise[1].CreatedAt, dummyEnterprise[1].UpdatedAt, dummyEnterprise[1].Description))

		enterpriseRepository := repository.NewEnterpriseRepository(db)
		enterprises, totalData, err := enterpriseRepository.FindAll("e", 0, 1)
		if err != nil {
			assert.Error(t, err)
		}
		assert.NoError(t, err)
		assert.NotNil(t, totalData)
		assert.NotNil(t, enterprises)
	}

	page = 2
	if page == 2 {
		mock.ExpectQuery("SELECT * FROM `enterprises` WHERE name LIKE ?").
			WithArgs("%e%").
			WillReturnRows(sqlMock.
				NewRows([]string{"id", "name", "user_id", "number_phone",
					"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
				AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
					dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
					dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description).
				AddRow(dummyEnterprise[1].ID, dummyEnterprise[1].Name, dummyEnterprise[1].UserID, dummyEnterprise[1].NumberPhone,
					dummyEnterprise[1].Address, dummyEnterprise[1].Status, dummyEnterprise[1].Postcode, dummyEnterprise[1].Longitude,
					dummyEnterprise[1].Latitude, dummyEnterprise[1].CreatedAt, dummyEnterprise[1].UpdatedAt, dummyEnterprise[1].Description))

		mock.ExpectQuery("SELECT * FROM `enterprises` WHERE name LIKE ? LIMIT 1 OFFSET 1").
			WithArgs("%e%").
			WillReturnRows(sqlMock.
				NewRows([]string{"id", "name", "user_id", "number_phone",
					"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
				AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
					dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
					dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description).
				AddRow(dummyEnterprise[1].ID, dummyEnterprise[1].Name, dummyEnterprise[1].UserID, dummyEnterprise[1].NumberPhone,
					dummyEnterprise[1].Address, dummyEnterprise[1].Status, dummyEnterprise[1].Postcode, dummyEnterprise[1].Longitude,
					dummyEnterprise[1].Latitude, dummyEnterprise[1].CreatedAt, dummyEnterprise[1].UpdatedAt, dummyEnterprise[1].Description))

		enterpriseRepository := repository.NewEnterpriseRepository(db)
		enterprises, totalData, err := enterpriseRepository.FindAll("e", 2, 1)
		if err != nil {
			assert.Error(t, err)
		}
		assert.NoError(t, err)
		assert.NotNil(t, totalData)
		assert.NotNil(t, enterprises)
	}
}

func TestEnterpriseRepository_FindByID(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `enterprises` WHERE id = ?").
		WithArgs(dummyEnterprise[0].ID).
		WillReturnRows(sqlMock.
			NewRows([]string{"id", "name", "user_id", "number_phone",
				"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
			AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
				dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
				dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description))

	enterpriseRepository := repository.NewEnterpriseRepository(db)
	enterprise, err := enterpriseRepository.FindByID(dummyEnterprise[0].ID.String())
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, enterprise)
}

func TestEnterpriseRepository_FindByUserID(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `enterprises` WHERE user_id = ?").
		WithArgs(dummyEnterprise[0].UserID).
		WillReturnRows(sqlMock.
			NewRows([]string{"id", "name", "user_id", "number_phone",
				"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
			AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
				dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
				dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description))

	enterpriseRepository := repository.NewEnterpriseRepository(db)
	enterprise, err := enterpriseRepository.FindByUserID(dummyEnterprise[0].UserID.String())
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, enterprise)
}

func TestEnterpriseRepository_FindByIDs(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `enterprises` WHERE id IN (?,?)").
		WithArgs(dummyEnterprise[0].ID, dummyEnterprise[1].ID).
		WillReturnRows(sqlMock.
			NewRows([]string{"id", "name", "user_id", "number_phone",
				"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
			AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
				dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
				dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description).
			AddRow(dummyEnterprise[1].ID, dummyEnterprise[1].Name, dummyEnterprise[1].UserID, dummyEnterprise[1].NumberPhone,
				dummyEnterprise[1].Address, dummyEnterprise[1].Status, dummyEnterprise[1].Postcode, dummyEnterprise[1].Longitude,
				dummyEnterprise[1].Latitude, dummyEnterprise[1].CreatedAt, dummyEnterprise[1].UpdatedAt, dummyEnterprise[1].Description))

	enterpriseRepository := repository.NewEnterpriseRepository(db)
	enterprises, err := enterpriseRepository.FindByIDs([]string{dummyEnterprise[0].ID.String(), dummyEnterprise[1].ID.String()})
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, enterprises)
}

func TestEnterpriseRepository_FindByStatusDraft(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `enterprises` WHERE status = ?").
		WithArgs(0).
		WillReturnRows(sqlMock.
			NewRows([]string{"id", "name", "user_id", "number_phone",
				"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
			AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
				dummyEnterprise[0].Address, dummyEnterprise[0].Status, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
				dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description))

	enterpriseRepository := repository.NewEnterpriseRepository(db)
	enterprise, err := enterpriseRepository.FindByStatusDraft()
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, enterprise)
}

func TestEnterpriseRepository_FindByStatusPublish(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `enterprises` WHERE status = ?").
		WithArgs(1).
		WillReturnRows(sqlMock.
			NewRows([]string{"id", "name", "user_id", "number_phone",
				"address", "status", "postcode", "longitude", "latitude", "created_at", "updated_at", "description"}).
			AddRow(dummyEnterprise[0].ID, dummyEnterprise[0].Name, dummyEnterprise[0].UserID, dummyEnterprise[0].NumberPhone,
				dummyEnterprise[0].Address, 1, dummyEnterprise[0].Postcode, dummyEnterprise[0].Longitude,
				dummyEnterprise[0].Latitude, dummyEnterprise[0].CreatedAt, dummyEnterprise[0].UpdatedAt, dummyEnterprise[0].Description))

	enterpriseRepository := repository.NewEnterpriseRepository(db)
	enterprise, err := enterpriseRepository.FindByStatusPublish()
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, enterprise)
}

func TestEnterpriseRepository_Save(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `enterprises` (`id`,`user_id`,`name`,`number_phone`,`address`,`postcode`,`latitude`,`longitude`,`description`,`status`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)").
		WithArgs(dummyEnterprise[0].ID, dummyEnterprise[0].UserID, dummyEnterprise[0].Name, dummyEnterprise[0].NumberPhone,
			dummyEnterprise[0].Address, int(dummyEnterprise[0].Postcode),
			dummyEnterprise[0].Latitude, dummyEnterprise[0].Longitude, dummyEnterprise[0].Description, int(dummyEnterprise[0].Status), AnyTime{},
			AnyTime{}).WillReturnResult(sqlMock.NewErrorResult(nil))
	mock.ExpectCommit()
	mock.ExpectClose()
	enterpriseRepository := repository.NewEnterpriseRepository(db)
	enterprise, err := enterpriseRepository.Save(dummyEnterprise[0])
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, enterprise)
}

func TestEnterpriseRepository_UpdateStatusByID(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `enterprises` SET `status`=?,`updated_at`=? WHERE id = ?").
		WithArgs(1, AnyTime{}, dummyEnterprise[0].ID).WillReturnResult(sqlMock.NewErrorResult(nil))
	mock.ExpectCommit()
	mock.ExpectClose()
	enterpriseRepository := repository.NewEnterpriseRepository(db)
	enterprise, err := enterpriseRepository.UpdateStatusByID(dummyEnterprise[0].ID.String(), 1)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, enterprise)
}

func TestEnterpriseRepository_Delete(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `enterprise_tags` WHERE `enterprise_tags`.`enterprise_id` IN (NULL)").WillReturnResult(sqlMock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM `enterprises` WHERE id = ?").WithArgs(dummyEnterprise[0].ID).WillReturnResult(sqlMock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectClose()

	enterpriseRepository := repository.NewEnterpriseRepository(db)
	err = enterpriseRepository.Delete(dummyEnterprise[0])
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
}

func TestEnterpriseRepository_Update(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `enterprises` SET `name`=?,`number_phone`=?,`address`=?,`postcode`=?,`latitude`=?,`longitude`=?,`description`=?,`updated_at`=? WHERE id = ?").
		WithArgs(dummyEnterprise[0].Name, dummyEnterprise[0].NumberPhone, dummyEnterprise[0].Address,
			int64(dummyEnterprise[0].Postcode), dummyEnterprise[0].Latitude, dummyEnterprise[0].Longitude, dummyEnterprise[0].Description,
			AnyTime{}, dummyEnterprise[0].ID).WillReturnResult(sqlMock.NewResult(1, 1))
	mock.ExpectCommit()

	enterpriseRepository := repository.NewEnterpriseRepository(db)
	_, err = enterpriseRepository.Update(dummyEnterprise[0])
	if err != nil {
		assert.Error(t, err)
	}
	//assert.NoError(t, err)
	//assert.NotNil(t, enterprise)
}
