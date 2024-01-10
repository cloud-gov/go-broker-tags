package brokertags

import (
	"strings"
	"time"
)

const (
	BrokerTagKey              = "broker"
	ClientTagKey              = "client"
	EnvironmentTagKey         = "environment"
	OrganizationGUIDTagKey    = "Organization GUID"
	OrganizationNameTagKey    = "Organization name"
	ServiceInstanceGUIDTagKey = "Instance GUID"
	ServiceInstanceNameTagKey = "Instance name"
	ServiceNameTagKey         = "Service offering name"
	ServicePlanName           = "Service plan name"
	SpaceGUIDTagKey           = "Space GUID"
	SpaceNameTagKey           = "Space name"
)

type TagManager interface {
	GenerateTags(
		action Action,
		environment string,
		serviceName string,
		servicePlanName string,
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
	environment string,
	serviceName string,
	planName string,
	organizationGUID string,
	spaceGUID string,
	instanceGUID string,
) (map[string]string, error) {
	tags := make(map[string]string)

	tags[ClientTagKey] = "Cloud Foundry"

	tags[action.getTagKey()] = time.Now().Format(time.RFC3339)

	if t.broker != "" {
		tags[BrokerTagKey] = t.broker
	}

	if environment != "" {
		tags[EnvironmentTagKey] = strings.ToLower(environment)
	}

	if serviceName != "" {
		tags[ServiceNameTagKey] = serviceName
	}

	if planName != "" {
		tags[ServicePlanName] = planName
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

		spaceName, err := t.cfNameResolver.getSpaceName(spaceGUID)
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
