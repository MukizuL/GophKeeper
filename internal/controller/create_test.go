package controller

import (
	"context"
	"testing"

	"github.com/MukizuL/GophKeeper/internal/errs"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	mockpb "github.com/MukizuL/GophKeeper/internal/proto/mocks"
	mockservices "github.com/MukizuL/GophKeeper/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestCreatePassword(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		data         []byte
		mockServices func(m *mockservices.MockServicesI)
		code         codes.Code
	}{
		{
			name: "internal error",
			data: []byte("some-secret"),
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().CreatePassword(gomock.Any(), []byte("some-secret")).
					Return(errs.ErrInternalServerError)
			},
			code: codes.Internal,
		},
		{
			name: "success",
			data: []byte("some-secret"),
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().CreatePassword(gomock.Any(), []byte("some-secret")).
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

			c := Controller{services: mockServices}

			resp, err := c.CreatePassword(ctx, &pb.CreatePasswordRequest{
				Data: tt.data,
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

func TestCreateBank(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		data         []byte
		mockServices func(m *mockservices.MockServicesI)
		code         codes.Code
	}{
		{
			name: "internal error",
			data: []byte("some-secret"),
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().CreateBank(gomock.Any(), []byte("some-secret")).
					Return(errs.ErrInternalServerError)
			},
			code: codes.Internal,
		},
		{
			name: "success",
			data: []byte("some-secret"),
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().CreateBank(gomock.Any(), []byte("some-secret")).
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

			c := Controller{services: mockServices}

			resp, err := c.CreateBank(ctx, &pb.CreateBankRequest{
				Data: tt.data,
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

func TestCreateText(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		data         []byte
		mockServices func(m *mockservices.MockServicesI)
		code         codes.Code
	}{
		{
			name: "internal error",
			data: []byte("some-secret"),
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().CreateTextual(gomock.Any(), []byte("some-secret")).
					Return(errs.ErrInternalServerError)
			},
			code: codes.Internal,
		},
		{
			name: "success",
			data: []byte("some-secret"),
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().CreateTextual(gomock.Any(), []byte("some-secret")).
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

			c := Controller{services: mockServices}

			resp, err := c.CreateText(ctx, &pb.CreateTextRequest{
				Data: tt.data,
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

func TestCreateData(t *testing.T) {
	type CreateDataServer = mockpb.MockGophkeeper_CreateDataServer[pb.CreateDataRequest, pb.CreateDataResponse]
	tests := []struct {
		name         string
		mockStream   func(m *CreateDataServer)
		mockServices func(m *mockservices.MockServicesI)
		code         codes.Code
	}{
		{
			name: "unauthenticated - no token in ctx",
			mockStream: func(m *CreateDataServer) {
				m.EXPECT().Context().Return(context.Background()).AnyTimes()
			},
			mockServices: nil,
			code:         codes.Unauthenticated,
		},
		{
			name: "not authorized",
			mockStream: func(m *CreateDataServer) {
				md := metadata.Pairs("access-token", "valid-token")
				ctx := metadata.NewIncomingContext(context.Background(), md)
				m.EXPECT().Context().Return(ctx).AnyTimes()
			},
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().CreateData(gomock.Any(), "valid-token", gomock.Any()).
					Return(errs.ErrNotAuthorized)
			},
			code: codes.Unauthenticated,
		},
		{
			name: "internal error",
			mockStream: func(m *CreateDataServer) {
				md := metadata.Pairs("access-token", "valid-token")
				ctx := metadata.NewIncomingContext(context.Background(), md)
				m.EXPECT().Context().Return(ctx).AnyTimes()
			},
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().CreateData(gomock.Any(), "valid-token", gomock.Any()).
					Return(errs.ErrInternalServerError)
			},
			code: codes.Internal,
		},
		{
			name: "success",
			mockStream: func(m *CreateDataServer) {
				md := metadata.Pairs("access-token", "valid-token")
				ctx := metadata.NewIncomingContext(context.Background(), md)
				m.EXPECT().Context().Return(ctx).AnyTimes()
			},
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().CreateData(gomock.Any(), "valid-token", gomock.Any()).
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

			stream := mockpb.NewMockGophkeeper_CreateDataServer[pb.CreateDataRequest, pb.CreateDataResponse](ctrl)
			if tt.mockStream != nil {
				tt.mockStream(stream)
			}

			c := Controller{services: mockServices}

			err := c.CreateData(stream)

			if tt.code == codes.OK {
				assert.NoError(t, err)
			} else {
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.code, st.Code())
			}
		})
	}
}
