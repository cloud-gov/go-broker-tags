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
	ServiceNameTagKey         = "Service offering name"
	ServicePlanName           = "Service plan name"
	SpaceGUIDTagKey           = "Space GUID"
	SpaceNameTagKey           = "Space name"
)

type TagManager interface {
	GenerateTags(
		action Action,
		serviceName string,
		servicePlanName string,
		organizationGUID string,
		spaceGUID string,
		instanceGUID string,
	) (map[string]string, error)
}

type CfTagManager struct {
	broker         string
	environment    string
	cfNameResolver NameResolver
}

func NewCFTagManager(
	broker string,
	environment string,
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
		environment,
		cfNameResolver,
	}, nil
}

func (t *CfTagManager) GenerateTags(
	action Action,
	serviceName string,
	planName string,
	instanceGUID string,
	spaceGUID string,
	organizationGUID string,
) (map[string]string, error) {
	tags := make(map[string]string)

	tags[ClientTagKey] = "Cloud Foundry"

	tags[action.getTagKey()] = time.Now().Format(time.RFC3339)

	if t.broker != "" {
		tags[BrokerTagKey] = t.broker
	}

	if t.environment != "" {
		tags[EnvironmentTagKey] = strings.ToLower(t.environment)
	}

	if serviceName != "" {
		tags[ServiceNameTagKey] = serviceName
	}

	if planName != "" {
		tags[ServicePlanName] = planName
	}

	if instanceGUID != "" {
		tags[ServiceInstanceGUIDTagKey] = instanceGUID
	}

	if instanceGUID != "" && spaceGUID == "" {
		instance, err := t.cfNameResolver.getServiceInstance(instanceGUID)
		if err != nil {
			return nil, err
		}
		spaceGUID = instance.Relationships.Space.Data.GUID
	}

	if spaceGUID != "" {
		tags[SpaceGUIDTagKey] = spaceGUID

		spaceName, err := t.cfNameResolver.getSpaceName(spaceGUID)
		if err != nil {
			return nil, err
		}
		tags[SpaceNameTagKey] = spaceName
	}

	if organizationGUID != "" {
		tags[OrganizationGUIDTagKey] = organizationGUID

		organizationName, err := t.cfNameResolver.getOrganizationName(organizationGUID)
		if err != nil {
			return nil, err
		}
		tags[OrganizationNameTagKey] = organizationName
	}

	return tags, nil
}
