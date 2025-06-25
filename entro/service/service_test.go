package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFileSecrets(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		typeExpcted string
		wantFound   bool
	}{
		{
			name:        "AccesskeyFound",
			input:       "secretASIAJAZ4HRG3CPA63XEQtest",
			typeExpcted: "access_key",
			wantFound:   true,
		},
		{
			name:        "AccesskeyNotFound",
			input:       "secretASIJAZ4HRG3CPA63XEQtest",
			typeExpcted: "",
			wantFound:   false,
		},
		{
			name:        "SecretkeyFound",
			input:       "secreteDcCc9H6oCkGUSp3Rhmsx8NIfVG8kO2T/3jORxuZYtest",
			typeExpcted: "secret_key",
			wantFound:   true,
		},
		{
			name:        "SeessioKeyFound",
			input:       "SessionFQoDYXdzEPP//////////wEaDPv5GPAhRW8pw6/nsiKsAZu7sZDCXPtEBEurxmvyV1r+nWy1I4VPbdIJV+iDnotwS3PKIyj+yDnOeigMf2yp9y2Dg9D7r51vWUyUQQfceZi9/8Ghy38RcOnWImhNdVP5zl1zh85FHz6ytePo+puHZwfTkuAQHj38gy6VF/14GU17qDcPTfjhbETGqEmh8QX6xfmWlO0ZrTmsAo4ZHav8yzbbl3oYdCLICOjMhOO1oY+B/DiURk3ZLPjaXyoo2Iql2QU=",
			typeExpcted: "session_key",
			wantFound:   true,
		},
	}

	for _, test := range tests {
		res := findFileSecrets(test.input)
		assert.Equal(t, test.wantFound, len(res) > 0, fmt.Sprintf(" test  %s expected %v, got %v, ", test.name, test.wantFound, res))
		if test.wantFound {
			assert.Equal(t, test.typeExpcted, res[0].Type)
		}
	}
}
