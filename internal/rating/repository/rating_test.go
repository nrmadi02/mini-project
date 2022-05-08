package repository_test

import (
	"database/sql"
	"database/sql/driver"
	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/internal/rating/repository"
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

var dummyRating = []domain.RatingEnterprise{
	domain.RatingEnterprise{
		ID:           uuid.FromStringOrNil("1"),
		Rating:       3,
		EnterpriseID: uuid.FromStringOrNil("1"),
		UserID:       uuid.FromStringOrNil("1"),
	},
	domain.RatingEnterprise{
		ID:           uuid.FromStringOrNil("2"),
		Rating:       4,
		EnterpriseID: uuid.FromStringOrNil("1"),
		UserID:       uuid.FromStringOrNil("2"),
	},
}

func TestRatingRepository_FindRatingByIDUserAndEnterprise(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `rating_enterprises` WHERE enterprise_id = ? AND user_id = ?").
		WithArgs(dummyRating[0].EnterpriseID, dummyRating[0].UserID).
		WillReturnRows(sqlMock.NewRows([]string{"id", "rating", "enterprise_id", "user_id"}).
			AddRow(dummyRating[0].ID, dummyRating[0].Rating, dummyRating[0].EnterpriseID, dummyRating[0].UserID))

	ratingRepository := repository.NewRatingRepository(db)
	rating, err := ratingRepository.FindRatingByIDUserAndEnterprise(dummyRating[0].EnterpriseID.String(), dummyRating[0].UserID.String())
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, rating)
}

func TestRatingRepository_GetAllRatingByEnterpriseID(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `rating_enterprises` WHERE enterprise_id = ?").
		WithArgs(dummyRating[0].EnterpriseID).
		WillReturnRows(sqlMock.NewRows([]string{"id", "rating", "enterprise_id", "user_id"}).
			AddRow(dummyRating[0].ID, dummyRating[0].Rating, dummyRating[0].EnterpriseID, dummyRating[0].UserID).
			AddRow(dummyRating[1].ID, dummyRating[1].Rating, dummyRating[1].EnterpriseID, dummyRating[1].UserID))

	ratingRepository := repository.NewRatingRepository(db)
	ratings, err := ratingRepository.GetAllRatingByEnterpriseID(dummyRating[0].EnterpriseID.String())
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, ratings)
}

func TestRatingRepository_AddRating(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `rating_enterprises` (`id`,`rating`,`enterprise_id`,`user_id`) VALUES (?,?,?,?)").
		WithArgs(dummyRating[0].ID, int(dummyRating[0].Rating), dummyRating[0].EnterpriseID, dummyRating[0].UserID).WillReturnResult(sqlMock.NewErrorResult(nil))
	mock.ExpectCommit()

	ratingRepository := repository.NewRatingRepository(db)
	rating, err := ratingRepository.AddRating(dummyRating[0])
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, rating)
}

func TestRatingRepository_UpdateRating(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `rating_enterprises` SET `rating`=? WHERE enterprise_id = ? AND user_id = ?").
		WithArgs(int(dummyRating[0].Rating), dummyRating[0].EnterpriseID, dummyRating[0].UserID).WillReturnResult(sqlMock.NewErrorResult(nil))
	mock.ExpectCommit()

	ratingRepository := repository.NewRatingRepository(db)
	rating, err := ratingRepository.UpdateRating(dummyRating[0].EnterpriseID.String(), dummyRating[0].UserID.String(), dummyRating[0].Rating)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, rating)
}

func TestRatingRepository_DeleteRating(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `rating_enterprises` WHERE enterprise_id = ? AND user_id = ?").
		WithArgs(dummyRating[0].EnterpriseID, dummyRating[0].UserID).WillReturnResult(sqlMock.NewErrorResult(nil))
	mock.ExpectCommit()

	ratingRepository := repository.NewRatingRepository(db)
	err = ratingRepository.DeleteRating(dummyRating[0])
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
}
