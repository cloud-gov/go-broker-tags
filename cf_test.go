package brokertags

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/resource"
	"github.com/google/go-cmp/cmp"
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

type mockServiceInstances struct {
	getServiceInstanceErr error
	instanceName          string
	instanceGUID          string
}

func (si *mockServiceInstances) Get(ctx context.Context, guid string) (*resource.ServiceInstance, error) {
	if si.getServiceInstanceErr != nil {
		return nil, si.getServiceInstanceErr
	}
	if guid != si.instanceGUID {
		return nil, fmt.Errorf("guid argument: %s does not match expected guid: %s", guid, si.instanceGUID)
	}
	return &resource.ServiceInstance{
		Name: si.instanceName,
	}, nil
}

func TestGetOrganization(t *testing.T) {
	testCases := map[string]struct {
		cfResourceGetter     *cfResourceGetter
		expectedOrganization *resource.Organization
		expectedErr          error
		organizationGuid     string
	}{
		"success": {
			cfResourceGetter: &cfResourceGetter{
				Organizations: &mockOrganizations{
					organizationName: "org-1",
					organizationGuid: "guid-1",
				},
			},
			organizationGuid: "guid-1",
			expectedOrganization: &resource.Organization{
				Name: "org-1",
			},
		},
		"error": {
			cfResourceGetter: &cfResourceGetter{
				Organizations: &mockOrganizations{
					getOrganizationErr: errors.New("error getting organization"),
				},
			},
			expectedErr: errors.New("error getting organization"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			organization, err := test.cfResourceGetter.getOrganization(test.organizationGuid)
			if !cmp.Equal(organization, test.expectedOrganization) {
				t.Errorf(cmp.Diff(organization, test.expectedOrganization))
			}
			if (test.expectedErr != nil && err == nil) ||
				(err != nil && err.Error() != test.expectedErr.Error()) {
				t.Fatalf("expected error: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}

func TestGetSpace(t *testing.T) {
	testCases := map[string]struct {
		cfResourceGetter *cfResourceGetter
		expectedSpace    *resource.Space
		expectedErr      error
		spaceGuid        string
	}{
		"success": {
			cfResourceGetter: &cfResourceGetter{
				Spaces: &mockSpaces{
					spaceName: "space-1",
					spaceGuid: "guid-1",
				},
			},
			spaceGuid: "guid-1",
			expectedSpace: &resource.Space{
				Name: "space-1",
			},
		},
		"error": {
			cfResourceGetter: &cfResourceGetter{
				Spaces: &mockSpaces{
					getSpaceErr: errors.New("error getting space"),
				},
			},
			expectedErr: errors.New("error getting space"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			space, err := test.cfResourceGetter.getSpace(test.spaceGuid)
			if !cmp.Equal(space, test.expectedSpace) {
				t.Errorf(cmp.Diff(space, test.expectedSpace))
			}
			if (test.expectedErr != nil && err == nil) ||
				(err != nil && err.Error() != test.expectedErr.Error()) {
				t.Fatalf("expected error: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}

func TestGetServiceInstance(t *testing.T) {
	testCases := map[string]struct {
		cfResourceGetter        *cfResourceGetter
		expectedServiceInstance *resource.ServiceInstance
		expectedErr             error
		instanceGUID            string
	}{
		"success": {
			cfResourceGetter: &cfResourceGetter{
				ServiceInstances: &mockServiceInstances{
					instanceName: "instance-1",
					instanceGUID: "guid-1",
				},
			},
			instanceGUID: "guid-1",
			expectedServiceInstance: &resource.ServiceInstance{
				Name: "instance-1",
			},
		},
		"error": {
			cfResourceGetter: &cfResourceGetter{
				ServiceInstances: &mockServiceInstances{
					getServiceInstanceErr: errors.New("error getting instance"),
				},
			},
			expectedErr: errors.New("error getting instance"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			instance, err := test.cfResourceGetter.getServiceInstance(test.instanceGUID)
			if !cmp.Equal(instance, test.expectedServiceInstance) {
				t.Errorf(cmp.Diff(instance, test.expectedServiceInstance))
			}
			if (test.expectedErr != nil && err == nil) ||
				(err != nil && err.Error() != test.expectedErr.Error()) {
				t.Fatalf("expected error: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}
