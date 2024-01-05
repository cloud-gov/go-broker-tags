package brokerTags

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type mockOrganizations struct {
	getOrganizationErr error
	organizationName   string
	organizationGuid   string
}

func (o *mockOrganizations) Get(ctx context.Context, guid string) (*resource.Organization, error) {
	if o.getOrganizationErr != nil {
		return nil, o.getOrganizationErr
	}
	if guid != o.organizationGuid {
		return nil, fmt.Errorf("guid argument: %s does not match expected guid: %s", guid, o.organizationGuid)
	}
	return &resource.Organization{
		Name: o.organizationName,
	}, nil
}

type mockServiceInstances struct {
	getServiceInstanceErr error
	serviceInstanceName   string
	serviceInstanceGuid   string
}

func (si *mockServiceInstances) Get(ctx context.Context, guid string) (*resource.ServiceInstance, error) {
	if si.getServiceInstanceErr != nil {
		return nil, si.getServiceInstanceErr
	}
	if guid != si.serviceInstanceGuid {
		return nil, fmt.Errorf("guid argument: %s does not match expected guid: %s", guid, si.serviceInstanceGuid)
	}
	return &resource.ServiceInstance{
		Name: si.serviceInstanceName,
	}, nil
}

type mockServiceOfferings struct {
	getServiceOfferingErr error
	serviceOfferingName   string
	serviceOfferingGuid   string
}

func (so *mockServiceOfferings) Get(ctx context.Context, guid string) (*resource.ServiceOffering, error) {
	if so.getServiceOfferingErr != nil {
		return nil, so.getServiceOfferingErr
	}
	if guid != so.serviceOfferingGuid {
		return nil, fmt.Errorf("guid argument: %s does not match expected guid: %s", guid, so.serviceOfferingGuid)
	}
	return &resource.ServiceOffering{
		Name: so.serviceOfferingName,
	}, nil
}

type mockServicePlans struct {
	getServicePlanErr error
	servicePlanName   string
	servicePlanGuid   string
}

func (sp *mockServicePlans) Get(ctx context.Context, guid string) (*resource.ServicePlan, error) {
	if sp.getServicePlanErr != nil {
		return nil, sp.getServicePlanErr
	}
	if guid != sp.servicePlanGuid {
		return nil, fmt.Errorf("guid argument: %s does not match expected guid: %s", guid, sp.servicePlanGuid)
	}
	return &resource.ServicePlan{
		Name: sp.servicePlanName,
	}, nil
}

type mockSpaces struct {
	getSpaceErr error
	spaceName   string
	spaceGuid   string
}

func (s *mockSpaces) Get(ctx context.Context, guid string) (*resource.Space, error) {
	if s.getSpaceErr != nil {
		return nil, s.getSpaceErr
	}
	if guid != s.spaceGuid {
		return nil, fmt.Errorf("guid argument: %s does not match expected guid: %s", guid, s.spaceGuid)
	}
	return &resource.Space{
		Name: s.spaceName,
	}, nil
}

func TestGetRequiredEnvVars(t *testing.T) {
	testCases := map[string]struct {
		envVars   map[string]string
		expectErr bool
	}{
		"no env vars": {
			expectErr: true,
		},
		"one env var set": {
			expectErr: true,
			envVars: map[string]string{
				"CF_API_URL": "api-1",
			},
		},
		"two env vars set": {
			expectErr: true,
			envVars: map[string]string{
				"CF_API_URL":       "api-1",
				"CF_API_CLIENT_ID": "client-1",
			},
		},
		"all env vars set": {
			expectErr: false,
			envVars: map[string]string{
				"CF_API_URL":           "api-1",
				"CF_API_CLIENT_ID":     "client-1",
				"CF_API_CLIENT_SECRET": "secret",
			},
		},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			for envVar, value := range test.envVars {
				os.Setenv(envVar, value)
			}

			_, err := getRequiredEnvVars()
			if !test.expectErr && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if test.expectErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			for envVar := range test.envVars {
				os.Unsetenv(envVar)
			}
		})
	}
}

func TestGetOrganizationName(t *testing.T) {
	testCases := map[string]struct {
		cfClientWrapper          *cfNameResolver
		expectedOrganizationName string
		expectedErr              error
		organizationGuid         string
	}{
		"success": {
			cfClientWrapper: &cfNameResolver{
				Organizations: &mockOrganizations{
					organizationName: "org-1",
					organizationGuid: "guid-1",
				},
				ServiceInstances: &mockServiceInstances{},
				ServiceOfferings: &mockServiceOfferings{},
				ServicePlans:     &mockServicePlans{},
				Spaces:           &mockSpaces{},
			},
			organizationGuid:         "guid-1",
			expectedOrganizationName: "org-1",
		},
		"error": {
			cfClientWrapper: &cfNameResolver{
				Organizations: &mockOrganizations{},
				ServiceInstances: &mockServiceInstances{
					getServiceInstanceErr: errors.New("error getting organization"),
				},
				ServiceOfferings: &mockServiceOfferings{},
				ServicePlans:     &mockServicePlans{},
				Spaces:           &mockSpaces{},
			},
			expectedErr: errors.New("error getting organization"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			organizationName, err := test.cfClientWrapper.getOrganizationName(test.organizationGuid)
			if organizationName != test.expectedOrganizationName {
				t.Fatalf("expected organization name: %s, got: %s", test.expectedOrganizationName, organizationName)
			}
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Fatalf("expected error: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}

func TestGetServiceInstanceName(t *testing.T) {
	testCases := map[string]struct {
		cfClientWrapper      *cfNameResolver
		expectedInstanceName string
		expectedErr          error
		serviceInstanceGuid  string
	}{
		"success": {
			cfClientWrapper: &cfNameResolver{
				Organizations: &mockOrganizations{},
				ServiceInstances: &mockServiceInstances{
					serviceInstanceName: "instance-1",
					serviceInstanceGuid: "guid-1",
				},
				ServiceOfferings: &mockServiceOfferings{},
				ServicePlans:     &mockServicePlans{},
				Spaces:           &mockSpaces{},
			},
			serviceInstanceGuid:  "guid-1",
			expectedInstanceName: "instance-1",
		},
		"error": {
			cfClientWrapper: &cfNameResolver{
				Organizations: &mockOrganizations{},
				ServiceInstances: &mockServiceInstances{
					getServiceInstanceErr: errors.New("error getting service instance"),
				},
				ServiceOfferings: &mockServiceOfferings{},
				ServicePlans:     &mockServicePlans{},
				Spaces:           &mockSpaces{},
			},
			expectedErr: errors.New("error getting service instance"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			offeringName, err := test.cfClientWrapper.getServiceInstanceName(test.serviceInstanceGuid)
			if offeringName != test.expectedInstanceName {
				t.Fatalf("expected instance name: %s, got: %s", test.expectedInstanceName, offeringName)
			}
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Fatalf("expected error: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}

func TestGetServiceOfferingName(t *testing.T) {
	testCases := map[string]struct {
		cfClientWrapper      *cfNameResolver
		expectedOfferingName string
		expectedErr          error
		serviceOfferingGuid  string
	}{
		"success": {
			cfClientWrapper: &cfNameResolver{
				Organizations:    &mockOrganizations{},
				ServiceInstances: &mockServiceInstances{},
				ServiceOfferings: &mockServiceOfferings{
					serviceOfferingName: "offering-1",
					serviceOfferingGuid: "guid-1",
				},
				ServicePlans: &mockServicePlans{},
				Spaces:       &mockSpaces{},
			},
			serviceOfferingGuid:  "guid-1",
			expectedOfferingName: "offering-1",
		},
		"error": {
			cfClientWrapper: &cfNameResolver{
				Organizations:    &mockOrganizations{},
				ServiceInstances: &mockServiceInstances{},
				ServiceOfferings: &mockServiceOfferings{
					getServiceOfferingErr: errors.New("error getting service offering"),
				},
				ServicePlans: &mockServicePlans{},
				Spaces:       &mockSpaces{},
			},
			expectedErr: errors.New("error getting service offering"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			offeringName, err := test.cfClientWrapper.getServiceOfferingName(test.serviceOfferingGuid)
			if offeringName != test.expectedOfferingName {
				t.Fatalf("expected offering name: %s, got: %s", test.expectedOfferingName, offeringName)
			}
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Fatalf("expected error: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}

func TestGetServicePlanName(t *testing.T) {
	testCases := map[string]struct {
		cfClientWrapper  *cfNameResolver
		expectedPlanName string
		expectedErr      error
		servicePlanGuid  string
	}{
		"success": {
			cfClientWrapper: &cfNameResolver{
				Organizations:    &mockOrganizations{},
				ServiceInstances: &mockServiceInstances{},
				ServiceOfferings: &mockServiceOfferings{},
				ServicePlans: &mockServicePlans{
					servicePlanName: "plan-1",
					servicePlanGuid: "guid-1",
				},
				Spaces: &mockSpaces{},
			},
			servicePlanGuid:  "guid-1",
			expectedPlanName: "plan-1",
		},
		"error": {
			cfClientWrapper: &cfNameResolver{
				Organizations:    &mockOrganizations{},
				ServiceInstances: &mockServiceInstances{},
				ServiceOfferings: &mockServiceOfferings{},
				ServicePlans: &mockServicePlans{
					getServicePlanErr: errors.New("error getting service plan"),
				},
				Spaces: &mockSpaces{},
			},
			expectedErr: errors.New("error getting service plan"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			planName, err := test.cfClientWrapper.getServicePlanName(test.servicePlanGuid)
			if planName != test.expectedPlanName {
				t.Fatalf("expected plan name: %s, got: %s", test.expectedPlanName, planName)
			}
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Fatalf("expected error: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}

func TestGetSpaceName(t *testing.T) {
	testCases := map[string]struct {
		cfClientWrapper   *cfNameResolver
		expectedSpaceName string
		expectedErr       error
		spaceGuid         string
	}{
		"success": {
			cfClientWrapper: &cfNameResolver{
				Organizations:    &mockOrganizations{},
				ServiceInstances: &mockServiceInstances{},
				ServiceOfferings: &mockServiceOfferings{},
				ServicePlans:     &mockServicePlans{},
				Spaces: &mockSpaces{
					spaceName: "plan-1",
					spaceGuid: "guid-1",
				},
			},
			spaceGuid:         "guid-1",
			expectedSpaceName: "plan-1",
		},
		"error": {
			cfClientWrapper: &cfNameResolver{
				Organizations:    &mockOrganizations{},
				ServiceInstances: &mockServiceInstances{},
				ServiceOfferings: &mockServiceOfferings{},
				ServicePlans:     &mockServicePlans{},
				Spaces: &mockSpaces{
					getSpaceErr: errors.New("error getting space"),
				},
			},
			expectedErr: errors.New("error getting space"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			spaceName, err := test.cfClientWrapper.getSpaceName(test.spaceGuid)
			if spaceName != test.expectedSpaceName {
				t.Fatalf("expected space name: %s, got: %s", test.expectedSpaceName, spaceName)
			}
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Fatalf("expected error: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}
