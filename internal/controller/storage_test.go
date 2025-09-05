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

func TestGetPasswords(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		mockServices func(m *mockservices.MockServicesI)
		expectedData [][]byte
		code         codes.Code
	}{
		{
			name: "internal error",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().GetPasswords(gomock.Any()).
					Return(nil, errs.ErrInternalServerError)
			},
			expectedData: nil,
			code:         codes.Internal,
		},
		{
			name: "success",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().GetPasswords(gomock.Any()).
					Return([][]byte{[]byte("pass1"), []byte("pass2")}, nil)
			},
			expectedData: [][]byte{[]byte("pass1"), []byte("pass2")},
			code:         codes.OK,
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

			resp, err := c.GetPasswords(ctx, &pb.GetPasswordsRequest{})

			if tt.code == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedData, resp.Data)
			} else {
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.code, st.Code())
			}
		})
	}
}

func TestGetBank(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		mockServices func(m *mockservices.MockServicesI)
		expectedData [][]byte
		code         codes.Code
	}{
		{
			name: "internal error",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().GetBank(gomock.Any()).
					Return(nil, errs.ErrInternalServerError)
			},
			expectedData: nil,
			code:         codes.Internal,
		},
		{
			name: "success",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().GetBank(gomock.Any()).
					Return([][]byte{[]byte("card1"), []byte("card2")}, nil)
			},
			expectedData: [][]byte{[]byte("card1"), []byte("card2")},
			code:         codes.OK,
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

			resp, err := c.GetBank(ctx, &pb.GetBankRequest{})

			if tt.code == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedData, resp.Data)
			} else {
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.code, st.Code())
			}
		})
	}
}

func TestGetText(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		mockServices func(m *mockservices.MockServicesI)
		expectedData [][]byte
		code         codes.Code
	}{
		{
			name: "internal error",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().GetText(gomock.Any()).
					Return(nil, errs.ErrInternalServerError)
			},
			expectedData: nil,
			code:         codes.Internal,
		},
		{
			name: "success",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().GetText(gomock.Any()).
					Return([][]byte{[]byte("text1"), []byte("text2")}, nil)
			},
			expectedData: [][]byte{[]byte("text1"), []byte("text2")},
			code:         codes.OK,
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

			resp, err := c.GetText(ctx, &pb.GetTextRequest{})

			if tt.code == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedData, resp.Data)
			} else {
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.code, st.Code())
			}
		})
	}
}

func TestGetData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		mockServices func(m *mockservices.MockServicesI)
		expectedData []*pb.File
		code         codes.Code
	}{
		{
			name: "internal error",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().GetData(gomock.Any()).
					Return(nil, errs.ErrInternalServerError)
			},
			expectedData: nil,
			code:         codes.Internal,
		},
		{
			name: "success",
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().GetData(gomock.Any()).
					Return([]*pb.File{
						&pb.File{Id: "1", Filename: []byte("filename1")},
						&pb.File{Id: "2", Filename: []byte("filename2")},
					}, nil)
			},
			expectedData: []*pb.File{
				&pb.File{Id: "1", Filename: []byte("filename1")},
				&pb.File{Id: "2", Filename: []byte("filename2")},
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

			resp, err := c.GetData(ctx, &pb.GetDataRequest{})

			if tt.code == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedData, resp.File)
			} else {
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.code, st.Code())
			}
		})
	}
}

func TestDownload(t *testing.T) {
	ctxWithToken := func(token string) context.Context {
		return metadata.NewIncomingContext(context.Background(),
			metadata.Pairs("access-token", token))
	}

	type DownloadServer = mockpb.MockGophkeeper_DownloadServer[pb.DownloadResponse]

	tests := []struct {
		name         string
		mockStream   func(m *DownloadServer)
		mockServices func(m *mockservices.MockServicesI)
		code         codes.Code
	}{
		{
			name: "unauthenticated - no token",
			mockStream: func(m *DownloadServer) {
				m.EXPECT().Context().Return(context.Background()).AnyTimes()
			},
			mockServices: nil,
			code:         codes.Unauthenticated,
		},
		{
			name: "not authorized",
			mockStream: func(m *DownloadServer) {
				m.EXPECT().Context().Return(ctxWithToken("valid-token")).AnyTimes()
			},
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().
					Download(gomock.Any(), "valid-token", "file-id", gomock.Any()).
					Return(errs.ErrNotAuthorized)
			},
			code: codes.Unauthenticated,
		},
		{
			name: "internal error",
			mockStream: func(m *DownloadServer) {
				m.EXPECT().Context().Return(ctxWithToken("valid-token")).AnyTimes()
			},
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().
					Download(gomock.Any(), "valid-token", "file-id", gomock.Any()).
					Return(errs.ErrInternalServerError)
			},
			code: codes.Internal,
		},
		{
			name: "success",
			mockStream: func(m *DownloadServer) {
				m.EXPECT().Context().Return(ctxWithToken("valid-token")).AnyTimes()
			},
			mockServices: func(m *mockservices.MockServicesI) {
				m.EXPECT().
					Download(gomock.Any(), "valid-token", "file-id", gomock.Any()).
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

			stream := mockpb.NewMockGophkeeper_DownloadServer[pb.DownloadResponse](ctrl)
			if tt.mockStream != nil {
				tt.mockStream(stream)
			}

			c := Controller{services: mockServices}

			err := c.Download(&pb.DownloadRequest{Id: "file-id"}, stream)

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
