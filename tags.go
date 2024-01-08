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
	ServiceOfferingGUIDTagKey = "Service GUID"
	ServiceOfferingNameTagKey = "Service offering name"
	ServicePlanGUIDTagKey     = "Plan GUID"
	ServicePlanNameTagKey     = "Service plan name"
	SpaceGUIDTagKey           = "Space GUID"
	SpaceNameTagKey           = "Space name"
)

type TagGenerator interface {
	GenerateTags(
		action Action,
		serviceGUID string,
		servicePlanGUID string,
		organizationGUID string,
		spaceGUID string,
		instanceGUID string,
	) (map[string]string, error)
}

type TagManager struct {
	broker         string
	cfNameResolver NameResolver
}

func NewManager(
	cfApiUrl string,
	cfApiClientId string,
	cfApiClientSecret string,
) (*TagManager, error) {
	cfNameResolver, err := newCFNameResolver(
		cfApiUrl,
		cfApiClientId,
		cfApiClientSecret,
	)
	if err != nil {
		return nil, err
	}
	return &TagManager{
		cfNameResolver: cfNameResolver,
	}, nil
}

func (t *TagManager) GenerateTags(
	action Action,
	serviceGUID string,
	servicePlanGUID string,
	organizationGUID string,
	spaceGUID string,
	instanceGUID string,
) (map[string]string, error) {
	tags := make(map[string]string)

	tags[ClientTagKey] = "Cloud Foundry"

	tags[BrokerTagKey] = t.broker

	tags[action.getTagKey()] = time.Now().Format(time.RFC3339)

	if serviceGUID != "" {
		tags[ServiceOfferingGUIDTagKey] = serviceGUID

		serviceOfferingName, err := t.cfNameResolver.getServiceOfferingName(serviceGUID)
		if err != nil {
			return nil, err
		}
		tags[ServiceOfferingNameTagKey] = serviceOfferingName
	}

	if servicePlanGUID != "" {
		tags[ServicePlanGUIDTagKey] = servicePlanGUID

		servicePlanName, err := t.cfNameResolver.getServicePlanName(servicePlanGUID)
		if err != nil {
			return nil, err
		}
		tags[ServicePlanNameTagKey] = servicePlanName
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
