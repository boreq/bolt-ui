package auth

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatedUsername(t *testing.T) {
	testCases := []struct {
		Username      string
		ExpectedError error
	}{
		{
			Username:      "",
			ExpectedError: errors.New("username can't be empty"),
		},
		{
			Username:      "Username123-_",
			ExpectedError: nil,
		},
		{
			Username:      "username/",
			ExpectedError: errors.New("username must conform to the following regexp: ^[a-zA-Z0-9_-]+$"),
		},
		{
			Username:      strings.Repeat("a", 101),
			ExpectedError: errors.New("username length can't exceed 100 characters"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Username, func(t *testing.T) {
			_, err := NewValidatedUsername(testCase.Username)
			if testCase.ExpectedError == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, testCase.ExpectedError.Error())
			}
		})
	}

}
