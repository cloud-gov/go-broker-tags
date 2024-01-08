package brokertags

import (
	"context"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type NameResolver interface {
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

type cfNameResolver struct {
	Organizations    OrganizationGetter
	ServiceInstances ServiceInstanceGetter
	ServiceOfferings ServiceOfferingGetter
	ServicePlans     ServicePlanGetter
	Spaces           SpaceGetter
}

func newCFNameResolver(
	cfApiUrl string,
	cfApiClientId string,
	cfApiClientSecret string,
) (*cfNameResolver, error) {
	cfg, err := config.NewClientSecret(
		cfApiUrl,
		cfApiClientId,
		cfApiClientSecret,
	)
	if err != nil {
		return nil, err
	}
	cf, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	return &cfNameResolver{
		Organizations:    cf.Organizations,
		ServiceInstances: cf.ServiceInstances,
		ServiceOfferings: cf.ServiceOfferings,
		ServicePlans:     cf.ServicePlans,
		Spaces:           cf.Spaces,
	}, nil
}

func (c *cfNameResolver) getServiceInstanceName(instanceGUID string) (string, error) {
	instance, err := c.ServiceInstances.Get(context.Background(), instanceGUID)
	if err != nil {
		return "", err
	}
	return instance.Name, nil
}

func (c *cfNameResolver) getServiceOfferingName(serviceGUID string) (string, error) {
	service, err := c.ServiceOfferings.Get(context.Background(), serviceGUID)
	if err != nil {
		return "", err
	}
	return service.Name, nil
}

func (c *cfNameResolver) getServicePlanName(servicePlanGUID string) (string, error) {
	servicePlan, err := c.ServicePlans.Get(context.Background(), servicePlanGUID)
	if err != nil {
		return "", err
	}
	return servicePlan.Name, nil
}

func (c *cfNameResolver) getOrganizationName(organizationGUID string) (string, error) {
	organization, err := c.Organizations.Get(context.Background(), organizationGUID)
	if err != nil {
		return "", err
	}
	return organization.Name, nil
}

func (c *cfNameResolver) getSpaceName(spaceGUID string) (string, error) {
	space, err := c.Spaces.Get(context.Background(), spaceGUID)
	if err != nil {
		return "", err
	}
	return space.Name, nil
}
