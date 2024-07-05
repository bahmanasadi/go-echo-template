package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"goechotemplate/api/internal/account/model"
	"goechotemplate/api/internal/account/repository"
	"goechotemplate/api/internal/config"
	"testing"
	"time"
)

type PersonServiceTestSuite struct {
	suite.Suite
	service PersonService
	db      *sql.DB
}

func (suite *PersonServiceTestSuite) SetupSuite() {
	cnf, err := config.Load()
	if err != nil {
		suite.T().Fatal(err)
	}

	pgxCfg, err := pgx.ParseConfig(cnf.DBConnectionString)
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.db = stdlib.OpenDB(*pgxCfg)

	suite.service = NewPersonService(repository.NewPersonRepository(suite.db))
}

func (suite *PersonServiceTestSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *PersonServiceTestSuite) TestCreatePerson() {
	ctx := context.Background()
	newPerson := model.Person{
		ExternalID: uuid.NewString(),
		Email:      fmt.Sprintf("%s@x.com", uuid.NewString()),
		Password:   nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	person, err := suite.service.CreatePerson(ctx, &newPerson)
	suite.NoError(err)
	suite.Equal(newPerson.ExternalID, person.ExternalID)
	suite.Equal(newPerson.Email, person.Email)
}

func (suite *PersonServiceTestSuite) TestGetPersonExists() {
	ctx := context.Background()
	person, err := suite.service.GetPersonByExternalID(ctx, "random-id")
	suite.NoError(err)
	suite.Equal("random-id", person.ExternalID)
	suite.Equal("random@x.com", person.Email)
}

func (suite *PersonServiceTestSuite) TestGetPersonNotExists() {
	ctx := context.Background()
	_, err := suite.service.GetPersonByExternalID(ctx, "123")
	suite.Error(err)
	suite.Equal("GetPersonByExternalID: sql: no rows in result set", err.Error())
}

func TestPersonServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PersonServiceTestSuite))
}
