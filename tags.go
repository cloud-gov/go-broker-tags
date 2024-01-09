package brokertags

import (
	"time"
)

const (
	BrokerTagKey              = "broker"
	ClientTagKey              = "client"
	OrganizationGUIDTagKey    = "Organization GUID"
	OrganizationNameTagKey    = "Organization name"
	ServiceInstanceGUIDTagKey = "Instance GUID"
	ServiceInstanceNameTagKey = "Instance name"
	ServiceIDTagKey           = "Service ID"
	ServicePlanIdTagKey       = "Plan ID"
	SpaceGUIDTagKey           = "Space GUID"
	SpaceNameTagKey           = "Space name"
)

type TagManager interface {
	GenerateTags(
		action Action,
		serviceGUID string,
		servicePlanGUID string,
		organizationGUID string,
		spaceGUID string,
		instanceGUID string,
	) (map[string]string, error)
}

type CfTagManager struct {
	broker         string
	cfNameResolver NameResolver
}

func NewCFTagManager(
	broker string,
	cfApiUrl string,
	cfApiClientId string,
	cfApiClientSecret string,
) (*CfTagManager, error) {
	cfNameResolver, err := newCFNameResolver(
		cfApiUrl,
		cfApiClientId,
		cfApiClientSecret,
	)
	if err != nil {
		return nil, err
	}
	return &CfTagManager{
		broker,
		cfNameResolver,
	}, nil
}

func (t *CfTagManager) GenerateTags(
	action Action,
	serviceID string,
	planID string,
	organizationGUID string,
	spaceGUID string,
	instanceGUID string,
) (map[string]string, error) {
	tags := make(map[string]string)

	tags[ClientTagKey] = "Cloud Foundry"

	tags[BrokerTagKey] = t.broker

	tags[action.getTagKey()] = time.Now().Format(time.RFC3339)

	if serviceID != "" {
		tags[ServiceIDTagKey] = serviceID
	}

	if planID != "" {
		tags[ServicePlanIdTagKey] = planID
	}

	if organizationGUID != "" {
		tags[OrganizationGUIDTagKey] = organizationGUID

		organizationName, err := t.cfNameResolver.getOrganizationName(organizationGUID)
		if err != nil {
			return nil, err
		}
		tags[OrganizationNameTagKey] = organizationName
	}

	if spaceGUID != "" {
		tags[SpaceGUIDTagKey] = spaceGUID

		spaceName, err := t.cfNameResolver.getSpaceName(organizationGUID)
		if err != nil {
			return nil, err
		}
		tags[SpaceNameTagKey] = spaceName
	}

	if instanceGUID != "" {
		tags[ServiceInstanceGUIDTagKey] = instanceGUID

		instanceName, err := t.cfNameResolver.getServiceInstanceName(instanceGUID)
		if err != nil {
			return nil, err
		}
		tags[ServiceInstanceNameTagKey] = instanceName
	}

	return tags, nil
}
