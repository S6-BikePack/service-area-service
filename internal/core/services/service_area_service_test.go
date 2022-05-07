package services

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"service-area-service/config"
	"service-area-service/internal/core/domain"
	"service-area-service/internal/core/interfaces"
	mock "service-area-service/internal/mock"
	"testing"
)

type ServiceAreaServiceTestSuite struct {
	suite.Suite
	Cfg              *config.Config
	MockRepository   *mock.ServiceAreaRepository
	MockRMQPublisher *mock.RabbitMQPublisher
	TestService      interfaces.ServiceAreaService
	TestData         struct {
		ServiceAreas []domain.ServiceArea
	}
}

func (suite *ServiceAreaServiceTestSuite) SetupSuite() {
	cfgPath := "../../../test/service-area.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	mockRepository := new(mock.ServiceAreaRepository)
	mockRMQPublisher := new(mock.RabbitMQPublisher)

	service := NewServiceAreaService(cfg, mockRMQPublisher, mockRepository)

	serviceAreas := []domain.ServiceArea{
		domain.NewServiceArea(1, "tst-1", "test-area-1", domain.NewArea([][]float64{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}})),
		domain.NewServiceArea(2, "tst-2", "test-area-2", domain.NewArea([][]float64{{1, 1}, {0, 1}, {0, 0}, {1, 0}, {1, 1}})),
	}

	suite.Cfg = cfg
	suite.MockRepository = mockRepository
	suite.MockRMQPublisher = mockRMQPublisher
	suite.TestService = service
	suite.TestData = struct {
		ServiceAreas []domain.ServiceArea
	}{
		ServiceAreas: serviceAreas,
	}
}

func (suite *ServiceAreaServiceTestSuite) TestService_Create() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockRepository.On("Save", testArea).Return(testArea, nil)
	suite.MockRMQPublisher.On("CreateServiceArea", testArea).Return(nil)

	result, err := suite.TestService.Create(context.Background(), testArea.ID, testArea.Identifier, testArea.Name, testArea.Area)

	suite.NoError(err)

	suite.EqualValues(testArea, result)
}

func (suite *ServiceAreaServiceTestSuite) TestService_Create_ErrorSaving() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockRepository.On("Save", testArea).Return(domain.ServiceArea{}, errors.New("could not save service-area"))
	suite.MockRMQPublisher.AssertNotCalled(suite.T(), "CreateServiceArea")

	result, err := suite.TestService.Create(context.Background(), testArea.ID, testArea.Identifier, testArea.Name, testArea.Area)

	suite.Error(err)
	suite.EqualValues(domain.ServiceArea{}, result)
}

func (suite *ServiceAreaServiceTestSuite) TestService_Get() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockRepository.On("Get", testArea.ID).Return(testArea, nil)

	result, err := suite.TestService.Get(context.Background(), testArea.ID)

	suite.NoError(err)

	suite.EqualValues(testArea, result)
}

func (suite *ServiceAreaServiceTestSuite) TestService_Get_NotFound() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockRepository.On("Get", testArea.ID).Return(domain.ServiceArea{}, errors.New("not found"))

	result, err := suite.TestService.Get(context.Background(), testArea.ID)

	suite.Error(err)
	suite.EqualValues(domain.ServiceArea{}, result)
}

func (suite *ServiceAreaServiceTestSuite) TestService_GetAll() {
	suite.MockRepository.On("GetAll").Return(suite.TestData.ServiceAreas, nil)

	result, err := suite.TestService.GetAll(context.Background())

	suite.NoError(err)

	suite.Equal(2, len(result))

	suite.EqualValues(suite.TestData.ServiceAreas[0], result[0])
	suite.EqualValues(suite.TestData.ServiceAreas[1], result[1])
}

func (suite *ServiceAreaServiceTestSuite) TestService_Update() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockRepository.On("Get", testArea.ID).Return(testArea, nil)

	updatedArea := testArea
	updatedArea.Name = "updated-area"
	updatedArea.Identifier = "u-area"
	updatedArea.Area = domain.NewArea([][]float64{{1, 1}, {0, 1}, {0, 0}, {1, 0}, {1, 1}})

	suite.MockRepository.On("Update", updatedArea).Return(updatedArea, nil)
	suite.MockRMQPublisher.On("UpdateServiceArea", updatedArea).Return(nil)

	result, err := suite.TestService.Update(context.Background(), updatedArea)

	suite.NoError(err)

	suite.EqualValues(updatedArea, result)
}

func (suite *ServiceAreaServiceTestSuite) TestService_Update_OnlyName() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockRepository.On("Get", testArea.ID).Return(testArea, nil)

	updatedArea := domain.ServiceArea{
		ID:   testArea.ID,
		Name: "updated-area",
	}

	expected := testArea
	expected.Name = updatedArea.Name

	suite.MockRepository.On("Update", expected).Return(expected, nil)
	suite.MockRMQPublisher.On("UpdateServiceArea", expected).Return(nil)

	result, err := suite.TestService.Update(context.Background(), updatedArea)

	suite.NoError(err)

	suite.EqualValues(expected.Name, result.Name)
	suite.EqualValues(expected.Identifier, result.Identifier)
	suite.EqualValues(expected.Area, result.Area)
}

func (suite *ServiceAreaServiceTestSuite) TestService_Update_OnlyIdentifier() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockRepository.On("Get", testArea.ID).Return(testArea, nil)

	updatedArea := domain.ServiceArea{
		ID:         testArea.ID,
		Identifier: "u-area",
	}

	expected := testArea
	expected.Identifier = updatedArea.Identifier

	suite.MockRepository.On("Update", expected).Return(expected, nil)
	suite.MockRMQPublisher.On("UpdateServiceArea", expected).Return(nil)

	result, err := suite.TestService.Update(context.Background(), updatedArea)

	suite.NoError(err)

	suite.EqualValues(expected.Name, result.Name)
	suite.EqualValues(expected.Identifier, result.Identifier)
	suite.EqualValues(expected.Area, result.Area)
}

func (suite *ServiceAreaServiceTestSuite) TestService_Update_OnlyArea() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockRepository.On("Get", testArea.ID).Return(testArea, nil)

	updatedArea := domain.ServiceArea{
		ID:   testArea.ID,
		Area: domain.NewArea([][]float64{{1, 1}, {0, 1}, {0, 0}, {1, 0}, {1, 1}}),
	}

	expected := testArea
	expected.Area = updatedArea.Area

	suite.MockRepository.On("Update", expected).Return(expected, nil)
	suite.MockRMQPublisher.On("UpdateServiceArea", expected).Return(nil)

	result, err := suite.TestService.Update(context.Background(), updatedArea)

	suite.NoError(err)

	suite.EqualValues(expected.Name, result.Name)
	suite.EqualValues(expected.Identifier, result.Identifier)
	suite.EqualValues(expected.Area, result.Area)
}

func (suite *ServiceAreaServiceTestSuite) TestService_Update_NotFound() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockRepository.On("Get", testArea.ID).Return(domain.ServiceArea{}, errors.New("not found"))
	suite.MockRepository.AssertNotCalled(suite.T(), "Update")
	suite.MockRMQPublisher.AssertNotCalled(suite.T(), "UpdateServiceArea")

	result, err := suite.TestService.Update(context.Background(), testArea)

	suite.Error(err)
	suite.EqualValues(domain.ServiceArea{}, result)
}

func (suite *ServiceAreaServiceTestSuite) TestService_Update_ErrorSaving() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockRepository.On("Get", testArea.ID).Return(testArea, nil)

	updatedArea := testArea
	updatedArea.Name = "updated-area"
	updatedArea.Identifier = "u-area"
	updatedArea.Area = domain.NewArea([][]float64{{1, 1}, {0, 1}, {0, 0}, {1, 0}, {1, 1}})

	suite.MockRepository.On("Update", updatedArea).Return(domain.ServiceArea{}, errors.New("could not save service-area"))
	suite.MockRMQPublisher.AssertNotCalled(suite.T(), "UpdateServiceArea")

	result, err := suite.TestService.Update(context.Background(), updatedArea)

	suite.Error(err)
	suite.EqualValues(testArea, result)
}

func (suite *ServiceAreaServiceTestSuite) SetupTest() {
	suite.MockRMQPublisher.ExpectedCalls = nil
	suite.MockRepository.ExpectedCalls = nil
}

func TestIntegration_ServiceAreaServiceTestSuite(t *testing.T) {
	repoSuite := new(ServiceAreaServiceTestSuite)
	suite.Run(t, repoSuite)
}
