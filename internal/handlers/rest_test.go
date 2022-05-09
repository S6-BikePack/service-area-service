package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"service-area-service/config"
	"service-area-service/internal/core/domain"
	"service-area-service/internal/mock"
	"service-area-service/pkg/dto"
	"service-area-service/pkg/logging"
	"strings"
	"testing"
)

type RestHandlerTestSuite struct {
	suite.Suite
	MockService *mock.ServiceAreaService
	TestHandler *HTTPHandler
	TestRouter  *gin.Engine
	Cfg         *config.Config
	TestData    struct {
		ServiceAreas []domain.ServiceArea
	}
}

func (suite *RestHandlerTestSuite) SetupSuite() {
	cfgPath := "../../test/service-area.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	logger := logging.MockLogger{}

	mockService := new(mock.ServiceAreaService)

	router := gin.New()
	gin.SetMode(gin.TestMode)

	deliveryHandler := NewHTTPHandler(mockService, router, logger, cfg)
	deliveryHandler.SetupEndpoints()

	serviceAreas := []domain.ServiceArea{
		domain.NewServiceArea(1, "tst-1", "test-area-1", domain.NewArea([][]float64{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}})),
		domain.NewServiceArea(2, "tst-2", "test-area-2", domain.NewArea([][]float64{{1, 1}, {0, 1}, {0, 0}, {1, 0}, {1, 1}})),
	}

	suite.Cfg = cfg
	suite.MockService = mockService
	suite.TestRouter = router
	suite.TestHandler = deliveryHandler
	suite.TestData = struct {
		ServiceAreas []domain.ServiceArea
	}{
		ServiceAreas: serviceAreas,
	}
}

func (suite *RestHandlerTestSuite) SetupTest() {
	suite.MockService.ExpectedCalls = nil
}

func (suite *RestHandlerTestSuite) TestHandler_GetAll() {
	suite.MockService.On("GetAll").Return(suite.TestData.ServiceAreas, nil)

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/service-areas", nil)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.ServiceAreaListResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.Len(responseObject, 2)

	suite.EqualValues(suite.TestData.ServiceAreas[0].Name, responseObject[0].Name)
	suite.EqualValues(suite.TestData.ServiceAreas[0].Identifier, responseObject[0].Identifier)
	suite.EqualValues(suite.TestData.ServiceAreas[0].ID, responseObject[0].ID)
}

func (suite *RestHandlerTestSuite) TestHandler_GetAll_NoneFound() {
	suite.MockService.On("GetAll").Return([]domain.ServiceArea{}, errors.New("Not found"))

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/service-areas", nil)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusNotFound, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Get() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockService.On("Get", 1).Return(testArea, nil)

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/service-areas/%d", testArea.ID), nil)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.ServiceAreaResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.EqualValues(testArea.Name, responseObject.Name)
	suite.EqualValues(testArea.Identifier, responseObject.Identifier)
	suite.EqualValues(testArea.ID, responseObject.ID)
	suite.EqualValues(testArea.Area, responseObject.Area)
}

func (suite *RestHandlerTestSuite) TestHandler_Get_BadID() {
	suite.MockService.On("Get", "test").Return(domain.ServiceArea{}, nil)

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/service-areas/%s", "test"), nil)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Get_NotFound() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockService.On("Get", 1).Return(domain.ServiceArea{}, errors.New("Not found"))

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/service-areas/%d", testArea.ID), nil)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusNotFound, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Create() {
	testArea := suite.TestData.ServiceAreas[0]

	suite.MockService.On("Create", testArea.ID, testArea.Identifier, testArea.Name, testArea.Area).Return(suite.TestData.ServiceAreas[0], nil)

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateServiceArea{
		ID:         testArea.ID,
		Identifier: testArea.Identifier,
		Name:       testArea.Name,
		Area:       testArea.Area,
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPost, "/api/service-areas", strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.ServiceAreaResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.EqualValues(testArea.Name, responseObject.Name)
	suite.EqualValues(testArea.Identifier, responseObject.Identifier)
	suite.EqualValues(testArea.ID, responseObject.ID)
	suite.EqualValues(testArea.Area, responseObject.Area)
}

func (suite *RestHandlerTestSuite) TestHandler_Create_BadInput() {
	testArea := suite.TestData.ServiceAreas[0]

	rr := httptest.NewRecorder()

	data, err := json.Marshal(struct {
		AreaName string
	}{
		AreaName: testArea.Name,
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPost, "/api/service-areas", strings.NewReader(string(data)))

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Create_CouldNotCreate() {
	testArea := suite.TestData.ServiceAreas[1]

	suite.MockService.On("Create", testArea.ID, testArea.Identifier, testArea.Name, testArea.Area).Return(domain.ServiceArea{}, errors.New("could not create"))

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateServiceArea{
		ID:         testArea.ID,
		Identifier: testArea.Identifier,
		Name:       testArea.Name,
		Area:       testArea.Area,
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPost, "/api/service-areas", strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusInternalServerError, rr.Code)
}

func TestIntegration_RestHandlerTestSuite(t *testing.T) {
	repoSuite := new(RestHandlerTestSuite)
	suite.Run(t, repoSuite)
}
