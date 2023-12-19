package sharedBrokerUtils

import (
	"context"
	"os"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
)

type CFClientWrapper interface {
	getServiceOfferingName(serviceGUID string) (string, error)
	getServicePlanName(servicePlanGUID string) (string, error)
	getOrganizationName(organizationGUID string) (string, error)
	getSpaceName(spaceGUID string) (string, error)
	getInstanceName(instanceGUID string) (string, error)
}

type cfClientWrapper struct {
	cf *client.Client
}

func NewCFClient() (*cfClientWrapper, error) {
	cfg, err := config.NewClientSecret(
		os.Getenv("CF_API_URL"),
		os.Getenv("CF_API_CLIENT_ID"),
		os.Getenv("CF_API_CLIENT_SECRET"),
	)
	if err != nil {
		return nil, err
	}
	cf, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	return &cfClientWrapper{
		cf: cf,
	}, nil
}

func (c *cfClientWrapper) getServiceOfferingName(serviceGUID string) (string, error) {
	service, err := c.cf.ServiceOfferings.Get(context.Background(), serviceGUID)
	if err != nil {
		return "", err
	}
	return service.Name, nil
}

func (c *cfClientWrapper) getServicePlanName(servicePlanGUID string) (string, error) {
	servicePlan, err := c.cf.ServicePlans.Get(context.Background(), servicePlanGUID)
	if err != nil {
		return "", err
	}
	return servicePlan.Name, nil
}

func (c *cfClientWrapper) getOrganizationName(organizationGUID string) (string, error) {
	organization, err := c.cf.Organizations.Get(context.Background(), organizationGUID)
	if err != nil {
		return "", err
	}
	return organization.Name, nil
}

func (c *cfClientWrapper) getSpaceName(spaceGUID string) (string, error) {
	space, err := c.cf.Spaces.Get(context.Background(), spaceGUID)
	if err != nil {
		return "", err
	}
	return space.Name, nil
}

func (c *cfClientWrapper) getInstanceName(instanceGUID string) (string, error) {
	instance, err := c.cf.ServiceInstances.Get(context.Background(), instanceGUID)
	if err != nil {
		return "", err
	}
	return instance.Name, nil
}
