package service

import (
	"context"
	"database/sql"
	"fmt"
	"goechotemplate/api/internal/model"
	"goechotemplate/api/internal/repo"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"

	db "goechotemplate/api/db/model"
	"goechotemplate/api/internal/config"
)

type ServiceTestSuite struct {
	suite.Suite
	service PersonService
	db      *sql.DB
}

func (suite *ServiceTestSuite) SetupSuite() {
	cnf, err := config.Load()
	if err != nil {
		suite.T().Fatal(err)
	}

	pgxCfg, err := pgx.ParseConfig(cnf.DBConnectionString)
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.db = stdlib.OpenDB(*pgxCfg)

	suite.service = NewPersonService(repo.NewPersonRepo(db.New(suite.db)))
}

func (suite *ServiceTestSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *ServiceTestSuite) TestCreatePerson() {
	ctx := context.Background()
	newPerson := model.Person{
		ExternalID: uuid.NewString(),
		Email:      fmt.Sprintf("%s@x.com", uuid.NewString()),
		Password:   nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	person, err := suite.service.Create(ctx, &newPerson)
	suite.NoError(err)
	suite.Equal(newPerson.ExternalID, person.ExternalID)
	suite.Equal(newPerson.Email, person.Email)
}

func (suite *ServiceTestSuite) TestGetPersonExists() {
	ctx := context.Background()
	person, err := suite.service.GetByExternalID(ctx, "random-id")
	suite.NoError(err)
	suite.Equal("random-id", person.ExternalID)
	suite.Equal("random@x.com", person.Email)
}

func (suite *ServiceTestSuite) TestGetPersonNotExists() {
	ctx := context.Background()
	_, err := suite.service.GetByExternalID(ctx, "123")
	suite.Error(err)
	suite.Equal("PersonService.GetByExternalID: PersonRepo.GetByExternalID: sql: no rows in result set", err.Error())
}

func TestPersonServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
