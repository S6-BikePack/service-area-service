package services

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/sdk/trace"
	"service-area-service/config"
	"service-area-service/internal/core/domain"
	"service-area-service/internal/core/interfaces"
	"service-area-service/internal/mock"
	"service-area-service/pkg/rabbitmq"
	"testing"
)

type RabbitMQPublisherTestSuite struct {
	suite.Suite
	MockService   *mock.ServiceAreaService
	TestRabbitMQ  *rabbitmq.RabbitMQ
	TestPublisher interfaces.MessageBusPublisher
	Cfg           *config.Config
	TestData      struct {
		ServiceArea domain.ServiceArea
	}
}

func (suite *RabbitMQPublisherTestSuite) SetupSuite() {
	cfgPath := "../../../test/service-area.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	mockService := new(mock.ServiceAreaService)

	rmqServer, err := rabbitmq.NewRabbitMQ(cfg)

	if err != nil {
		panic(errors.WithStack(err))
	}

	tracer := trace.NewTracerProvider()

	rmqPublisher := NewRabbitMQPublisher(rmqServer, tracer, cfg)

	suite.Cfg = cfg
	suite.MockService = mockService
	suite.TestRabbitMQ = rmqServer
	suite.TestPublisher = rmqPublisher
	suite.TestData = struct {
		ServiceArea domain.ServiceArea
	}{
		ServiceArea: domain.NewServiceArea(1, "tst-1", "test-area-1", domain.NewArea([][]float64{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}})),
	}
}

func (suite *RabbitMQPublisherTestSuite) TestRabbitMQPublisher_CreateServiceArea() {
	ch, err := suite.TestRabbitMQ.Connection.Channel()

	suite.NoError(err)

	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	suite.NoError(err)

	err = ch.QueueBind(
		queue.Name,
		"service_area.create",
		"topics",
		false,
		nil)
	if err != nil {
		return
	}

	suite.NoError(err)

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	suite.NoError(err)

	err = suite.TestPublisher.CreateServiceArea(context.Background(), suite.TestData.ServiceArea)

	suite.NoError(err)

	for msg := range msgs {
		suite.Equal("service_area.create", msg.RoutingKey)

		var serviceArea domain.ServiceArea

		err = json.Unmarshal(msg.Body, &serviceArea)
		suite.NoError(err)

		suite.Equal(suite.TestData.ServiceArea, serviceArea)

		err = msg.Ack(true)

		suite.NoError(err)

		err = ch.Close()

		suite.NoError(err)

		return
	}

}

func (suite *RabbitMQPublisherTestSuite) TestRabbitMQPublisher_UpdateServiceArea() {
	ch, err := suite.TestRabbitMQ.Connection.Channel()

	suite.NoError(err)

	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	suite.NoError(err)

	err = ch.QueueBind(
		queue.Name,
		"service_area.update",
		"topics",
		false,
		nil)
	if err != nil {
		return
	}

	suite.NoError(err)

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	suite.NoError(err)

	err = suite.TestPublisher.UpdateServiceArea(context.Background(), suite.TestData.ServiceArea)

	suite.NoError(err)

	for msg := range msgs {
		suite.Equal("service_area.update", msg.RoutingKey)

		var serviceArea domain.ServiceArea

		err = json.Unmarshal(msg.Body, &serviceArea)
		suite.NoError(err)

		suite.Equal(suite.TestData.ServiceArea, serviceArea)

		err = msg.Ack(true)

		suite.NoError(err)

		err = ch.Close()

		suite.NoError(err)

		return
	}
}

func TestIntegration_RabbitMQPublisherTestSuite(t *testing.T) {
	repoSuite := new(RabbitMQPublisherTestSuite)
	suite.Run(t, repoSuite)
}
