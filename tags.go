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

type ResourceGUIDs struct {
	instanceGUID     string
	spaceGUID        string
	organizationGUID string
}

func (t *CfTagManager) GenerateTags(
	action Action,
	serviceName string,
	planName string,
	resourceGUIDs ResourceGUIDs,
	getMissingResources bool,
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

	if resourceGUIDs.instanceGUID != "" {
		tags[ServiceInstanceGUIDTagKey] = resourceGUIDs.instanceGUID
	}

	var (
		spaceGUID        string
		organizationGUID string
		err              error
	)

	spaceGUID = resourceGUIDs.spaceGUID
	if spaceGUID == "" && resourceGUIDs.instanceGUID != "" && getMissingResources {
		spaceGUID, err = t.getSpaceGuid(resourceGUIDs.instanceGUID)
		if err != nil {
			return nil, err
		}
	}

	if spaceGUID != "" {
		tags[SpaceGUIDTagKey] = spaceGUID

		space, err := t.cfResourceGetter.getSpace(spaceGUID)
		if err != nil {
			return nil, err
		}
		tags[SpaceNameTagKey] = space.Name
	}

	organizationGUID = resourceGUIDs.organizationGUID
	if organizationGUID == "" && spaceGUID != "" && getMissingResources {
		organizationGUID, err = t.getOrganizationGuid(spaceGUID)
		if err != nil {
			return nil, err
		}
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

func (t *CfTagManager) getSpaceGuid(instanceGUID string) (string, error) {
	instance, err := t.cfResourceGetter.getServiceInstance(instanceGUID)
	if err != nil {
		return "", err
	}
	spaceGUID := instance.Relationships.Space.Data.GUID
	return spaceGUID, nil
}

func (t *CfTagManager) getOrganizationGuid(spaceGUID string) (string, error) {
	space, err := t.cfResourceGetter.getSpace(spaceGUID)
	if err != nil {
		return "", err
	}
	organizationGUID := space.Relationships.Organization.Data.GUID
	return organizationGUID, nil
}
