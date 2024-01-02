package brokerTags

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type mockOrganizations struct {
	getOrganizationErr error
	organizationName   string
}

func (o *mockOrganizations) Get(ctx context.Context, guid string) (*resource.Organization, error) {
	if o.getOrganizationErr != nil {
		return nil, o.getOrganizationErr
	}
	return &resource.Organization{
		Name: o.organizationName,
	}, nil
}

type mockServiceOfferings struct {
	getServiceOfferingErr error
	serviceOfferingName   string
	serviceGuid           string
}

func (so *mockServiceOfferings) Get(ctx context.Context, guid string) (*resource.ServiceOffering, error) {
	if so.getServiceOfferingErr != nil {
		return nil, so.getServiceOfferingErr
	}
	if guid != so.serviceGuid {
		return nil, fmt.Errorf("guid argument: %s does not match expected guid: %s", guid, so.serviceGuid)
	}
	return &resource.ServiceOffering{
		Name: so.serviceOfferingName,
	}, nil
}

type mockServicePlans struct {
	getServicePlanErr error
	servicePlanName   string
}

func (sp *mockServicePlans) Get(ctx context.Context, guid string) (*resource.ServicePlan, error) {
	if sp.getServicePlanErr != nil {
		return nil, sp.getServicePlanErr
	}
	return &resource.ServicePlan{
		Name: sp.servicePlanName,
	}, nil
}

type mockServiceInstances struct {
	getServicePlanErr error
	servicePlanName   string
}

func (si *mockServiceInstances) Get(ctx context.Context, guid string) (*resource.ServiceInstance, error) {
	if si.getServicePlanErr != nil {
		return nil, si.getServicePlanErr
	}
	return &resource.ServiceInstance{
		Name: si.servicePlanName,
	}, nil
}

type mockSpaces struct {
	getSpaceErr error
	spaceName   string
}

func (s *mockSpaces) Get(ctx context.Context, guid string) (*resource.Space, error) {
	if s.getSpaceErr != nil {
		return nil, s.getSpaceErr
	}
	return &resource.Space{
		Name: s.spaceName,
	}, nil
}

func TestGetServiceOfferingName(t *testing.T) {
	testCases := map[string]struct {
		cfClientWrapper      *cfClientWrapper
		expectedOfferingName string
		expectedErr          error
		serviceOfferingGuid  string
	}{
		"success": {
			cfClientWrapper: &cfClientWrapper{
				Organizations:    &mockOrganizations{},
				ServiceInstances: &mockServiceInstances{},
				ServiceOfferings: &mockServiceOfferings{
					serviceOfferingName: "offering-1",
					serviceGuid:         "guid-1",
				},
				ServicePlans: &mockServicePlans{},
				Spaces:       &mockSpaces{},
			},
			serviceOfferingGuid:  "guid-1",
			expectedOfferingName: "offering-1",
		},
		"error": {
			cfClientWrapper: &cfClientWrapper{
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
