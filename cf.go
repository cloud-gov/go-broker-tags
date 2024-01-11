package brokertags

import (
	"context"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type NameResolver interface {
	getOrganization(organizationGUID string) (*resource.Organization, error)
	getSpace(spaceGUID string) (*resource.Space, error)
	getServiceInstance(instanceGUID string) (*resource.ServiceInstance, error)
}

type OrganizationGetter interface {
	Get(ctx context.Context, guid string) (*resource.Organization, error)
}

type SpaceGetter interface {
	Get(ctx context.Context, guid string) (*resource.Space, error)
}

type ServiceInstanceGetter interface {
	Get(ctx context.Context, guid string) (*resource.ServiceInstance, error)
}

type cfNameResolver struct {
	Organizations    OrganizationGetter
	Spaces           SpaceGetter
	ServiceInstances ServiceInstanceGetter
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
		Spaces:           cf.Spaces,
		ServiceInstances: cf.ServiceInstances,
	}, nil
}

func (c *cfNameResolver) getOrganization(organizationGUID string) (*resource.Organization, error) {
	organization, err := c.Organizations.Get(context.Background(), organizationGUID)
	if err != nil {
		return nil, err
	}
	return organization, nil
}

func (c *cfNameResolver) getSpace(spaceGUID string) (*resource.Space, error) {
	space, err := c.Spaces.Get(context.Background(), spaceGUID)
	if err != nil {
		return nil, err
	}
	return space, nil
}

func (c *cfNameResolver) getServiceInstance(instanceGUID string) (*resource.ServiceInstance, error) {
	instance, err := c.ServiceInstances.Get(context.Background(), instanceGUID)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
