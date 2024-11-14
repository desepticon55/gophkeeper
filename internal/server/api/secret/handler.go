package secret

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

// Handler to upload user secret
func UploadSecretHandler(logger *zap.Logger, service secretService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, fmt.Sprintf("Method '%s' is not allowed", request.Method), http.StatusBadRequest)
			return
		}

		var secret model.Secret
		if err := json.NewDecoder(request.Body).Decode(&secret); err != nil {
			logger.Error("Error decode request", zap.Error(err))
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if err := request.Body.Close(); err != nil {
			logger.Error("Error closing response body", zap.Error(err))
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		err := service.CreateSecret(request.Context(), secret)
		if err != nil {
			if errors.Is(err, model.ErrSecretNameIsEmpty) {
				http.Error(writer, "Secret name is not filled", http.StatusBadRequest)
				return
			}

			if errors.Is(err, model.ErrSecretExistToCurrentUser) {
				http.Error(writer, "Secret already exists to current user", http.StatusConflict)
				return
			}

			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

// Handler to read user secret by secret name
func ReadOneSecretHandler(logger *zap.Logger, service secretService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			http.Error(writer, fmt.Sprintf("Method '%s' is not allowed", request.Method), http.StatusBadRequest)
			return
		}

		secretName := chi.URLParam(request, "name")
		if secretName == "" {
			http.Error(writer, "Secret name should be filled", http.StatusBadRequest)
			return
		}

		secret, err := service.FindSecret(request.Context(), secretName)
		if err != nil {
			if errors.Is(err, model.ErrSecretWasNotFound) {
				http.Error(writer, "Secret was not found", http.StatusNotFound)
				return
			}
			logger.Error("Error during find secret.", zap.String("secretName", secretName), zap.Error(err))
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(secret)
		if err != nil {
			logger.Error("Error during marshal secrets.", zap.Error(err))
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if _, err = writer.Write(bytes); err != nil {
			logger.Error("Error write secrets.", zap.Error(err))
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

// Handler to delete user secret by name
func DeleteSecretHandler(logger *zap.Logger, service secretService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodDelete {
			http.Error(writer, fmt.Sprintf("Method '%s' is not allowed", request.Method), http.StatusBadRequest)
			return
		}

		secretName := chi.URLParam(request, "name")
		if secretName == "" {
			http.Error(writer, "Secret name should be filled", http.StatusBadRequest)
			return
		}

		err := service.DeleteSecret(request.Context(), secretName)
		if err != nil {
			if errors.Is(err, model.ErrSecretWasNotFound) {
				http.Error(writer, "Secrets was not found", http.StatusNotFound)
				return
			}
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

// Handler to read all user secrets
func ReadAllSecretsHandler(logger *zap.Logger, service secretService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			http.Error(writer, fmt.Sprintf("Method '%s' is not allowed", request.Method), http.StatusBadRequest)
			return
		}

		secrets, err := service.FindAllSecrets(request.Context())
		if err != nil {
			if errors.Is(err, model.ErrSecretsWasNotFound) {
				http.Error(writer, "Secrets was not found", http.StatusNoContent)
				return
			}
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(secrets)
		if err != nil {
			logger.Error("Error during marshal secrets.", zap.Error(err))
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if _, err = writer.Write(bytes); err != nil {
			logger.Error("Error write secrets.", zap.Error(err))
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
