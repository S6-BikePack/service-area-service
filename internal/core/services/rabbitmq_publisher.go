package services

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"service-area-service/config"
	"service-area-service/internal/core/domain"
	"service-area-service/pkg/rabbitmq"
)

type rabbitmqPublisher struct {
	rabbitmq *rabbitmq.RabbitMQ
	tracer   trace.Tracer
	config   *config.Config
}

func NewRabbitMQPublisher(rabbitmq *rabbitmq.RabbitMQ, tracerProvider trace.TracerProvider, cfg *config.Config) *rabbitmqPublisher {
	return &rabbitmqPublisher{rabbitmq: rabbitmq, tracer: tracerProvider.Tracer("RabbitMQ.Publisher"), config: cfg}
}

func (rmq *rabbitmqPublisher) CreateServiceArea(ctx context.Context, serviceArea domain.ServiceArea) error {
	return rmq.publishJson(ctx, "create", serviceArea)
}

func (rmq *rabbitmqPublisher) UpdateServiceArea(ctx context.Context, serviceArea domain.ServiceArea) error {
	return rmq.publishJson(ctx, "update", serviceArea)
}

func (rmq *rabbitmqPublisher) publishJson(ctx context.Context, topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	ctx, span := rmq.tracer.Start(ctx, "publish")
	span.AddEvent(
		"Published message to rabbitmq",
		trace.WithAttributes(
			attribute.String("topic", topic),
			attribute.String("body", string(js))))
	span.End()

	err = rmq.rabbitmq.Channel.Publish(
		rmq.config.RabbitMQ.Exchange,
		fmt.Sprintf("service_area.%s", topic),
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         js,
		},
	)

	return err
}
