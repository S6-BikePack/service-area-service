package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"service-area-service/config"
	"service-area-service/internal/core/domain"
	"service-area-service/pkg/azure"
)

type azurePublisher struct {
	serviceBus *azure.ServiceBus
	sender     *azservicebus.Sender
	config     *config.Config
}

func NewAzurePublisher(serviceBus *azure.ServiceBus, cfg *config.Config) *azurePublisher {
	return &azurePublisher{serviceBus: serviceBus, config: cfg}
}

func (az *azurePublisher) CreateServiceArea(ctx context.Context, serviceArea domain.ServiceArea) error {
	return az.publishJson(ctx, "create", serviceArea)
}

func (az *azurePublisher) UpdateServiceArea(ctx context.Context, serviceArea domain.ServiceArea) error {
	return az.publishJson(ctx, "update", serviceArea)
}

func (az *azurePublisher) publishJson(ctx context.Context, topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	topic = fmt.Sprintf("service_area.%s", topic)

	sender, err := az.serviceBus.Client.NewSender(topic, nil)

	defer func(sender *azservicebus.Sender, ctx context.Context) {
		_ = sender.Close(ctx)
	}(sender, ctx)

	if err != nil {
		return err
	}

	err = sender.SendMessage(ctx, &azservicebus.Message{
		Body:    js,
		Subject: &topic,
	}, nil)

	if err != nil {
		return err
	}

	return err
}
