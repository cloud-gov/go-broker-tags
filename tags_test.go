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

func TestGenerateTags(t *testing.T) {
	testCases := map[string]struct {
		tagManager          *CfTagManager
		expectedTags        map[string]string
		action              Action
		serviceOfferingName string
		servicePlanName     string
		organizationGUID    string
		spaceGUID           string
		instanceGUID        string
	}{
		"Create": {
			action:              Create,
			serviceOfferingName: "abc1",
			servicePlanName:     "abc2",
			organizationGUID:    "abc3",
			spaceGUID:           "abc4",
			instanceGUID:        "abc5",
			tagManager: &CfTagManager{
				broker:      "AWS Broker",
				environment: "testing",
				cfNameResolver: &mockCFClientWrapper{
					organizationName: "org-1",
					spaceName:        "space-1",
					spaceGUID:        "abc4",
					organizationGUID: "abc3",
					instanceGUID:     "abc5",
				},
			},
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
			organizationGUID:    "abc3",
			spaceGUID:           "abc4",
			instanceGUID:        "abc5",
			tagManager: &CfTagManager{
				broker:      "AWS Broker",
				environment: "testing",
				cfNameResolver: &mockCFClientWrapper{
					organizationName: "org-1",
					spaceName:        "space-1",
					spaceGUID:        "abc4",
					organizationGUID: "abc3",
					instanceGUID:     "abc5",
				},
			},
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
			organizationGUID:    "abc3",
			spaceGUID:           "abc4",
			instanceGUID:        "abc5",
			tagManager: &CfTagManager{
				environment: "testing",
				cfNameResolver: &mockCFClientWrapper{
					organizationName: "org-1",
					spaceName:        "space-1",
					spaceGUID:        "abc4",
					organizationGUID: "abc3",
					instanceGUID:     "abc5",
				},
			},
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
			organizationGUID:    "abc3",
			spaceGUID:           "abc4",
			instanceGUID:        "abc5",
			tagManager: &CfTagManager{
				broker: "AWS Broker",
				cfNameResolver: &mockCFClientWrapper{
					organizationName: "org-1",
					spaceName:        "space-1",
					spaceGUID:        "abc4",
					organizationGUID: "abc3",
					instanceGUID:     "abc5",
				},
			},
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
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			tags, err := test.tagManager.GenerateTags(
				test.action,
				test.serviceOfferingName,
				test.servicePlanName,
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
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := test.tagManager.GenerateTags(
				Create,
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
