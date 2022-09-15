package internal

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidate(t *testing.T) {
	tCase := []struct {
		name     string
		username string
		email    string
		password string
		err      error
	}{
		{"empty username fails", "", "test_1@hotmail.com", "Password1@", errors.New("username should not be empty")},
		{"invalid mail fails", "Daniel", "test_1@hotmail@", "Password1@", errors.New("invalid email")},
		{"invalid password fails", "Daniel", "test_1@hotmail.com", "Passwo", errors.New("invalid password")},
		{"all fields valid success", "Daniel", "test_1@hotmail.com", "Password1@", nil},
	}

	for _, test := range tCase {
		t.Run(test.name, func(t *testing.T) {
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
