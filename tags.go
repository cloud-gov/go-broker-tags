package brokerTags

import (
	"time"
)

const (
	BrokerTagKey              = "broker"
	ClientTagKey              = "client"
	CreatedAtTagKey           = "Created at"
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
	UpdatedAtTagKey           = "Updated at"
)

// Action - Custom type to hold value for broker action
type Action int

const (
	Create Action = iota // EnumIndex = 0
	Update               // EnumIndex = 1
)

func (a Action) String() string {
	return [...]string{"Created", "Updated"}[a]
}

type TagManager struct {
	broker         string
	cfNameResolver NameResolver
}

func NewTagManager() (*TagManager, error) {
	cfNameResolver, err := newCFNameResolver()
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

	if action == Create {
		tags[CreatedAtTagKey] = time.Now().Format(time.RFC822Z)
	} else if action == Update {
		tags[UpdatedAtTagKey] = time.Now().Format(time.RFC822Z)
	}

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
