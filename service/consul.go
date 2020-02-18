package service

import (
	"consul-service-tag-web/utils"

	"github.com/hashicorp/consul/api"
)

// ConsulClient ...
type ConsulClient struct {
	client *api.Client
}

// ConsulService ...
type ConsulService struct {
	FilterTag string
	Service   []*ConsulServiceTag
}

// ConsulServiceTag ...
type ConsulServiceTag struct {
	Name string
	Tags []string
}

// NewConsulClient ...
func NewConsulClient() (*ConsulClient, error) {
	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())

	if err != nil {
		return nil, err
	}

	return &ConsulClient{client: client}, nil
}

// GetConsulServices ...
func (c *ConsulClient) GetConsulServices(filterTags []string) ([]*ConsulService, error) {

	consulService := []*ConsulService{}
	consulServiceTag := []*ConsulServiceTag{}

	// Bind to Catalog
	consulCatalog := c.client.Catalog()

	// Query consul for services
	svcs, _, err := consulCatalog.Services(nil)
	if err != nil {
		return nil, err
	}

	// Range over each tag to filter
	for _, tag := range filterTags {

		consulServiceTag = []*ConsulServiceTag{}

		// Range over each service
		for k, v := range svcs {

			// Check if tag to filter matches with current one
			if _, isFound := utils.Find(v, tag); !isFound {
				// Tag not found, loop over next iteration
				continue
			}

			// Store service tags
			consulServiceTag = append(consulServiceTag, &ConsulServiceTag{
				Name: k,
				Tags: v,
			})

		}

		// Store service
		consulService = append(consulService, &ConsulService{
			FilterTag: tag,
			Service:   consulServiceTag,
		})
	}

	return consulService, nil
}
