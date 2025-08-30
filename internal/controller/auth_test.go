package controller

import (
	"context"
	"testing"

	"github.com/MukizuL/GophKeeper/internal/errs"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	mockservices "github.com/MukizuL/GophKeeper/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type fakeStream struct {
	grpc.ServerTransportStream
	header metadata.MD
}

func (f *fakeStream) SetHeader(md metadata.MD) error {
	f.header = md
	return nil
}
func (f *fakeStream) SendHeader(md metadata.MD) error { return nil }
func (f *fakeStream) SetTrailer(md metadata.MD) error { return nil }

func TestRegister(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		login        string
		password     string
		mockServices func(m *mockservices.MockServicesI)
		code         codes.Code
	}{
		{
			name:         "too short login",
			login:        "ab",
			password:     "12345678",
			mockServices: nil,
			code:         codes.InvalidArgument,
		},
		{
			name:         "too long login",
			login:        string(make([]rune, 256)),
			password:     "12345678",
			mockServices: nil,
			code:         codes.InvalidArgument},
		{
			name:         "too short password",
			login:        "validlogin",
			password:     "1234567",
			mockServices: nil,
			code:         codes.InvalidArgument,
		},
		{
			name:         "too long password runes",
			login:        "validlogin",
			password:     string(make([]rune, 37)),
			mockServices: nil,
			code:         codes.InvalidArgument,
		},
		{
			name:         "too long password bytes",
			login:        "validlogin",
			password:     string(make([]byte, 73)),
			mockServices: nil,
			code:         codes.InvalidArgument,
		},
		{
			name:     "duplicate login",
			login:    "validlogin",
			password: "validpass",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().CreateNewUser(gomock.Any(), "validlogin", "validpass").
					Return(errs.ErrDuplicateLogin)
			},
			code: codes.FailedPrecondition,
		},
		{
			name:     "internal error",
			login:    "validlogin",
			password: "validpass",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().CreateNewUser(gomock.Any(), "validlogin", "validpass").
					Return(errs.ErrInternalServerError)
			},
			code: codes.Internal,
		},
		{
			name:     "success",
			login:    "validlogin",
			password: "validpass",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().CreateNewUser(gomock.Any(), "validlogin", "validpass").
					Return(nil)
			},
			code: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockServices := mockservices.NewMockServicesI(ctrl)
			if tt.mockServices != nil {
				tt.mockServices(mockServices)
			}

			c := Controller{
				services: mockServices,
			}

			resp, err := c.Register(ctx, &pb.RegisterRequest{
				Login:    tt.login,
				Password: tt.password,
			})

			if tt.code == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			} else {
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.code, st.Code())
			}
		})
	}
}

func TestAuthorize(t *testing.T) {
	tests := []struct {
		name         string
		login        string
		password     string
		mockServices func(m *mockservices.MockServicesI)
		code         codes.Code
	}{
		{
			name:     "wrong credentials",
			login:    "user",
			password: "wrongpass",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().Login(gomock.Any(), "user", "wrongpass").
					Return("", nil, errs.ErrWrongCredentials)
			},
			code: codes.Unauthenticated,
		},
		{
			name:     "signing token error",
			login:    "user",
			password: "validpass",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().Login(gomock.Any(), "user", "validpass").
					Return("", nil, errs.ErrSigningToken)
			},
			code: codes.Internal,
		},
		{
			name:     "internal error",
			login:    "user",
			password: "validpass",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().Login(gomock.Any(), "user", "validpass").
					Return("", nil, errs.ErrInternalServerError)
			},
			code: codes.Internal,
		},
		{
			name:     "success",
			login:    "user",
			password: "validpass",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().Login(gomock.Any(), "user", "validpass").
					Return("token123", []byte("derivedkey"), nil)
			},
			code: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockServices := mockservices.NewMockServicesI(ctrl)
			if tt.mockServices != nil {
				tt.mockServices(mockServices)
			}

			c := Controller{services: mockServices}

			fs := &fakeStream{}
			ctx := grpc.NewContextWithServerTransportStream(context.Background(), fs)

			resp, err := c.Authorize(ctx, &pb.AuthRequest{
				Login:    tt.login,
				Password: tt.password,
			})

			if tt.code == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			} else {
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.code, st.Code())
			}
		})
	}
}
