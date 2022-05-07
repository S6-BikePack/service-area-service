package repositories

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"service-area-service/config"
	"service-area-service/internal/core/domain"
	"service-area-service/pkg/logging"
	"testing"
)

type ServiceAreaRepositoryTestSuite struct {
	suite.Suite
	TestDb   *gorm.DB
	TestRepo *ServiceAreaRepository
	Cfg      *config.Config
	TestData struct {
		Area        domain.Area
		ServiceArea domain.ServiceArea
	}
}

func (suite *ServiceAreaRepositoryTestSuite) SetupSuite() {
	cfgPath := "../../test/service-area.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	logger, err := logging.NewSugaredOtelZap(cfg)
	defer logger.Close()

	if err != nil {
		panic(errors.WithStack(err))
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic(errors.WithStack(err))
	}

	repository, err := NewServiceAreaRepository(db)

	if err != nil {
		panic(errors.WithStack(err))
	}

	db.Exec("DELETE FROM public.service_areas")

	area := domain.NewArea([][]float64{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}})
	serviceArea := domain.NewServiceArea(1, "tst", "test-area", area)

	suite.Cfg = cfg
	suite.TestDb = db
	suite.TestRepo = repository
	suite.TestData = struct {
		Area        domain.Area
		ServiceArea domain.ServiceArea
	}{
		Area:        area,
		ServiceArea: serviceArea,
	}
}

func (suite *ServiceAreaRepositoryTestSuite) TestRepository_CreateServiceArea() {

	_, err := suite.TestRepo.Save(context.Background(), suite.TestData.ServiceArea)

	suite.NoError(err)

	queryResult := domain.ServiceArea{}
	suite.TestDb.Raw("SELECT * FROM public.service_areas WHERE id=?",
		suite.TestData.ServiceArea.ID).Scan(&queryResult)

	suite.EqualValues(queryResult.Name, suite.TestData.ServiceArea.Name)
	suite.EqualValues(queryResult.Identifier, suite.TestData.ServiceArea.Identifier)
	suite.EqualValues(queryResult.Area, suite.TestData.ServiceArea.Area)
}

func (suite *ServiceAreaRepositoryTestSuite) TestRepository_GetServiceArea() {
	result, err := suite.TestRepo.Get(context.Background(), 1)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.ServiceArea.Name, result.Name)
	suite.EqualValues(suite.TestData.ServiceArea.Identifier, result.Identifier)
}

func (suite *ServiceAreaRepositoryTestSuite) TestRepository_GetServiceAreas() {
	suite.TestDb.Exec("INSERT INTO public.service_areas (id, identifier, name, rider_coverage, area) VALUES (2, 'tst-2', 'test-area-2', 0, '0103000020E61000000100000005000000000000000000000000000000000000000000000000000000000000000000F03F000000000000F03F000000000000F03F000000000000F03F000000000000000000000000000000000000000000000000'::geometry(Polygon,4326))")

	result, err := suite.TestRepo.GetAll(context.Background())

	suite.NoError(err)

	suite.Equal(2, len(result))

	suite.EqualValues(suite.TestData.ServiceArea.Name, result[0].Name)
	suite.EqualValues(suite.TestData.ServiceArea.Identifier, result[0].Identifier)
	suite.EqualValues("test-area-2", result[1].Name)
	suite.EqualValues("tst-2", result[1].Identifier)
}

func (suite *ServiceAreaRepositoryTestSuite) TestRepository_UpdateServiceArea() {
	updatedArea := suite.TestData.ServiceArea
	updatedArea.Name = "updated-area"
	updatedArea.Identifier = "u-area"
	updatedArea.Area = domain.NewArea([][]float64{{1, 1}, {0, 1}, {0, 0}, {1, 0}, {1, 1}})

	_, err := suite.TestRepo.Update(context.Background(), updatedArea)

	suite.NoError(err)

	queryResult := domain.ServiceArea{}
	suite.TestDb.Raw("SELECT * FROM public.service_areas WHERE id=?",
		suite.TestData.ServiceArea.ID).Scan(&queryResult)

	suite.EqualValues(queryResult.Name, updatedArea.Name)
	suite.EqualValues(queryResult.Identifier, updatedArea.Identifier)
	suite.EqualValues(queryResult.Area, updatedArea.Area)
}

func TestIntegration_ServiceAreaRepositoryTestSuite(t *testing.T) {
	repoSuite := new(ServiceAreaRepositoryTestSuite)
	suite.Run(t, repoSuite)
}
