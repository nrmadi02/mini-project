package repository_test

import (
	"database/sql"
	"database/sql/driver"
	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/internal/review/repository"
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

var dummyReview = []domain.Review{
	domain.Review{
		ID:           uuid.FromStringOrNil("1"),
		Review:       "Bagus Sekali",
		EnterpriseID: uuid.FromStringOrNil("1"),
		UserID:       uuid.FromStringOrNil("1"),
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	},
	domain.Review{
		ID:           uuid.FromStringOrNil("2"),
		Review:       "Bagus Sekali 2",
		EnterpriseID: uuid.FromStringOrNil("2"),
		UserID:       uuid.FromStringOrNil("2"),
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	},
}

func TestReviewRepository_FindByID(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `reviews` WHERE id = ?").
		WithArgs(dummyReview[0].ID).
		WillReturnRows(sqlMock.NewRows([]string{"id", "review", "enterprise_id", "user_id", "created_at", "updated_at"}).
			AddRow(dummyReview[0].ID, dummyReview[0].Review, dummyReview[0].EnterpriseID, dummyReview[0].UserID, dummyReview[0].CreatedAt, dummyReview[0].UpdatedAt))

	reviewRepository := repository.NewReviewRepository(db)
	review, err := reviewRepository.FindByID(dummyReview[0].ID.String())
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, review)
}

func TestReviewRepository_FindByEnterpriseID(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `reviews` WHERE enterprise_id = ?").
		WithArgs(dummyReview[0].EnterpriseID).
		WillReturnRows(sqlMock.NewRows([]string{"id", "review", "enterprise_id", "user_id", "created_at", "updated_at"}).
			AddRow(dummyReview[0].ID, dummyReview[0].Review, dummyReview[0].EnterpriseID, dummyReview[0].UserID, dummyReview[0].CreatedAt, dummyReview[0].UpdatedAt))

	reviewRepository := repository.NewReviewRepository(db)
	review, err := reviewRepository.FindByEnterpriseID(dummyReview[0].EnterpriseID.String())
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, review)
}

func TestReviewRepository_FindByUserIDAndEnterpriseID(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectQuery("SELECT * FROM `reviews` WHERE enterprise_id = ? AND user_id = ?").
		WithArgs(dummyReview[0].EnterpriseID, dummyReview[0].UserID).
		WillReturnRows(sqlMock.NewRows([]string{"id", "review", "enterprise_id", "user_id", "created_at", "updated_at"}).
			AddRow(dummyReview[0].ID, dummyReview[0].Review, dummyReview[0].EnterpriseID, dummyReview[0].UserID, dummyReview[0].CreatedAt, dummyReview[0].UpdatedAt))

	reviewRepository := repository.NewReviewRepository(db)
	review, err := reviewRepository.FindByUserIDAndEnterpriseID(dummyReview[0].EnterpriseID.String(), dummyReview[0].UserID.String())
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, review)
}

func TestReviewRepository_Add(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reviews` (`id`,`review`,`enterprise_id`,`user_id`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)").
		WithArgs(dummyReview[0].ID, dummyReview[0].Review, dummyReview[0].EnterpriseID, dummyReview[0].UserID, AnyTime{}, AnyTime{}).WillReturnResult(sqlMock.NewErrorResult(nil))
	mock.ExpectCommit()

	userRepository := repository.NewReviewRepository(db)
	review, err := userRepository.Add(dummyReview[0])
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, review)
}

func TestReviewRepository_Update(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `reviews` SET `review`=?,`updated_at`=? WHERE enterprise_id = ? AND user_id = ?").
		WithArgs(dummyReview[1].Review, AnyTime{}, dummyReview[0].EnterpriseID, dummyReview[0].UserID).
		WillReturnResult(sqlMock.NewErrorResult(nil))
	mock.ExpectCommit()

	userRepository := repository.NewReviewRepository(db)
	review, err := userRepository.Update(dummyReview[0].EnterpriseID.String(), dummyReview[0].UserID.String(), dummyReview[1].Review)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, review)
}

func TestReviewRepository_Delete(t *testing.T) {
	dbMock, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").
		WillReturnResult(sqlMock.NewErrorResult(nil))
	mock.ExpectCommit()

	userRepository := repository.NewReviewRepository(db)
	err = userRepository.Delete(dummyReview[0])
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
}
