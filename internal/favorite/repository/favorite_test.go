package repository_test

import (
	"database/sql"
	"database/sql/driver"
	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/internal/favorite/repository"
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

var dummyEnterprise = []domain.Enterprise{
	domain.Enterprise{
		ID:               uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89a"),
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
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
	},
}

var dummyFavorite = []domain.Favorite{
	domain.Favorite{
		ID:          uuid.FromStringOrNil("1"),
		UserID:      uuid.FromStringOrNil("1"),
		Enterprises: nil,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}, domain.Favorite{
		ID:          uuid.FromStringOrNil("2"),
		UserID:      uuid.FromStringOrNil("1"),
		Enterprises: nil,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	},
}

var ID = uuid.NewV4()
var dummyFavorite2 = []domain.Favorite{
	domain.Favorite{
		ID:          ID,
		UserID:      uuid.FromStringOrNil("1"),
		Enterprises: nil,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}, domain.Favorite{
		ID:          ID,
		UserID:      uuid.FromStringOrNil("1"),
		Enterprises: nil,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	},
}

func TestFavoriteRepository_FindAll(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `favorites`").
		WillReturnRows(sqlMock.NewRows([]string{"id", "user_id", "created_at", "updated_at"}).
			AddRow(dummyFavorite[0].ID, dummyFavorite[0].UserID, dummyFavorite[0].CreatedAt, dummyFavorite[0].UpdatedAt).
			AddRow(dummyFavorite[1].ID, dummyFavorite[1].UserID, dummyFavorite[1].CreatedAt, dummyFavorite[1].UpdatedAt))

	favoriteRepository := repository.NewFavoriteRepository(db)
	favorites, err := favoriteRepository.FindAll()
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, favorites)
}

func TestFavoriteRepository_FindByID(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `favorites` WHERE id = ?").
		WithArgs(dummyFavorite[0].ID).
		WillReturnRows(sqlMock.NewRows([]string{"id", "user_id", "created_at", "updated_at"}).
			AddRow(dummyFavorite[0].ID, dummyFavorite[0].UserID, dummyFavorite[0].CreatedAt, dummyFavorite[0].UpdatedAt))

	favoriteRepository := repository.NewFavoriteRepository(db)
	favorite, err := favoriteRepository.FindByID(dummyFavorite[0].ID.String())
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, favorite)
}

func TestFavoriteRepository_FindByUserID(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `favorites` WHERE user_id = ?").
		WithArgs(dummyFavorite[0].UserID).
		WillReturnRows(sqlMock.NewRows([]string{"id", "user_id", "created_at", "updated_at"}).
			AddRow(dummyFavorite[0].ID, dummyFavorite[0].UserID, dummyFavorite[0].CreatedAt, dummyFavorite[0].UpdatedAt))

	favoriteRepository := repository.NewFavoriteRepository(db)
	favorite, err := favoriteRepository.FindByUserID(dummyFavorite[0].UserID.String())
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, favorite)
}

func TestFavoriteRepository_Add(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `favorites` (`id`,`user_id`,`created_at`,`updated_at`) VALUES (?,?,?,?)").
		WithArgs(dummyFavorite[0].ID, dummyFavorite[0].UserID, AnyTime{}, AnyTime{}).WillReturnResult(sqlMock.NewErrorResult(nil))
	mock.ExpectCommit()
	favoriteRepository := repository.NewFavoriteRepository(db)
	favorite, err := favoriteRepository.Add(dummyFavorite[0])
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, favorite)
}

func TestFavoriteRepository_Update(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)
	t.Run("remove favorite", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM `enterprise_favorites` WHERE `enterprise_favorites`.`favorite_id` = ? AND `enterprise_favorites`.`enterprise_id` = ?").
			WithArgs(dummyFavorite2[0].ID, dummyEnterprise[0].ID).WillReturnResult(sqlMock.NewResult(1, 1))
		mock.ExpectCommit()
		favoriteRepository := repository.NewFavoriteRepository(db)
		favorite, err := favoriteRepository.Update(dummyFavorite2[0], dummyEnterprise, "remove")
		if err != nil {
			assert.Error(t, err)
		}
		assert.NoError(t, err)
		assert.NotNil(t, favorite)
	})
	t.Run("add favorite", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `favorites` SET `updated_at`=? WHERE `id` = ?").
			WithArgs(AnyTime{}, dummyFavorite2[0].ID).WillReturnResult(sqlMock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO `enterprises` (`id`,`user_id`,`name`,`number_phone`,`address`,`postcode`,`latitude`,`longitude`,`description`,`status`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `id`=`id`").
			WithArgs(dummyEnterprise[0].ID, dummyEnterprise[0].UserID, dummyEnterprise[0].Name, dummyEnterprise[0].NumberPhone,
				dummyEnterprise[0].Address, int(dummyEnterprise[0].Postcode),
				dummyEnterprise[0].Latitude, dummyEnterprise[0].Longitude, dummyEnterprise[0].Description, int(dummyEnterprise[0].Status), AnyTime{}, AnyTime{}).
			WillReturnResult(sqlMock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO `enterprise_favorites` (`favorite_id`,`enterprise_id`) VALUES (?,?) ON DUPLICATE KEY UPDATE `favorite_id`=`favorite_id`").
			WithArgs(dummyFavorite2[0].ID, dummyEnterprise[0].ID).
			WillReturnResult(sqlMock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectClose()
		favoriteRepository := repository.NewFavoriteRepository(db)
		favorite, err := favoriteRepository.Update(dummyFavorite2[0], dummyEnterprise, "add")
		if err != nil {
			assert.Error(t, err)
		}
		assert.NoError(t, err)
		assert.NotNil(t, favorite)
	})
}
