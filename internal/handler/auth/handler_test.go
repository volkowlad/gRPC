package auth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/volkowlad/gRPC/protos/gen"
	"testing"

	mock_auth "github.com/volkowlad/gRPC/internal/handler/auth/mock"

	"github.com/golang/mock/gomock"
)

func TestNewHandlers(t *testing.T) {
	type mockBehaviour func(s *mock_auth.MockService, ctx context.Context, req *gen.RegisterRequest)

	type input struct {
		ctx context.Context
		req *gen.RegisterRequest
	}

	type expected struct {
		resp *gen.RegisterResponse
		err  error
	}

	testCases := []struct {
		name          string
		input         input
		mockBehaviour mockBehaviour
		expected      expected
		expErr        bool
	}{
		{
			name: "success",
			input: input{
				ctx: context.Background(),
				req: &gen.RegisterRequest{
					Username: "test",
					Password: "test",
				},
			},
			mockBehaviour: func(s *mock_auth.MockService, ctx context.Context, req *gen.RegisterRequest) {
				s.EXPECT().Register(ctx, req.GetUsername(), req.GetPassword()).Return("done", nil)
			},
			expected: expected{
				resp: &gen.RegisterResponse{
					Message: "done",
				},
				err: nil,
			},
			expErr: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_auth.NewMockService(c)
			test.mockBehaviour(service, test.input.ctx, test.input.req)

			handlers := &Server{
				service: service,
			}

			res, err := handlers.Register(test.input.ctx, test.input.req)
			if test.expErr {
				assert.ErrorContains(t, err, test.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected.resp, res)
			}
		})
	}
}
