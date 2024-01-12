package brokertags

import (
	"errors"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/google/go-cmp/cmp"
)

type mockCFClientWrapper struct {
	getOrganizationErr          error
	organizationName            string
	getSpaceErr                 error
	spaceName                   string
	getServiceInstanceErr       error
	instanceName                string
	spaceGUID                   string
	organizationGUID            string
	instanceGUID                string
	getServiceInstanceCallCount int
	getSpaceInstanceCallCount   int
}

func (m *mockCFClientWrapper) getOrganization(organizationGUID string) (*resource.Organization, error) {
	if m.getOrganizationErr != nil {
		return nil, m.getOrganizationErr
	}
	if m.organizationGUID != "" && m.organizationGUID != organizationGUID {
		return nil, errors.New("organization GUID does not match expected value")
	}
	return &resource.Organization{
		Name: m.organizationName,
	}, nil
}

func (m *mockCFClientWrapper) getSpace(spaceGUID string) (*resource.Space, error) {
	m.getSpaceInstanceCallCount++
	if m.getSpaceErr != nil {
		return nil, m.getSpaceErr
	}
	if m.spaceGUID != "" && m.spaceGUID != spaceGUID {
		return nil, errors.New("space GUID does not match expected value")
	}
	return &resource.Space{
		Name: m.spaceName,
		Relationships: &resource.SpaceRelationships{
			Organization: &resource.ToOneRelationship{
				Data: &resource.Relationship{
					GUID: m.organizationGUID,
				},
			},
		},
	}, nil
}

func (m *mockCFClientWrapper) getServiceInstance(instanceGUID string) (*resource.ServiceInstance, error) {
	m.getServiceInstanceCallCount++
	if m.getServiceInstanceErr != nil {
		return nil, m.getServiceInstanceErr
	}
	if m.instanceGUID != "" && m.instanceGUID != instanceGUID {
		return nil, errors.New("instance GUID does not match expected value")
	}
	return &resource.ServiceInstance{
		Name: m.instanceName,
		Relationships: resource.ServiceInstanceRelationships{
			Space: &resource.ToOneRelationship{
				Data: &resource.Relationship{
					GUID: m.spaceGUID,
				},
			},
		},
	}, nil
}

func TestGenerateTags(t *testing.T) {
	testCases := map[string]struct {
		tagManager                          *CfTagManager
		expectedTags                        map[string]string
		action                              Action
		serviceOfferingName                 string
		servicePlanName                     string
		resourceGUIDS                       ResourceGUIDs
		getMissingResources                 bool
		expectedGetServiceInstanceCallCount int
		expectedGetSpaceInstanceCallCount   int
	}{
		"Create": {
			action:              Create,
			serviceOfferingName: "abc1",
			servicePlanName:     "abc2",
			resourceGUIDS: ResourceGUIDs{
				OrganizationGUID: "abc3",
				SpaceGUID:        "abc4",
				InstanceGUID:     "abc5",
			},
			tagManager: &CfTagManager{
				broker:      "AWS Broker",
				environment: "testing",
				cfResourceGetter: &mockCFClientWrapper{
					organizationName: "org-1",
					spaceName:        "space-1",
					spaceGUID:        "abc4",
					organizationGUID: "abc3",
					instanceGUID:     "abc5",
				},
			},
			expectedGetSpaceInstanceCallCount: 1,
			expectedTags: map[string]string{
				"client":                "Cloud Foundry",
				"broker":                "AWS Broker",
				"environment":           "testing",
				"Service offering name": "abc1",
				"Service plan name":     "abc2",
				"Organization GUID":     "abc3",
				"Space GUID":            "abc4",
				"Instance GUID":         "abc5",
				"Organization name":     "org-1",
				"Space name":            "space-1",
			},
		},
		"Update": {
			action:              Update,
			serviceOfferingName: "abc1",
			servicePlanName:     "abc2",
			resourceGUIDS: ResourceGUIDs{
				OrganizationGUID: "abc3",
				SpaceGUID:        "abc4",
				InstanceGUID:     "abc5",
			},
			tagManager: &CfTagManager{
				broker:      "AWS Broker",
				environment: "testing",
				cfResourceGetter: &mockCFClientWrapper{
					organizationName: "org-1",
					spaceName:        "space-1",
					spaceGUID:        "abc4",
					organizationGUID: "abc3",
					instanceGUID:     "abc5",
				},
			},
			expectedGetSpaceInstanceCallCount: 1,
			expectedTags: map[string]string{
				"client":                "Cloud Foundry",
				"broker":                "AWS Broker",
				"environment":           "testing",
				"Service offering name": "abc1",
				"Service plan name":     "abc2",
				"Organization GUID":     "abc3",
				"Space GUID":            "abc4",
				"Instance GUID":         "abc5",
				"Organization name":     "org-1",
				"Space name":            "space-1",
			},
		},
		"no broker name": {
			action:              Create,
			serviceOfferingName: "abc1",
			servicePlanName:     "abc2",
			resourceGUIDS: ResourceGUIDs{
				OrganizationGUID: "abc3",
				SpaceGUID:        "abc4",
				InstanceGUID:     "abc5",
			},
			tagManager: &CfTagManager{
				environment: "testing",
				cfResourceGetter: &mockCFClientWrapper{
					organizationName: "org-1",
					spaceName:        "space-1",
					spaceGUID:        "abc4",
					organizationGUID: "abc3",
					instanceGUID:     "abc5",
				},
			},
			expectedGetSpaceInstanceCallCount: 1,
			expectedTags: map[string]string{
				"client":                "Cloud Foundry",
				"environment":           "testing",
				"Service offering name": "abc1",
				"Service plan name":     "abc2",
				"Organization GUID":     "abc3",
				"Space GUID":            "abc4",
				"Instance GUID":         "abc5",
				"Organization name":     "org-1",
				"Space name":            "space-1",
			},
		},
		"no environment tag": {
			action:              Create,
			serviceOfferingName: "abc1",
			servicePlanName:     "abc2",
			resourceGUIDS: ResourceGUIDs{
				OrganizationGUID: "abc3",
				SpaceGUID:        "abc4",
				InstanceGUID:     "abc5",
			},
			tagManager: &CfTagManager{
				broker: "AWS Broker",
				cfResourceGetter: &mockCFClientWrapper{
					organizationName: "org-1",
					spaceName:        "space-1",
					spaceGUID:        "abc4",
					organizationGUID: "abc3",
					instanceGUID:     "abc5",
				},
			},
			expectedGetSpaceInstanceCallCount: 1,
			expectedTags: map[string]string{
				"client":                "Cloud Foundry",
				"broker":                "AWS Broker",
				"Service offering name": "abc1",
				"Service plan name":     "abc2",
				"Organization GUID":     "abc3",
				"Space GUID":            "abc4",
				"Instance GUID":         "abc5",
				"Organization name":     "org-1",
				"Space name":            "space-1",
			},
		},
		"get missing organization": {
			action:              Create,
			serviceOfferingName: "abc1",
			servicePlanName:     "abc2",
			resourceGUIDS: ResourceGUIDs{
				SpaceGUID:    "abc4",
				InstanceGUID: "abc5",
			},
			getMissingResources: true,
			tagManager: &CfTagManager{
				broker:      "AWS Broker",
				environment: "testing",
				cfResourceGetter: &mockCFClientWrapper{
					organizationGUID: "abc3",
					organizationName: "org-1",
					spaceName:        "space-1",
					spaceGUID:        "abc4",
					instanceGUID:     "abc5",
				},
			},
			expectedGetSpaceInstanceCallCount: 1,
			expectedTags: map[string]string{
				"client":                "Cloud Foundry",
				"broker":                "AWS Broker",
				"environment":           "testing",
				"Service offering name": "abc1",
				"Service plan name":     "abc2",
				"Organization GUID":     "abc3",
				"Space GUID":            "abc4",
				"Instance GUID":         "abc5",
				"Organization name":     "org-1",
				"Space name":            "space-1",
			},
		},
		"get missing space and organization": {
			action:              Create,
			serviceOfferingName: "abc1",
			servicePlanName:     "abc2",
			resourceGUIDS: ResourceGUIDs{
				InstanceGUID: "abc5",
			},
			getMissingResources: true,
			tagManager: &CfTagManager{
				broker:      "AWS Broker",
				environment: "testing",
				cfResourceGetter: &mockCFClientWrapper{
					spaceGUID:        "abc4",
					organizationGUID: "abc3",
					organizationName: "org-1",
					spaceName:        "space-1",
					instanceGUID:     "abc5",
				},
			},
			expectedGetSpaceInstanceCallCount:   1,
			expectedGetServiceInstanceCallCount: 1,
			expectedTags: map[string]string{
				"client":                "Cloud Foundry",
				"broker":                "AWS Broker",
				"environment":           "testing",
				"Service offering name": "abc1",
				"Service plan name":     "abc2",
				"Organization GUID":     "abc3",
				"Space GUID":            "abc4",
				"Instance GUID":         "abc5",
				"Organization name":     "org-1",
				"Space name":            "space-1",
			},
		},
		"no organization GUID, do not get missing resource": {
			action:              Create,
			serviceOfferingName: "abc1",
			servicePlanName:     "abc2",
			resourceGUIDS: ResourceGUIDs{
				SpaceGUID:    "abc4",
				InstanceGUID: "abc5",
			},
			tagManager: &CfTagManager{
				broker:      "AWS Broker",
				environment: "testing",
				cfResourceGetter: &mockCFClientWrapper{
					spaceGUID:    "abc4",
					spaceName:    "space-1",
					instanceGUID: "abc5",
				},
			},
			expectedGetSpaceInstanceCallCount:   1,
			expectedGetServiceInstanceCallCount: 0,
			expectedTags: map[string]string{
				"client":                "Cloud Foundry",
				"broker":                "AWS Broker",
				"environment":           "testing",
				"Service offering name": "abc1",
				"Service plan name":     "abc2",
				"Space GUID":            "abc4",
				"Instance GUID":         "abc5",
				"Space name":            "space-1",
			},
		},
		"no space or organization GUID, do not get missing resource": {
			action:              Create,
			serviceOfferingName: "abc1",
			servicePlanName:     "abc2",
			resourceGUIDS: ResourceGUIDs{
				InstanceGUID: "abc5",
			},
			tagManager: &CfTagManager{
				broker:      "AWS Broker",
				environment: "testing",
				cfResourceGetter: &mockCFClientWrapper{
					instanceGUID: "abc5",
				},
			},
			expectedGetSpaceInstanceCallCount:   0,
			expectedGetServiceInstanceCallCount: 0,
			expectedTags: map[string]string{
				"client":                "Cloud Foundry",
				"broker":                "AWS Broker",
				"environment":           "testing",
				"Service offering name": "abc1",
				"Service plan name":     "abc2",
				"Instance GUID":         "abc5",
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			tags, err := test.tagManager.GenerateTags(
				test.action,
				test.serviceOfferingName,
				test.servicePlanName,
				test.resourceGUIDS,
				test.getMissingResources,
			)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			actionTagKey := test.action.getTagKey()
			if tags[actionTagKey] == "" {
				t.Fatalf("Expected a value for %s tag", actionTagKey)
			}
			delete(tags, actionTagKey)

			if !cmp.Equal(tags, test.expectedTags) {
				t.Errorf(cmp.Diff(tags, test.expectedTags))
			}

			mockCfResourceGetter, ok := test.tagManager.cfResourceGetter.(*mockCFClientWrapper)
			if ok {
				if test.expectedGetServiceInstanceCallCount != mockCfResourceGetter.getServiceInstanceCallCount {
					t.Fatalf("Expected %d calls to getServiceInstance, got %d", test.expectedGetServiceInstanceCallCount, mockCfResourceGetter.getServiceInstanceCallCount)
				}
				if test.expectedGetSpaceInstanceCallCount != mockCfResourceGetter.getSpaceInstanceCallCount {
					t.Fatalf("Expected %d calls to getSpace, got %d", test.expectedGetSpaceInstanceCallCount, mockCfResourceGetter.getSpaceInstanceCallCount)
				}
			}
		})
	}
}

func TestGenerateTagsHandleErrors(t *testing.T) {
	testCases := map[string]struct {
		tagManager  *CfTagManager
		expectedErr error
	}{
		"error getting organization": {
			tagManager: &CfTagManager{
				cfResourceGetter: &mockCFClientWrapper{
					getOrganizationErr: errors.New("error getting organization"),
				},
			},
			expectedErr: errors.New("error getting organization"),
		},
		"error getting space": {
			tagManager: &CfTagManager{
				broker: "AWS Broker",
				cfResourceGetter: &mockCFClientWrapper{
					getSpaceErr: errors.New("error getting space"),
				},
			},
			expectedErr: errors.New("error getting space"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := test.tagManager.GenerateTags(
				Create,
				"abc1",
				"abc2",
				ResourceGUIDs{
					OrganizationGUID: "org-1",
					InstanceGUID:     "instance-1",
					SpaceGUID:        "space-1",
				},
				false,
			)
			if err == nil || err.Error() != test.expectedErr.Error() {
				t.Fatalf("did not received expected err: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}

func TestGetSpaceGuid(t *testing.T) {
	testCases := map[string]struct {
		tagManager   *CfTagManager
		instanceGUID string
		expectedGuid string
		expectedErr  error
	}{
		"success": {
			tagManager: &CfTagManager{
				cfResourceGetter: &mockCFClientWrapper{
					instanceGUID: "instance-1",
					spaceGUID:    "space-1",
				},
			},
			instanceGUID: "instance-1",
			expectedGuid: "space-1",
		},
		"error": {
			tagManager: &CfTagManager{
				cfResourceGetter: &mockCFClientWrapper{
					getServiceInstanceErr: errors.New("error getting service instance"),
				},
			},
			expectedErr: errors.New("error getting service instance"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			spaceGUID, err := test.tagManager.getSpaceGuid(test.instanceGUID)
			if spaceGUID != test.expectedGuid {
				t.Errorf("expected: %s, got: %s", test.expectedGuid, spaceGUID)
			}
			if (test.expectedErr != nil && err == nil) ||
				(err != nil && err.Error() != test.expectedErr.Error()) {
				t.Errorf("expected error: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}

func TestGetOrganizationGuidFromSpace(t *testing.T) {
	testCases := map[string]struct {
		space        *resource.Space
		tagManager   *CfTagManager
		expectedGUID string
	}{
		"success": {
			tagManager: &CfTagManager{
				cfResourceGetter: &mockCFClientWrapper{},
			},
			space: &resource.Space{
				Name: "space-1",
				Relationships: &resource.SpaceRelationships{
					Organization: &resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org-guid-1",
						},
					},
				},
			},
			expectedGUID: "org-guid-1",
		},
		"nil space": {
			tagManager: &CfTagManager{
				cfResourceGetter: &mockCFClientWrapper{},
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			organizationGUID := test.tagManager.getOrganizationGuidFromSpace(test.space)
			if organizationGUID != test.expectedGUID {
				t.Errorf("expected GUID: %s, got: %s", test.expectedGUID, organizationGUID)
			}
		})
	}
}
