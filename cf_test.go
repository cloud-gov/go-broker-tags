package brokertags

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
				Spaces: &mockSpaces{},
			},
			organizationGuid:         "guid-1",
			expectedOrganizationName: "org-1",
		},
		"error": {
			cfClientWrapper: &cfNameResolver{
				Organizations: &mockOrganizations{},
				Spaces:        &mockSpaces{},
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

func TestGetSpaceName(t *testing.T) {
	testCases := map[string]struct {
		cfClientWrapper   *cfNameResolver
		expectedSpaceName string
		expectedErr       error
		spaceGuid         string
	}{
		"success": {
			cfClientWrapper: &cfNameResolver{
				Organizations: &mockOrganizations{},
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
				Organizations: &mockOrganizations{},
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
