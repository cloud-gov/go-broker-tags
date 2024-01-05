package brokerTags

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

type TagManager struct {
	broker          string
	cfClientWrapper CFClientWrapper
}

func NewTagManager() (*TagManager, error) {
	cfClientWrapper, err := NewCFClientWrapper()
	if err != nil {
		return nil, err
	}
	return &TagManager{
		cfClientWrapper: cfClientWrapper,
	}, nil
}

func (t *TagManager) GenerateTags(
	action string, // The action that about to occur for the tagged resource, e.g. "created", "updated"
	serviceGUID string,
	servicePlanGUID string,
	organizationGUID string,
	spaceGUID string,
	instanceGUID string,
) (map[string]string, error) {
	tags := make(map[string]string)

	tags[ClientTagKey] = "Cloud Foundry"

	tags[BrokerTagKey] = t.broker

	tags[action+" at"] = time.Now().Format(time.RFC822Z)

	if serviceGUID != "" {
		tags[ServiceOfferingGUIDTagKey] = serviceGUID

		serviceOfferingName, err := t.cfClientWrapper.getServiceOfferingName(serviceGUID)
		if err != nil {
			return nil, err
		}
		tags[ServiceOfferingNameTagKey] = serviceOfferingName
	}

	if servicePlanGUID != "" {
		tags[ServicePlanGUIDTagKey] = servicePlanGUID

		servicePlanName, err := t.cfClientWrapper.getServicePlanName(servicePlanGUID)
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
