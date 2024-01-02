package brokerTags

import (
	"time"
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
	action string,
	serviceGUID string,
	servicePlanGUID string,
	organizationGUID string,
	spaceGUID string,
	instanceGUID string,
) (map[string]string, error) {
	tags := make(map[string]string)

	tags["client"] = "Cloud Foundry"

	tags["broker"] = t.broker

	tags[action+" at"] = time.Now().Format(time.RFC822Z)

	if serviceGUID != "" {
		tags["Service GUID"] = serviceGUID

		serviceOfferingName, err := t.cfClientWrapper.getServiceOfferingName(serviceGUID)
		if err != nil {
			return nil, err
		}
		tags["Service offering name"] = serviceOfferingName
	}

	if servicePlanGUID != "" {
		tags["Plan GUID"] = servicePlanGUID

		servicePlanName, err := t.cfClientWrapper.getServicePlanName(servicePlanGUID)
		if err != nil {
			return nil, err
		}
		tags["Service plan name"] = servicePlanName
	}

	if organizationGUID != "" {
		tags["Organization GUID"] = organizationGUID

		organizationName, err := t.cfClientWrapper.getOrganizationName(organizationGUID)
		if err != nil {
			return nil, err
		}
		tags["Organization name"] = organizationName
	}

	if spaceGUID != "" {
		tags["Space GUID"] = spaceGUID

		spaceName, err := t.cfClientWrapper.getSpaceName(organizationGUID)
		if err != nil {
			return nil, err
		}
		tags["Space name"] = spaceName
	}

	if instanceGUID != "" {
		tags["Instance GUID"] = instanceGUID

		instanceName, err := t.cfClientWrapper.getInstanceName(instanceGUID)
		if err != nil {
			return nil, err
		}
		tags["Instance name"] = instanceName
	}

	return tags, nil
}
