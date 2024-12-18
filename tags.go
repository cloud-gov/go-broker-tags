package brokertags

import (
	"strings"
	"time"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
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
		serviceName string,
		servicePlanName string,
		resourceGUIDs ResourceGUIDs,
		getMissingResources bool,
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
	InstanceGUID     string
	SpaceGUID        string
	OrganizationGUID string
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

	var (
		instanceGUID     string
		instance         *resource.ServiceInstance
		spaceGUID        string
		space            *resource.Space
		organizationGUID string
		organization     *resource.Organization
		err              error
	)

	instanceGUID = resourceGUIDs.InstanceGUID
	if instanceGUID != "" {
		tags[ServiceInstanceGUIDTagKey] = instanceGUID

		instance, err = t.cfResourceGetter.getServiceInstance(instanceGUID)
		if err != nil {
			return nil, err
		}
	}

	if instance != nil {
		tags[ServiceInstanceNameTagKey] = instance.Name
	}

	spaceGUID = resourceGUIDs.SpaceGUID
	if spaceGUID == "" && instance != nil {
		spaceGUID = instance.Relationships.Space.Data.GUID
	}

	if spaceGUID != "" {
		tags[SpaceGUIDTagKey] = spaceGUID
		space, err = t.cfResourceGetter.getSpace(spaceGUID)
		if err != nil {
			return nil, err
		}
	}

	if space != nil {
		tags[SpaceNameTagKey] = space.Name
	}

	organizationGUID = resourceGUIDs.OrganizationGUID
	if organizationGUID == "" && getMissingResources {
		organizationGUID = t.getOrganizationGuidFromSpace(space)
	}

	if organizationGUID != "" {
		tags[OrganizationGUIDTagKey] = organizationGUID
		organization, err = t.cfResourceGetter.getOrganization(organizationGUID)
		if err != nil {
			return nil, err
		}
	}

	if organization != nil {
		tags[OrganizationNameTagKey] = organization.Name
	}

	return tags, nil
}

func (t *CfTagManager) getOrganizationGuidFromSpace(space *resource.Space) string {
	if space == nil {
		return ""
	}
	return space.Relationships.Organization.Data.GUID
}
