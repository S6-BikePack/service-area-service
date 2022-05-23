package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"service-area-service/config"
)

type ServiceBus struct {
	Client *azservicebus.Client
}

func NewAzureServiceBus(cfg *config.Config) (*ServiceBus, error) {
	connStr := cfg.AzureServiceBus.ConnectionString

	client, err := azservicebus.NewClientFromConnectionString(connStr, nil)

	if err != nil {
		return nil, err
	}

	return &ServiceBus{
		Client: client,
	}, nil
}

func (r *ServiceBus) Close() {
	_ = r.Client.Close(context.Background())
}
