package pattern_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wjiec/alchemy/cmd/protoc-gen-alchemy/internal/gengo/pattern"
)

func TestParse(t *testing.T) {
	cases := []struct {
		Pattern string
		Names   []string
		Error   bool
	}{
		{
			Pattern: "/api/users/{id}",
			Names:   []string{"id"},
		},
		{
			Pattern: "/api/users/{id:[0-5]+}",
			Names:   []string{"id"},
		},
		{
			Pattern: "/api/users/{id:[1-5]{8,}}",
			Names:   []string{"id"},
		},
		{
			Pattern: "/api/users/{id}/names/{name}",
			Names:   []string{"id", "name"},
		},
		{
			Pattern: "/api/users",
			Names:   []string{},
		},
		{
			Pattern: "/api/users/{id",
			Error:   true,
		},
		{
			Pattern: "/api/users/{id:{}",
			Error:   true,
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			names, err := pattern.Parse(tt.Pattern)
			if (tt.Error && assert.Error(t, err)) || assert.NoError(t, err) {
				assert.Equal(t, tt.Names, names)
			}
		})
	}
}
