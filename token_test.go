package openstack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsAdmin(t *testing.T) {
	tests := []struct {
		description    string
		testData       []map[string]string
		expectedResult bool
	}{
		{
			description: "testData contains admin role, expected to be admin",
			testData: []map[string]string{
				map[string]string{
					"name": "_member_",
					"id":   "fake_id1",
				},
				map[string]string{
					"name": "admin",
					"id":   "fake_id2",
				},
			},
			expectedResult: true,
		},
		{
			description: "testData does not contain admin role, not expected to be admin",
			testData: []map[string]string{
				map[string]string{
					"name": "_member_",
					"id":   "fake_id1",
				},
				map[string]string{
					"name": "_member_",
					"id":   "fake_id2",
				},
			},
			expectedResult: false,
		},
	}

	for _, testCase := range tests {
		tokenInfo := TokenInfo{Roles: testCase.testData}
		assert.Equal(t, testCase.expectedResult,
			tokenInfo.IsAdmin(), "isAdmin is not working properly")
	}
}
