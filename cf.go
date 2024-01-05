package brokerTags

import (
	"context"
	"os"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type CFClientWrapper interface {
	getOrganizationName(organizationGUID string) (string, error)
	getServiceInstanceName(instanceGUID string) (string, error)
	getServiceOfferingName(serviceGUID string) (string, error)
	getServicePlanName(servicePlanGUID string) (string, error)
	getSpaceName(spaceGUID string) (string, error)
}

type OrganizationGetter interface {
	Get(ctx context.Context, guid string) (*resource.Organization, error)
}

type ServiceInstanceGetter interface {
	Get(ctx context.Context, guid string) (*resource.ServiceInstance, error)
}

type ServiceOfferingGetter interface {
	Get(ctx context.Context, guid string) (*resource.ServiceOffering, error)
}

type ServicePlanGetter interface {
	Get(ctx context.Context, guid string) (*resource.ServicePlan, error)
}

type SpaceGetter interface {
	Get(ctx context.Context, guid string) (*resource.Space, error)
}

type cfClientWrapper struct {
	Organizations    OrganizationGetter
	ServiceInstances ServiceInstanceGetter
	ServiceOfferings ServiceOfferingGetter
	ServicePlans     ServicePlanGetter
	Spaces           SpaceGetter
}

func NewCFClientWrapper() (*cfClientWrapper, error) {
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
		Organizations:    cf.Organizations,
		ServiceInstances: cf.ServiceInstances,
		ServiceOfferings: cf.ServiceOfferings,
		ServicePlans:     cf.ServicePlans,
		Spaces:           cf.Spaces,
	}, nil
}

func (c *cfClientWrapper) getServiceInstanceName(instanceGUID string) (string, error) {
	instance, err := c.ServiceInstances.Get(context.Background(), instanceGUID)
	if err != nil {
		return "", err
	}
	return instance.Name, nil
}

func (c *cfClientWrapper) getServiceOfferingName(serviceGUID string) (string, error) {
	service, err := c.ServiceOfferings.Get(context.Background(), serviceGUID)
	if err != nil {
		return "", err
	}
	return service.Name, nil
}

func (c *cfClientWrapper) getServicePlanName(servicePlanGUID string) (string, error) {
	servicePlan, err := c.ServicePlans.Get(context.Background(), servicePlanGUID)
	if err != nil {
		return "", err
	}
	return servicePlan.Name, nil
}

func (c *cfClientWrapper) getOrganizationName(organizationGUID string) (string, error) {
	organization, err := c.Organizations.Get(context.Background(), organizationGUID)
	if err != nil {
		return "", err
	}
	return organization.Name, nil
}

func (c *cfClientWrapper) getSpaceName(spaceGUID string) (string, error) {
	space, err := c.Spaces.Get(context.Background(), spaceGUID)
	if err != nil {
		return "", err
	}
	return space.Name, nil
}
