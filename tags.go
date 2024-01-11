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
	broker           string
	environment      string
	cfResourceGetter ResourceGetter
}

func NewCFTagManager(
	broker string,
	environment string,
	cfApiUrl string,
	cfApiClientId string,
	cfApiClientSecret string,
) (*CfTagManager, error) {
	cfResourceGetter, err := newCFResourceGetter(
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
		cfResourceGetter,
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

	spaceGUID, err := t.getSpaceGuid(spaceGUID, instanceGUID)
	if err != nil {
		return nil, err
	}

	if spaceGUID != "" {
		tags[SpaceGUIDTagKey] = spaceGUID

		space, err := t.cfResourceGetter.getSpace(spaceGUID)
		if err != nil {
			return nil, err
		}
		tags[SpaceNameTagKey] = space.Name
	}

	organizationGUID, err = t.getOrganizationGuid(organizationGUID, spaceGUID)
	if err != nil {
		return nil, err
	}

	if organizationGUID != "" {
		tags[OrganizationGUIDTagKey] = organizationGUID

		organization, err := t.cfResourceGetter.getOrganization(organizationGUID)
		if err != nil {
			return nil, err
		}
		tags[OrganizationNameTagKey] = organization.Name
	}

	return tags, nil
}

func (t *CfTagManager) getSpaceGuid(
	spaceGUID string,
	instanceGUID string,
) (string, error) {
	if spaceGUID != "" {
		return spaceGUID, nil
	}
	if instanceGUID != "" {
		instance, err := t.cfResourceGetter.getServiceInstance(instanceGUID)
		if err != nil {
			return spaceGUID, err
		}
		spaceGUID = instance.Relationships.Space.Data.GUID
	}
	return spaceGUID, nil
}

func (t *CfTagManager) getOrganizationGuid(
	organizationGUID string,
	spaceGUID string,
) (string, error) {
	if organizationGUID != "" {
		return organizationGUID, nil
	}
	if spaceGUID != "" {
		space, err := t.cfResourceGetter.getSpace(spaceGUID)
		if err != nil {
			return organizationGUID, err
		}
		organizationGUID = space.Relationships.Organization.Data.GUID
	}
	return organizationGUID, nil
}
