package services

import (
	"context"
	"errors"
	"testing"

	"github.com/MukizuL/GophKeeper/internal/ctxutil"
	"github.com/MukizuL/GophKeeper/internal/dto"
	"github.com/MukizuL/GophKeeper/internal/errs"
	mockjwt "github.com/MukizuL/GophKeeper/internal/jwt/mocks"
	"github.com/MukizuL/GophKeeper/internal/models"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	mockpb "github.com/MukizuL/GophKeeper/internal/proto/mocks"
	mockstorage "github.com/MukizuL/GophKeeper/internal/storage/mocks"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateNewUser(t *testing.T) {
	tests := []struct {
		name        string
		login       string
		password    string
		mockStorage func(m *mockstorage.MockRepository)
		wantErr     error
	}{
		{
			name:     "success",
			login:    "user1",
			password: "password",
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					CreateNewUser(gomock.Any(), "user1", gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:     "duplicate login",
			login:    "user2",
			password: "password",
			mockStorage: func(m *mockstorage.MockRepository) {
				pgErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation}
				m.EXPECT().
					CreateNewUser(gomock.Any(), "user2", gomock.Any(), gomock.Any()).
					Return(pgErr)
			},
			wantErr: errs.ErrDuplicateLogin,
		},
		{
			name:     "other pg error",
			login:    "user3",
			password: "password",
			mockStorage: func(m *mockstorage.MockRepository) {
				pgErr := &pgconn.PgError{Code: "some-other-wantErr"}
				m.EXPECT().
					CreateNewUser(gomock.Any(), "user3", gomock.Any(), gomock.Any()).
					Return(pgErr)
			},
			wantErr: errs.ErrInternalServerError,
		},
		{
			name:     "generic error",
			login:    "user4",
			password: "password",
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					CreateNewUser(gomock.Any(), "user4", gomock.Any(), gomock.Any()).
					Return(errors.New("db down"))
			},
			wantErr: errs.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockstorage.NewMockRepository(ctrl)
			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}

			s := Services{
				storage: mockRepo,
				logger:  zap.NewNop(), // use a no-op logger for tests
			}

			err := s.CreateNewUser(context.Background(), tt.login, tt.password)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name        string
		login       string
		password    string
		mockStorage func(m *mockstorage.MockRepository)
		mockJWT     func(m *mockjwt.MockServiceI)
		wantErr     error
		wantToken   string
		wantSalt    []byte
	}{
		{
			name:     "success",
			login:    "user1",
			password: "password",
			mockStorage: func(m *mockstorage.MockRepository) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				m.EXPECT().
					GetUserByLogin(gomock.Any(), "user1").
					Return(&models.User{
						ID:       "123",
						Password: hash,
						Salt:     []byte("salty"),
					}, nil)
			},
			mockJWT: func(m *mockjwt.MockServiceI) {
				m.EXPECT().
					CreateToken("123").
					Return("valid-token", nil)
			},
			wantErr:   nil,
			wantToken: "valid-token",
			wantSalt:  []byte("salty"),
		},
		{
			name:     "wrong login",
			login:    "nouser",
			password: "password",
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					GetUserByLogin(gomock.Any(), "nouser").
					Return(&models.User{}, pgx.ErrNoRows)
			},
			mockJWT:   nil,
			wantErr:   errs.ErrWrongCredentials,
			wantToken: "",
			wantSalt:  nil,
		},
		{
			name:     "db error",
			login:    "dberror",
			password: "password",
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					GetUserByLogin(gomock.Any(), "dberror").
					Return(&models.User{}, errors.New("db down"))
			},
			mockJWT:   nil,
			wantErr:   errs.ErrInternalServerError,
			wantToken: "",
			wantSalt:  nil,
		},
		{
			name:     "wrong password",
			login:    "user2",
			password: "wrongpass",
			mockStorage: func(m *mockstorage.MockRepository) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
				m.EXPECT().
					GetUserByLogin(gomock.Any(), "user2").
					Return(&models.User{
						ID:       "234",
						Password: hash,
						Salt:     []byte("salt2"),
					}, nil)
			},
			mockJWT:   nil,
			wantErr:   errs.ErrWrongCredentials,
			wantToken: "",
			wantSalt:  nil,
		},
		{
			name:     "jwt creation error",
			login:    "user3",
			password: "password",
			mockStorage: func(m *mockstorage.MockRepository) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				m.EXPECT().
					GetUserByLogin(gomock.Any(), "user3").
					Return(&models.User{
						ID:       "345",
						Password: hash,
						Salt:     []byte("salt3"),
					}, nil)
			},
			mockJWT: func(m *mockjwt.MockServiceI) {
				m.EXPECT().
					CreateToken("345").
					Return("", errs.ErrSigningToken)
			},
			wantErr:   errs.ErrSigningToken,
			wantToken: "",
			wantSalt:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockstorage.NewMockRepository(ctrl)
			mockJWT := mockjwt.NewMockServiceI(ctrl)

			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}
			if tt.mockJWT != nil {
				tt.mockJWT(mockJWT)
			}

			s := Services{
				storage:    mockRepo,
				jwtService: mockJWT,
				logger:     zap.NewNop(),
			}

			token, salt, err := s.Login(context.Background(), tt.login, tt.password)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
				assert.Equal(t, tt.wantSalt, salt)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, "", token)
				assert.Nil(t, salt)
			}
		})
	}
}

func TestCreatePassword(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		data        []byte
		mockStorage func(m *mockstorage.MockRepository)
		wantErr     error
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-123"),
			data: []byte("secret-data"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					CreatePassword(gomock.Any(), "user-123", []byte("secret-data")).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:        "invalid context userID missing",
			ctx:         context.Background(),
			data:        []byte("irrelevant"),
			mockStorage: nil, // storage never called
			wantErr:     errs.ErrInternalServerError,
		},
		{
			name: "invalid context wrong type",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, 12345),
			data: []byte("irrelevant"),
			// storage never called
			wantErr: errs.ErrInternalServerError,
		},
		{
			name: "storage error",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-456"),
			data: []byte("some-data"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					CreatePassword(gomock.Any(), "user-456", []byte("some-data")).
					Return(errors.New("db fail"))
			},
			wantErr: errs.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockstorage.NewMockRepository(ctrl)
			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}

			s := Services{
				storage: mockRepo,
				logger:  zap.NewNop(),
			}

			err := s.CreatePassword(tt.ctx, tt.data)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestCreateBank(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		data        []byte
		mockStorage func(m *mockstorage.MockRepository)
		wantErr     error
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-123"),
			data: []byte("secret-data"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					CreateBank(gomock.Any(), "user-123", []byte("secret-data")).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:        "invalid context userID missing",
			ctx:         context.Background(),
			data:        []byte("irrelevant"),
			mockStorage: nil, // storage never called
			wantErr:     errs.ErrInternalServerError,
		},
		{
			name: "invalid context wrong type",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, 12345),
			data: []byte("irrelevant"),
			// storage never called
			wantErr: errs.ErrInternalServerError,
		},
		{
			name: "storage error",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-456"),
			data: []byte("some-data"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					CreateBank(gomock.Any(), "user-456", []byte("some-data")).
					Return(errors.New("db fail"))
			},
			wantErr: errs.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockstorage.NewMockRepository(ctrl)
			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}

			s := Services{
				storage: mockRepo,
				logger:  zap.NewNop(),
			}

			err := s.CreateBank(tt.ctx, tt.data)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestCreateText(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		data        []byte
		mockStorage func(m *mockstorage.MockRepository)
		wantErr     error
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-123"),
			data: []byte("secret-data"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					CreateTextual(gomock.Any(), "user-123", []byte("secret-data")).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:        "invalid context userID missing",
			ctx:         context.Background(),
			data:        []byte("irrelevant"),
			mockStorage: nil, // storage never called
			wantErr:     errs.ErrInternalServerError,
		},
		{
			name: "invalid context wrong type",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, 12345),
			data: []byte("irrelevant"),
			// storage never called
			wantErr: errs.ErrInternalServerError,
		},
		{
			name: "storage error",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-456"),
			data: []byte("some-data"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					CreateTextual(gomock.Any(), "user-456", []byte("some-data")).
					Return(errors.New("db fail"))
			},
			wantErr: errs.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockstorage.NewMockRepository(ctrl)
			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}

			s := Services{
				storage: mockRepo,
				logger:  zap.NewNop(),
			}

			err := s.CreateTextual(tt.ctx, tt.data)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestCreateData(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		mockJWT     func(m *mockjwt.MockServiceI)
		mockStorage func(m *mockstorage.MockRepository)
		wantErr     error
	}{
		{
			name:  "success",
			token: "valid-token",
			mockJWT: func(m *mockjwt.MockServiceI) {
				m.EXPECT().
					ValidateToken("valid-token").
					Return("user-123", nil)
			},
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					CreateData(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]byte("encrypted-filename"), nil)

				m.EXPECT().
					CreateReference(gomock.Any(), "user-123", gomock.Any(), []byte("encrypted-filename")).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:  "jwt validation fails",
			token: "invalid-token",
			mockJWT: func(m *mockjwt.MockServiceI) {
				m.EXPECT().
					ValidateToken("invalid-token").
					Return("", errors.New("invalid token"))
			},
			mockStorage: nil, // storage never called
			wantErr:     errors.New("invalid token"),
		},
		{
			name:  "create data fails",
			token: "valid-token",
			mockJWT: func(m *mockjwt.MockServiceI) {
				m.EXPECT().
					ValidateToken("valid-token").
					Return("user-123", nil)
			},
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					CreateData(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("disk error"))
			},
			wantErr: errs.ErrInternalServerError,
		},
		{
			name:  "create reference fails",
			token: "valid-token",
			mockJWT: func(m *mockjwt.MockServiceI) {
				m.EXPECT().
					ValidateToken("valid-token").
					Return("user-123", nil)
			},
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					CreateData(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]byte("encrypted-filename"), nil)

				m.EXPECT().
					CreateReference(gomock.Any(), "user-123", gomock.Any(), []byte("encrypted-filename")).
					Return(errors.New("db error"))
			},
			wantErr: errs.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockstorage.NewMockRepository(ctrl)
			mockJWT := mockjwt.NewMockServiceI(ctrl)
			mockStream := mockpb.NewMockGophkeeper_CreateDataServer[pb.CreateDataRequest, pb.CreateDataResponse](ctrl)

			if tt.mockJWT != nil {
				tt.mockJWT(mockJWT)
			}
			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}

			s := Services{
				storage:    mockRepo,
				jwtService: mockJWT,
				logger:     zap.NewNop(),
			}

			err := s.CreateData(context.Background(), tt.token, mockStream)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			}
		})
	}
}

func TestGetPasswords(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		mockStorage func(m *mockstorage.MockRepository)
		wantData    [][]byte
		wantErr     error
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-123"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					GetPasswordsByUserID(gomock.Any(), "user-123").
					Return([][]byte{[]byte("secret1"), []byte("secret2")}, nil)
			},
			wantData: [][]byte{[]byte("secret1"), []byte("secret2")},
			wantErr:  nil,
		},
		{
			name:        "userID wrong type",
			ctx:         context.WithValue(context.Background(), ctxutil.UserIDContextKey, 123),
			mockStorage: nil,
			wantData:    nil,
			wantErr:     errs.ErrInternalServerError,
		},
		{
			name: "storage error",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-123"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					GetPasswordsByUserID(gomock.Any(), "user-123").
					Return(nil, errors.New("db down"))
			},
			wantData: nil,
			wantErr:  errs.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockstorage.NewMockRepository(ctrl)

			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}

			s := Services{
				storage: mockRepo,
				logger:  zap.NewNop(),
			}

			got, err := s.GetPasswords(tt.ctx)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantData, got)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			}
		})
	}
}

func TestGetBank(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		mockStorage func(m *mockstorage.MockRepository)
		wantData    [][]byte
		wantErr     error
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-123"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					GetBankByUserID(gomock.Any(), "user-123").
					Return([][]byte{[]byte("secret1"), []byte("secret2")}, nil)
			},
			wantData: [][]byte{[]byte("secret1"), []byte("secret2")},
			wantErr:  nil,
		},
		{
			name:        "userID wrong type",
			ctx:         context.WithValue(context.Background(), ctxutil.UserIDContextKey, 123),
			mockStorage: nil,
			wantData:    nil,
			wantErr:     errs.ErrInternalServerError,
		},
		{
			name: "storage error",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-123"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					GetBankByUserID(gomock.Any(), "user-123").
					Return(nil, errors.New("db down"))
			},
			wantData: nil,
			wantErr:  errs.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockstorage.NewMockRepository(ctrl)

			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}

			s := Services{
				storage: mockRepo,
				logger:  zap.NewNop(),
			}

			got, err := s.GetBank(tt.ctx)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantData, got)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			}
		})
	}
}

func TestGetText(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		mockStorage func(m *mockstorage.MockRepository)
		wantData    [][]byte
		wantErr     error
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-123"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					GetTextualByUserID(gomock.Any(), "user-123").
					Return([][]byte{[]byte("secret1"), []byte("secret2")}, nil)
			},
			wantData: [][]byte{[]byte("secret1"), []byte("secret2")},
			wantErr:  nil,
		},
		{
			name:        "userID wrong type",
			ctx:         context.WithValue(context.Background(), ctxutil.UserIDContextKey, 123),
			mockStorage: nil,
			wantData:    nil,
			wantErr:     errs.ErrInternalServerError,
		},
		{
			name: "storage error",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-123"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					GetTextualByUserID(gomock.Any(), "user-123").
					Return(nil, errors.New("db down"))
			},
			wantData: nil,
			wantErr:  errs.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockstorage.NewMockRepository(ctrl)

			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}

			s := Services{
				storage: mockRepo,
				logger:  zap.NewNop(),
			}

			got, err := s.GetText(tt.ctx)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantData, got)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			}
		})
	}
}

func TestGetData(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		mockStorage func(m *mockstorage.MockRepository)
		wantData    []*pb.File
		wantErr     error
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-123"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					GetReferenceByUserID(gomock.Any(), "user-123").
					Return([]dto.FileReference{
						{ID: "1", Filename: []byte("filename1")},
						{ID: "2", Filename: []byte("filename2")},
					}, nil)
			},
			wantData: []*pb.File{
				{Id: "1", Filename: []byte("filename1")},
				{Id: "2", Filename: []byte("filename2")},
			},
			wantErr: nil,
		},
		{
			name:        "userID wrong type",
			ctx:         context.WithValue(context.Background(), ctxutil.UserIDContextKey, 123),
			mockStorage: nil,
			wantData:    nil,
			wantErr:     errs.ErrInternalServerError,
		},
		{
			name: "storage error",
			ctx:  context.WithValue(context.Background(), ctxutil.UserIDContextKey, "user-123"),
			mockStorage: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					GetReferenceByUserID(gomock.Any(), "user-123").
					Return(nil, errors.New("db down"))
			},
			wantData: nil,
			wantErr:  errs.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockstorage.NewMockRepository(ctrl)

			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}

			s := Services{
				storage: mockRepo,
				logger:  zap.NewNop(),
			}

			got, err := s.GetData(tt.ctx)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantData, got)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			}
		})
	}
}

func TestServices_Download(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		mockJWT  func(m *mockjwt.MockServiceI)
		mockRepo func(m *mockstorage.MockRepository)
		wantErr  error
	}{
		{
			name:  "invalid token",
			token: "bad-token",
			mockJWT: func(m *mockjwt.MockServiceI) {
				m.EXPECT().
					ValidateToken("bad-token").
					Return("", errs.ErrNotAuthorized)
			},
			mockRepo: nil,
			wantErr:  errs.ErrNotAuthorized,
		},
		{
			name:  "storage error",
			token: "valid-token",
			mockJWT: func(m *mockjwt.MockServiceI) {
				m.EXPECT().
					ValidateToken("valid-token").
					Return("user-123", nil)
			},
			mockRepo: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					Download(gomock.Any(), "file-123", gomock.Any()).
					Return(errors.New("disk error"))
			},
			wantErr: errs.ErrInternalServerError,
		},
		{
			name:  "success",
			token: "valid-token",
			mockJWT: func(m *mockjwt.MockServiceI) {
				m.EXPECT().
					ValidateToken("valid-token").
					Return("user-123", nil)
			},
			mockRepo: func(m *mockstorage.MockRepository) {
				m.EXPECT().
					Download(gomock.Any(), "file-123", gomock.Any()).
					Return(nil)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockJWT := mockjwt.NewMockServiceI(ctrl)
			mockRepo := mockstorage.NewMockRepository(ctrl)
			mockStream := mockpb.NewMockGophkeeper_DownloadServer[pb.DownloadResponse](ctrl)

			if tt.mockJWT != nil {
				tt.mockJWT(mockJWT)
			}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			s := Services{
				storage:    mockRepo,
				jwtService: mockJWT,
				logger:     zap.NewNop(),
			}

			err := s.Download(context.Background(), tt.token, "file-123", mockStream)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
