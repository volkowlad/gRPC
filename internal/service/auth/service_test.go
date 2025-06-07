package auth

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"

	mock_auth "github.com/volkowlad/gRPC/internal/handler/auth/mock"
)

func TestRegister(t *testing.T) {
	type mockBehaviour func(s *mock_auth.MockService, ctx context.Context, username, password string)

	type input struct {
		ctx      context.Context
		username string
		password string
	}

	type expected struct {
		message string
		err     error
	}

	testsCase := []struct {
		name          string
		input         input
		mockBehaviour mockBehaviour
		expected      expected
		expErr        bool
	}{
		{
			name: "success",
			input: input{
				ctx:      context.Background(),
				username: "test",
				password: "test",
			},
			mockBehaviour: func(s *mock_auth.MockService, ctx context.Context, username, password string) {
				s.EXPECT().Register(ctx, username, password).Return("done", nil)
			},
			expected: expected{
				message: "done",
				err:     nil,
			},
			expErr: false,
		},
		{
			name: "already exists",
			input: input{
				ctx:      context.Background(),
				username: "test",
				password: "test",
			},
			mockBehaviour: func(s *mock_auth.MockService, ctx context.Context, username, password string) {
				s.EXPECT().Register(ctx, username, password).Return("", errors.New("already exists"))
			},
			expected: expected{
				message: "",
				err:     errors.New("already exists"),
			},
			expErr: true,
		},
		{
			name: "failed to register",
			input: input{
				ctx:      context.Background(),
				username: "test",
				password: "test",
			},
			mockBehaviour: func(s *mock_auth.MockService, ctx context.Context, username, password string) {
				s.EXPECT().Register(ctx, username, password).Return("", errors.New("failed to register user"))
			},
			expected: expected{
				message: "",
				err:     errors.New("failed to register user"),
			},
			expErr: true,
		},
	}

	for _, test := range testsCase {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_auth.NewMockService(c)
			test.mockBehaviour(service, test.input.ctx, test.input.username, test.input.password)

			res, err := service.Register(test.input.ctx, test.input.username, test.input.password)
			if test.expErr {
				assert.ErrorContains(t, err, test.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected.message, res)
			}
		})
	}
}
