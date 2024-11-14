package secret

import (
	"context"
	"errors"
	"fmt"
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockSecretService struct {
	CreateSecretFunc   func(ctx context.Context, secret model.Secret) error
	FindSecretFunc     func(ctx context.Context, name string) (model.Secret, error)
	FindAllSecretsFunc func(ctx context.Context) ([]model.Secret, error)
	DeleteSecretFunc   func(ctx context.Context, name string) error
}

func (m *mockSecretService) CreateSecret(ctx context.Context, secret model.Secret) error {
	return m.CreateSecretFunc(ctx, secret)
}

func (m *mockSecretService) FindSecret(ctx context.Context, name string) (model.Secret, error) {
	return m.FindSecretFunc(ctx, name)
}

func (m *mockSecretService) FindAllSecrets(ctx context.Context) ([]model.Secret, error) {
	return m.FindAllSecretsFunc(ctx)
}

func (m *mockSecretService) DeleteSecret(ctx context.Context, name string) error {
	return m.DeleteSecretFunc(ctx, name)
}

func TestUploadSecretHandler(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	tests := []struct {
		name           string
		method         string
		body           string
		service        secretService
		expectedStatus int
	}{
		{
			name:   "Successful upload secret",
			method: http.MethodPost,
			body:   `{"name":"testSecret", "content":"dGVzdCBjb250ZW50", "type":"CREDENTIALS", "version":1}`,
			service: &mockSecretService{
				CreateSecretFunc: func(ctx context.Context, secret model.Secret) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid HTTP method",
			method:         http.MethodGet,
			body:           "",
			service:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Empty secret name",
			method: http.MethodPost,
			body:   `{"name":"", "content":"dGVzdCBjb250ZW50", "type":"CREDENTIALS", "version":1}`,
			service: &mockSecretService{
				CreateSecretFunc: func(ctx context.Context, secret model.Secret) error {
					return model.ErrSecretNameIsEmpty
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Secret already exists",
			method: http.MethodPost,
			body:   `{"name":"testSecret", "content":"dGVzdCBjb250ZW50", "type":"CREDENTIALS", "version":1}`,
			service: &mockSecretService{
				CreateSecretFunc: func(ctx context.Context, secret model.Secret) error {
					return model.ErrSecretExistToCurrentUser
				},
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/user/secret", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := UploadSecretHandler(logger, tt.service)
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestReadOneSecretHandler(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	mockSecret := model.Secret{
		Name:    "testSecret",
		Content: []byte("test content"),
		Type:    "password",
		Version: 1,
	}

	tests := []struct {
		name           string
		method         string
		urlParam       string
		service        secretService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "Successful operation",
			method:   http.MethodGet,
			urlParam: "testSecret",
			service: &mockSecretService{
				FindSecretFunc: func(ctx context.Context, name string) (model.Secret, error) {
					return mockSecret, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"name":"testSecret","content":"dGVzdCBjb250ZW50","type":"password"}`,
		},
		{
			name:           "Invalid HTTP method",
			method:         http.MethodPost,
			urlParam:       "testSecret",
			service:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty secret name",
			method:         http.MethodGet,
			urlParam:       "",
			service:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Secret was not found",
			method:   http.MethodGet,
			urlParam: "nonexistentSecret",
			service: &mockSecretService{
				FindSecretFunc: func(ctx context.Context, name string) (model.Secret, error) {
					return model.Secret{}, model.ErrSecretWasNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:     "Internal server error",
			method:   http.MethodGet,
			urlParam: "testSecret",
			service: &mockSecretService{
				FindSecretFunc: func(ctx context.Context, name string) (model.Secret, error) {
					return model.Secret{}, errors.New("database error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/user/secret"+tt.urlParam, nil)
			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("name", tt.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

			rec := httptest.NewRecorder()
			handler := ReadOneSecretHandler(logger, tt.service)
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedStatus == http.StatusOK {
				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, tt.expectedBody, string(body))
			}
		})
	}
}

func TestDeleteSecretHandler(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	tests := []struct {
		name           string
		method         string
		paramName      string
		service        secretService
		expectedStatus int
	}{
		{
			name:      "Successful delete secret",
			method:    http.MethodDelete,
			paramName: "testSecret",
			service: &mockSecretService{
				DeleteSecretFunc: func(ctx context.Context, name string) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Secret was not found",
			method:    http.MethodDelete,
			paramName: "nonexistent",
			service: &mockSecretService{
				DeleteSecretFunc: func(ctx context.Context, name string) error {
					return model.ErrSecretWasNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "Invalid HTTP method",
			method:    http.MethodGet,
			paramName: "testSecret",
			service: &mockSecretService{
				DeleteSecretFunc: func(ctx context.Context, name string) error {
					return nil
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Internal server error",
			method:    http.MethodDelete,
			paramName: "testSecret",
			service: &mockSecretService{
				DeleteSecretFunc: func(ctx context.Context, name string) error {
					return errors.New("unexpected error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, fmt.Sprintf("/api/user/secret/%s", tt.paramName), nil)

			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("name", tt.paramName)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			rec := httptest.NewRecorder()

			handler := DeleteSecretHandler(logger, tt.service)
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}
