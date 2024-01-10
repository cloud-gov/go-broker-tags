package brokertags

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type mockCFClientWrapper struct {
	getOrganizationErr error
	organizationName   string
	getSpaceErr        error
	spaceName          string
	getInstanceErr     error
	instanceName       string
	spaceGUID          string
	organizationGUID   string
	instanceGUID       string
}

func (m *mockCFClientWrapper) getOrganizationName(organizationGUID string) (string, error) {
	if m.getOrganizationErr != nil {
		return "", m.getOrganizationErr
	}
	if m.organizationGUID != "" && m.organizationGUID != organizationGUID {
		return "", errors.New("organization GUID does not match expected value")
	}
	return m.organizationName, nil
}

func (m *mockCFClientWrapper) getSpaceName(spaceGUID string) (string, error) {
	if m.getSpaceErr != nil {
		return "", m.getSpaceErr
	}
	if m.spaceGUID != "" && m.spaceGUID != spaceGUID {
		return "", errors.New("space GUID does not match expected value")
	}
	return m.spaceName, nil
}

func (m *mockCFClientWrapper) getServiceInstanceName(instanceGUID string) (string, error) {
	if m.getInstanceErr != nil {
		return "", m.getInstanceErr
	}
	if m.instanceGUID != "" && m.instanceGUID != instanceGUID {
		return "", errors.New("instance GUID does not match expected value")
	}
	return m.instanceName, nil
}

func TestGenerateTags(t *testing.T) {
	testCases := map[string]struct {
		tagManager       *CfTagManager
		expectedTags     map[string]string
		action           Action
		environment      string
		serviceID        string
		servicePlanID    string
		organizationGUID string
		spaceGUID        string
		instanceGUID     string
	}{
		"Create": {
			action:           Create,
			serviceID:        "abc1",
			servicePlanID:    "abc2",
			organizationGUID: "abc3",
			spaceGUID:        "abc4",
			instanceGUID:     "abc5",
			environment:      "testing",
			tagManager: &CfTagManager{
				broker: "AWS Broker",
				cfNameResolver: &mockCFClientWrapper{
					organizationName: "org-1",
					spaceName:        "space-1",
					instanceName:     "instance-1",
					spaceGUID:        "abc4",
					organizationGUID: "abc3",
					instanceGUID:     "abc5",
				},
			},
			expectedTags: map[string]string{
				"client":            "Cloud Foundry",
				"broker":            "AWS Broker",
				"environment":       "testing",
				"Service ID":        "abc1",
				"Plan ID":           "abc2",
				"Organization GUID": "abc3",
				"Space GUID":        "abc4",
				"Instance GUID":     "abc5",
				"Organization name": "org-1",
				"Space name":        "space-1",
				"Instance name":     "instance-1",
			},
		},
		"Update": {
			action:           Update,
			serviceID:        "abc1",
			servicePlanID:    "abc2",
			organizationGUID: "abc3",
			spaceGUID:        "abc4",
			instanceGUID:     "abc5",
			environment:      "testing",
			tagManager: &CfTagManager{
				broker: "AWS Broker",
				cfNameResolver: &mockCFClientWrapper{
					organizationName: "org-1",
					spaceName:        "space-1",
					instanceName:     "instance-1",
					spaceGUID:        "abc4",
					organizationGUID: "abc3",
					instanceGUID:     "abc5",
				},
			},
			expectedTags: map[string]string{
				"client":            "Cloud Foundry",
				"broker":            "AWS Broker",
				"environment":       "testing",
				"Service ID":        "abc1",
				"Plan ID":           "abc2",
				"Organization GUID": "abc3",
				"Space GUID":        "abc4",
				"Instance GUID":     "abc5",
				"Organization name": "org-1",
				"Space name":        "space-1",
				"Instance name":     "instance-1",
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			tags, err := test.tagManager.GenerateTags(
				test.action,
				test.environment,
				test.serviceID,
				test.servicePlanID,
				test.organizationGUID,
				test.spaceGUID,
				test.instanceGUID,
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
		})
	}
}

func TestGenerateTagsHandleErrors(t *testing.T) {
	testCases := map[string]struct {
		tagManager  *CfTagManager
		expectedErr error
	}{
		"error getting organization name": {
			tagManager: &CfTagManager{
				cfNameResolver: &mockCFClientWrapper{
					getOrganizationErr: errors.New("error getting organization name"),
				},
			},
			expectedErr: errors.New("error getting organization name"),
		},
		"error getting space name": {
			tagManager: &CfTagManager{
				broker: "AWS Broker",
				cfNameResolver: &mockCFClientWrapper{
					getSpaceErr: errors.New("error getting space name"),
				},
			},
			expectedErr: errors.New("error getting space name"),
		},
		"error getting instance name": {
			tagManager: &CfTagManager{
				broker: "AWS Broker",
				cfNameResolver: &mockCFClientWrapper{
					getInstanceErr: errors.New("error getting instance name"),
				},
			},
			expectedErr: errors.New("error getting instance name"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := test.tagManager.GenerateTags(
				Create,
				"testing",
				"abc1",
				"abc2",
				"abc3",
				"abc4",
				"abc5",
			)
			if err == nil || err.Error() != test.expectedErr.Error() {
				t.Fatalf("did not received expected err: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}
