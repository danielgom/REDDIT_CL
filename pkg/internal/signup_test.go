package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Parallel()
	tCase := []struct {
		name     string
		username string
		email    string
		password string
		err      error
	}{
		{"empty username fails", "", "test_1@hotmail.com", "Password1@", errEmptyUsername},
		{"invalid mail fails", "Daniel", "test_1@hotmail@", "Password1@", errInvalidEmail},
		{"invalid password fails", "Daniel", "test_1@hotmail.com", "Passwo", errInvalidPassword},
		{"all fields valid success", "Daniel", "test_1@hotmail.com", "Password1@", nil},
	}

	for _, test := range tCase {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			req := RegisterRequest{
				Username: test.username,
				Password: test.password,
				Email:    test.email,
			}
			err := req.Validate()
			require.Equal(t, err, test.err)
		})
	}
}
