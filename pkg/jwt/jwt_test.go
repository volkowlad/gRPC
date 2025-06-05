package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/volkowlad/gRPC/internal/config"
	"github.com/volkowlad/gRPC/internal/domain"
)

func TestNewAccessToken(t *testing.T) {
	type args struct {
		cfg  config.Token
		user domain.Users
	}

	type result struct {
		token string
		err   string
	}

	testCases := []struct {
		name   string
		args   args
		result result
		expErr bool
	}{
		{
			name: "success",
			args: args{
				cfg: config.Token{
					JWTSecret: "secret",
					AccessTTL: time.Hour,
				},
				user: domain.Users{
					ID:       uuid.UUID{},
					Username: "username",
				},
			},
			result: result{
				token: "token",
				err:   "",
			},
			expErr: false,
		},
		{
			name: "error",
			args: args{
				cfg: config.Token{
					JWTSecret: "",
					AccessTTL: time.Hour,
				},
				user: domain.Users{
					ID:       uuid.UUID{},
					Username: "username",
				},
			},
			result: result{
				token: "",
				err:   "jwt secret is required",
			},
			expErr: true,
		},
		{
			name: "error",
			args: args{
				cfg: config.Token{
					JWTSecret: "secret",
					AccessTTL: time.Minute,
				},
				user: domain.Users{
					ID:       uuid.UUID{},
					Username: "username",
				},
			},
			result: result{
				token: "",
				err:   "jwt token ttl is less than 3600",
			},
			expErr: true,
		},
	}

	for _, cases := range testCases {
		t.Run(cases.name, func(t *testing.T) {
			t.Parallel()

			token, err := NewAccessToken(cases.args.cfg, cases.args.user)
			if cases.expErr {
				assert.EqualError(t, err, cases.result.err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token, cases.result.token)
			}
		})
	}
}
