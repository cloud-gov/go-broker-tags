package brokertags

import "testing"

func TestCreateActionGetTagKey(t *testing.T) {
	testCases := map[string]struct {
		action         Action
		expectedTagKey string
	}{
		"Create": {
			action:         Create,
			expectedTagKey: "Created at",
		},
		"Update": {
			action:         Update,
			expectedTagKey: "Updated at",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			if test.action.getTagKey() != test.expectedTagKey {
				t.Errorf("expected tag key: %s, got: %s", test.expectedTagKey, test.action.getTagKey())
			}
		})
	}
}
