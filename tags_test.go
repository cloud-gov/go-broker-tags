package sharedBrokerUtils

import (
	"errors"
	"reflect"
	"testing"
)

type mockCFClientWrapper struct {
	getServiceOfferingErr error
	serviceOfferingName   string
	getServicePlanErr     error
	servicePlanName       string
	getOrganizationErr    error
	organizationName      string
	getSpaceErr           error
	spaceName             string
	getInstanceErr        error
	instanceName          string
}

func (m *mockCFClientWrapper) getServiceOfferingName(serviceGUID string) (string, error) {
	if m.getServiceOfferingErr != nil {
		return "", m.getServiceOfferingErr
	}
	return m.serviceOfferingName, nil
}

func (m *mockCFClientWrapper) getServicePlanName(servicePlanGUID string) (string, error) {
	if m.getServicePlanErr != nil {
		return "", m.getServicePlanErr
	}
	return m.servicePlanName, nil
}

func (m *mockCFClientWrapper) getOrganizationName(organizationGUID string) (string, error) {
	if m.getOrganizationErr != nil {
		return "", m.getOrganizationErr
	}
	return m.organizationName, nil
}

func (m *mockCFClientWrapper) getSpaceName(spaceGUID string) (string, error) {
	if m.getSpaceErr != nil {
		return "", m.getSpaceErr
	}
	return m.spaceName, nil
}

func (m *mockCFClientWrapper) getInstanceName(instanceGUID string) (string, error) {
	if m.getInstanceErr != nil {
		return "", m.getInstanceErr
	}
	return m.instanceName, nil
}

func TestGenerateTags(t *testing.T) {
	tagManager := &TagManager{
		cfClient: &mockCFClientWrapper{
			serviceOfferingName: "offering-1",
			servicePlanName:     "plan-1",
			organizationName:    "org-1",
			spaceName:           "space-1",
			instanceName:        "instance-1",
		},
	}
	tags, err := tagManager.GenerateTags(
		"AWS S3 Service Broker",
		"Created",
		"abc1",
		"abc2",
		"abc3",
		"abc4",
		"abc5",
	)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if tags["Created at"] == "" {
		t.Fatal("Expected a value for \"Created at\" tag")
	}
	delete(tags, "Created at")

	expectedTags := map[string]string{
		"Owner":                 "Cloud Foundry",
		"Created by":            "AWS S3 Service Broker",
		"Service GUID":          "abc1",
		"Plan GUID":             "abc2",
		"Organization GUID":     "abc3",
		"Space GUID":            "abc4",
		"Instance GUID":         "abc5",
		"Service offering name": "offering-1",
		"Service plan name":     "plan-1",
		"Organization name":     "org-1",
		"Space name":            "space-1",
		"Instance name":         "instance-1",
	}

	if !reflect.DeepEqual(tags, expectedTags) {
		t.Errorf("expected: %s, got: %s", expectedTags, tags)
	}
}

func TestGenerateTagsHandleErrors(t *testing.T) {
	testCases := map[string]struct {
		tagManager  *TagManager
		expectedErr error
	}{
		"error getting service offering name": {
			tagManager: &TagManager{
				cfClient: &mockCFClientWrapper{
					getServiceOfferingErr: errors.New("error getting service offering name"),
				},
			},
			expectedErr: errors.New("error getting service offering name"),
		},
		"error getting service plan name": {
			tagManager: &TagManager{
				cfClient: &mockCFClientWrapper{
					getServicePlanErr: errors.New("error getting service plan name"),
				},
			},
			expectedErr: errors.New("error getting service plan name"),
		},
		"error getting organization name": {
			tagManager: &TagManager{
				cfClient: &mockCFClientWrapper{
					getOrganizationErr: errors.New("error getting organization name"),
				},
			},
			expectedErr: errors.New("error getting organization name"),
		},
		"error getting space name": {
			tagManager: &TagManager{
				cfClient: &mockCFClientWrapper{
					getSpaceErr: errors.New("error getting space name"),
				},
			},
			expectedErr: errors.New("error getting space name"),
		},
		"error getting instance name": {
			tagManager: &TagManager{
				cfClient: &mockCFClientWrapper{
					getInstanceErr: errors.New("error getting instance name"),
				},
			},
			expectedErr: errors.New("error getting instance name"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := test.tagManager.GenerateTags(
				"AWS S3 Service Broker",
				"Created",
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
