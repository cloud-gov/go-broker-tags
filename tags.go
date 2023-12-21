package brokerTags

import (
	"time"
)

type TagManager struct {
	broker   string
	cfClient CFClientWrapper
}

func NewTagManager() (*TagManager, error) {
	cfClient, err := NewCFClient()
	if err != nil {
		return nil, err
	}
	return &TagManager{
		cfClient: cfClient,
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

		serviceOfferingName, err := t.cfClient.getServiceOfferingName(serviceGUID)
		if err != nil {
			return nil, err
		}
		tags["Service offering name"] = serviceOfferingName
	}

	if servicePlanGUID != "" {
		tags["Plan GUID"] = servicePlanGUID

		servicePlanName, err := t.cfClient.getServicePlanName(servicePlanGUID)
		if err != nil {
			return nil, err
		}
		tags["Service plan name"] = servicePlanName
	}

	if organizationGUID != "" {
		tags["Organization GUID"] = organizationGUID

		organizationName, err := t.cfClient.getOrganizationName(organizationGUID)
		if err != nil {
			return nil, err
		}
		tags["Organization name"] = organizationName
	}

	if spaceGUID != "" {
		tags["Space GUID"] = spaceGUID

		spaceName, err := t.cfClient.getSpaceName(organizationGUID)
		if err != nil {
			return nil, err
		}
		tags["Space name"] = spaceName
	}

	if instanceGUID != "" {
		tags["Instance GUID"] = instanceGUID

		instanceName, err := t.cfClient.getInstanceName(instanceGUID)
		if err != nil {
			return nil, err
		}
		tags["Instance name"] = instanceName
	}

	return tags, nil
}
