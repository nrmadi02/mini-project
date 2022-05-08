package repository_test

import (
	"database/sql"
	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/internal/tag/repository"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
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

func TestTagRepository_FindByName(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)
	tag := domain.Tag{
		ID:   uuid.FromStringOrNil("2"),
		Name: "Tag Satu",
	}
	mock.ExpectQuery("SELECT * FROM `tags` WHERE name = ? ORDER BY `tags`.`id` LIMIT 1").
		WithArgs(tag.Name).
		WillReturnRows(sqlMock.
			NewRows([]string{"id", "name"}).
			AddRow(tag.ID, tag.Name))
	tagRepository := repository.NewTagRepository(db)
	resTag, err := tagRepository.FindByName("Tag Satu")
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, resTag)
}

func TestTagRepository_FindByID(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)
	tag := domain.Tag{
		ID:   uuid.FromStringOrNil("2"),
		Name: "Tag Satu",
	}
	mock.ExpectQuery("SELECT * FROM `tags` WHERE id = ? ORDER BY `tags`.`id` LIMIT 1").
		WithArgs(tag.ID).
		WillReturnRows(sqlMock.
			NewRows([]string{"id", "name"}).
			AddRow(tag.ID, tag.Name))
	tagRepository := repository.NewTagRepository(db)
	resTag, err := tagRepository.FindByID(tag.ID.String())
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, resTag)
}

func TestTagRepository_FindAllTags(t *testing.T) {
	dbMock, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)
	tag := []domain.Tag{
		domain.Tag{
			ID:   uuid.FromStringOrNil("2"),
			Name: "Tag Satu",
		},
		domain.Tag{
			ID:   uuid.FromStringOrNil("3"),
			Name: "Tag Dua",
		},
	}
	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlMock.
			NewRows([]string{"id", "name"}).
			AddRow(tag[0].ID, tag[0].Name).
			AddRow(tag[1].ID, tag[1].Name))

	tagRepository := repository.NewTagRepository(db)
	resTags, err := tagRepository.FindAllTags()
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, resTags)
}

func TestTagRepository_FindByIDs(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)
	tag := []domain.Tag{
		domain.Tag{
			ID:   uuid.FromStringOrNil("2"),
			Name: "Tag Satu",
		},
		domain.Tag{
			ID:   uuid.FromStringOrNil("3"),
			Name: "Tag Dua",
		},
	}
	mock.ExpectQuery("SELECT * FROM `tags` WHERE id IN (?,?)").
		WithArgs(tag[0].ID, tag[1].ID).
		WillReturnRows(sqlMock.
			NewRows([]string{"id", "name"}).
			AddRow(tag[0].ID, tag[0].Name).
			AddRow(tag[1].ID, tag[1].Name))

	req := make([]string, 2)
	req[0] = tag[0].ID.String()
	req[1] = tag[1].ID.String()
	tagRepository := repository.NewTagRepository(db)
	resTags, err := tagRepository.FindByIDs(req)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, resTags)
}

func TestTagRepository_Delete(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)
	tag := []domain.Tag{
		domain.Tag{
			ID:   uuid.FromStringOrNil("2"),
			Name: "Tag Satu",
		},
		domain.Tag{
			ID:   uuid.FromStringOrNil("3"),
			Name: "Tag Dua",
		},
	}
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `tags` WHERE id = ?").
		WithArgs(tag[0].ID.String()).WillReturnResult(sqlMock.NewResult(1, 1))
	mock.ExpectCommit()

	tagRepository := repository.NewTagRepository(db)
	err = tagRepository.Delete(tag[0], tag[0].ID.String())
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
}

func TestTagRepository_Save(t *testing.T) {
	dbMock, mock, err := sqlMock.New(sqlMock.QueryMatcherOption(sqlMock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := SetupDBMock(dbMock)
	tag := []domain.Tag{
		domain.Tag{
			ID:   uuid.FromStringOrNil("2"),
			Name: "Tag Satu",
		},
		domain.Tag{
			ID:   uuid.FromStringOrNil("3"),
			Name: "Tag Dua",
		},
	}
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `tags` (`id`,`name`) VALUES (?,?)").
		WithArgs(tag[0].ID.String(), tag[0].Name).WillReturnResult(sqlMock.NewResult(1, 1))
	mock.ExpectCommit()

	tagRepository := repository.NewTagRepository(db)
	res, err := tagRepository.Save(tag[0])
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
