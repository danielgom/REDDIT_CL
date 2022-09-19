package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildLoginResponse(t *testing.T) {
	t.Parallel()

	want := &LoginResponse{
		Username: "daniel",
		Email:    "dga_355@hotmail.com",
		Token:    "asdsdasd",
	}

	response := BuildLoginResponse("daniel", "dga_355@hotmail.com", "asdsdasd")

	require.Equal(t, want, response)
}
