package auth

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	testCases := []struct {
		Password      string
		ExpectedError error
	}{
		{
			Password:      "",
			ExpectedError: errors.New("password can't be empty"),
		},
		{
			Password:      "somepassword",
			ExpectedError: nil,
		},
		{
			Password:      strings.Repeat("a", 10001),
			ExpectedError: errors.New("password length can't exceed 10000 characters"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Password, func(t *testing.T) {
			_, err := NewPassword(testCase.Password)
			if testCase.ExpectedError == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, testCase.ExpectedError.Error())
			}
		})
	}

}
