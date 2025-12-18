package response

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	logger "github.com/he-end/simproute/route_logger"
	"go.uber.org/zap"
)

// Response represents the standard API response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    Meta        `json:"meta"`
}

type ResponseNoContent string

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string `json:"code"`
	Details string `json:"details"`
}

// Meta contains metadata about the response
type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}

// ResponseHandler handles API responses
type ResponseHandler struct {
	logger *zap.Logger
	dev    bool
}

// NewWithGlobalLogger creates a new ResponseHandler using the global logger
func NewWithGlobalLogger() *ResponseHandler {
	dev := os.Getenv("ENV") == "development" || os.Getenv("ENV") == "dev"
	return &ResponseHandler{
		logger: logger.GetLogger(),
		dev:    dev,
	}
}

// Success sends a successful response
func (rh *ResponseHandler) Success(w http.ResponseWriter, message string, data interface{}) {
	// requestID := uuid.New().String()

	response := Response{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta: Meta{
			// RequestID: requestID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// rh.logger.Info("API Success Response",
	// zap.String("request_id", requestID),
	// zap.String("message", message),
	// zap.String("status", "success"),
	// )
	//
	rh.writeJSON(w, http.StatusOK, response)
}

// Success sends a successful Accepted
func (rh *ResponseHandler) Accepted(w http.ResponseWriter, message string, data interface{}) {
	// requestID := uuid.New().String()
	response := Response{
		Status:  "accepted",
		Message: message,
		Data:    data,
		Meta: Meta{
			// RequestID: requestID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// rh.logger.Info("API Success Response",
	// zap.String("request_id", requestID),
	// zap.String("message", message),
	// zap.String("status", "success"),
	// )
	//
	rh.writeJSON(w, http.StatusAccepted, response)
}

func (rh *ResponseHandler) Created(w http.ResponseWriter, resoureceLocation string, message string, data interface{}) {
	// requestID := uuid.New().String()

	response := Response{
		Status:  "created",
		Message: message,
		Data:    data,
		Meta: Meta{
			// RequestID: requestID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// rh.logger.Info("API Success Response",
	// zap.String("request_id", requestID),
	// zap.String("message", message),
	// zap.String("status", "created"),
	// )
	//
	w.Header().Add("Location", resoureceLocation)
	rh.writeJSON(w, http.StatusCreated, response)
}

// Success sends a successful response
func (rh *ResponseHandler) SuccessNoContent(w http.ResponseWriter, message string) {
	// requestID := uuid.New().String()

	// rh.logger.Info("API Success Response",
	// zap.String("request_id", requestID),
	// zap.String("message", message),
	// zap.String("status", "success"),
	// zap.String("method", http.MethodDelete),
	// )
	rh.writeNoContent(w, http.StatusNoContent)
	// rh.writeJSON(w, http.StatusNoContent, ResNoContent)
}

// Fail sends a failure response (business logic failure)
func (rh *ResponseHandler) Fail(w http.ResponseWriter, message string, errCode string, details string) {
	// requestID := uuid.New().String()

	errorInfo := &ErrorInfo{
		Code:    errCode,
		Details: details,
	}

	response := Response{
		Status:  "fail",
		Message: message,
		Error:   errorInfo,
		Meta: Meta{
			// RequestID: requestID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// rh.logger.Warn("API Fail Response",
	// zap.String("request_id", requestID),
	// zap.String("message", message),
	// zap.String("error_code", errCode),
	// zap.String("status", "fail"),
	// )
	rh.writeJSON(w, http.StatusBadRequest, response)
}

// Error sends an error response (system/server error)
func (rh *ResponseHandler) Error(w http.ResponseWriter, message string, errCode string, details string, httpStatus int) {
	// requestID := uuid.New().String()

	errorInfo := &ErrorInfo{
		Code:    errCode,
		Details: details,
	}

	response := Response{
		Status:  "error",
		Message: message,
		Error:   errorInfo,
		Meta: Meta{
			// RequestID: requestID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// In production, don't expose sensitive error details
	if !rh.dev && httpStatus >= 500 {
		response.Message = "Internal server error"
		response.Error.Details = "An unexpected error occurred"
	}

	// rh.logger.Error("API Error Response",
	// 	zap.String("request_id", requestID),
	// 	zap.String("message", message),
	// 	zap.String("error_code", errCode),
	// 	zap.Int("http_status", httpStatus),
	// 	zap.String("status", "error"),
	// )

	rh.writeJSON(w, httpStatus, response)
}

// writeJSON writes JSON response to the response writer
func (rh *ResponseHandler) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		rh.logger.Error("Failed to encode JSON response",
			zap.Error(err),
		)
	}
}

// writeJSON writes JSON response to the response writer
func (rh *ResponseHandler) writeNoContent(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

const (
	ResNoContent ResponseNoContent = "NO_CONTENT"
)

// Common error codes
const (
	ErrCodeInvalidJSON                = "INVALID_JSON"
	ErrCodeValidationError            = "VALIDATION_ERROR"
	ErrCodeInvalidCredentials         = "INVALID_CREDENTIALS"
	ErrCodeEmailInQueue               = "EMAIL_IN_QUEUE"
	ErrCodeEmailRegistered            = "EMAIL_REGISTERED"
	ErrCodeMissingToken               = "MISSING_TOKEN"
	ErrCodeVerificationFailed         = "VERIFICATION_FAILED"
	ErrCodeInternalError              = "INTERNAL_ERROR"
	ErrCodeDatabaseError              = "DATABASE_ERROR"
	ErrCodeEmailError                 = "EMAIL_ERROR"
	ErrCodeInvalidRefreshToken        = "INVALID_REFRESH_TOKEN"
	ErrCodeRefreshTokenExpired        = "REFRESH_TOKEN_EXPIRED"
	ErrCodeMissingAuthHeader          = "MISSING_AUTH_HEADER"
	ErrCodeInvalidAuthFormat          = "INVALID_AUTH_FORMAT"
	ErrCodeInvalidToken               = "INVALID_TOKEN"
	ErrCodeResendFailed               = "RESEND_FAILED"
	ErrCodeForgotPasswordFailed       = "FORGOT_PASSWORD_FAILED"
	ErrCodeResetPasswordFailed        = "RESET_PASSWORD_FAILED"
	ErrCodePasswordMismatch           = "PASSWORD_MISMATCH"
	ErrCodeErrSignup                  = "SIGNUP_ERROR"
	ErrCodeResendSignupFailed         = "RESEND_SIGNUP_FAILED"
	ErrCodeResendSigninFailed         = "RESEND_SIGNIN_FAILED"
	ErrCodeResendForgotPasswordFailed = "RESEND_FORGOT_PASSWORD_FAILED"
	ErrCodeResendAccountClosureFailed = "RESEND_ACCOUNT_CLOSURE_FAILED"
	ErrCodeAlreadyUsed                = "ALREADY_USED"
	ErrCodeUserAlreadyUse             = "USER_ALREADY_USE"
	ErrCodeNoFound                    = "NO_FOUND"
	ErrCodeInvalidRequest             = "INVALID_REQUEST"
	ErrCodeInvalidUUID                = "INVALID_UUID"
	ErrCodeTokenExpired               = "TOKEN_EXPIRED"
	ErrCodeInvalidHeader              = "INVALID_HEADER"
	ErrCodeURLInvalid                 = "INVALID_URL"
	ErrCodeDuplicatKey                = "DUPLICAT_KEY"
	ErrCodeErrUpdate                  = "ERROR_UPDATE"
	ErrCodeErrCreate                  = "ERROR_CREATE"
	ErrCodeErrDelete                  = "ERROR_DELETE"
	ErrCodeAlreadyExist               = "ALREADY_EXISTS"
	ErrCodeTypeUnsupported            = "TYPE_UNSUPPORTED"
	ErrCodeRetrieve                   = "RETRIEVE_ERROR"
	ErrCodeUpdateError                = "UPDATE_ERROR"
	ErrCodeNoFieldsUpdate             = "NO_FIELDS_UPDATE"
	ErrCodePayloadEmpty               = "PAYLOAD_EMPTY"
	ErrCodeMissingFieldJSON           = "MISSING_FIELDS_JSON"
	ErrCodeNoDataFound                = "NO_DATA"
	ErrCodeUnset                      = "UNSET"
)

// Common messages
const (
	// Success Messages
	MsgUpdateEndpointSuccess = "update endpoint success"
	MsgSignupSuccess         = "Registration successful, please check your email for verification"
	MsgSigninSuccess         = "Login successful"
	MsgVerificationSuccess   = "Email verification successful"
	MsgTokenRefreshSuccess   = "Token refreshed successfully"
	MsgSignoutSuccess        = "Signout successful"
	MsgUserDataRetrieved     = "User data retrieved successfully"
	MsgVerificationResent    = "Verification email sent successfully"
	MsgPasswordResetSuccess  = "Password reset successfully"
	MsgUpdateSuccess         = "Update successful"
	MsgDeleteSuccess         = "Delete successful"

	// Error messages
	MsgInvalidJSON         = "Invalid JSON format"
	MsgUserAlreadyUse      = "User already in used, please use another username"
	MsgValidationError     = "Validation failed"
	MsgInvalidCredentials  = "Invalid email or password"
	MsgEmailInQueue        = "Email is already in queue, please check your email"
	MsgUrlInvalid          = "url invalid"
	MsgEmailRegistered     = "Email already registerd"
	MsgEmailSame           = "Same email"
	MsgMissingToken        = "Verification token is required"
	MsgVerificationFailed  = "Email verification failed"
	MsgInternalError       = "Internal server error"
	MsgInvalidRefreshToken = "Invalid or expired refresh token"
	MsgPasswordResetSent   = "a password reset link has been sent"
	MsgPasswordMismatch    = "Password do not match"
	MsgNoFound             = "Resource no found"
	MsgUnauthorized        = "Unauthorized"
	MsgInvalidUUID         = "Invalid uuid"
	MsgTokenExpired        = "Token expired"
	MsgInvalidToken        = "Invalid token"
	MsgInvalidHeader       = "Invalid header"
	MsgDuplicatKey         = "duplicat key"
	MsgAPINotFound         = "APIs not found"
	MsgDefinitionNotFound  = "Definition not found"
	MsgDeleteError         = "Delete error"
	MsgCreateError         = "Create error"
	MsgNoFieldUpdate       = "no field update"
	MsgUpdateError         = "update error, please check your data and try again"
)
