package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bunkieproject/bunkie_be/controllers"
)

type RequestResetPassword struct {
	Email string `json:"email"`
}

type ResetPasswordInput struct {
	Request      RequestResetPassword `json:"request"`
	Testname     string               `json:"testname"`
	ResponseCode int                  `json:"response_code"`
}

func TestResetPassword(t *testing.T) {
	// Setup
	router := gin.Default()
	router.POST("/users/reset_password", controllers.ResetPassword())
	reset_password_entries := []ResetPasswordInput{
		ResetPasswordInput{
			Testname:     "Correct reset password",
			ResponseCode: http.StatusOK,
			Request: RequestResetPassword{
				Email: "testuser@gmail.com",
			},
		},
		ResetPasswordInput{
			Testname:     "No email",
			ResponseCode: http.StatusBadRequest,
			Request: RequestResetPassword{
				Email: "",
			},
		},
		ResetPasswordInput{
			Testname:     "Email is not valid",
			ResponseCode: http.StatusBadRequest,
			Request: RequestResetPassword{
				Email: "testuser",
			},
		},
		ResetPasswordInput{
			Testname:     "Email is not registered",
			ResponseCode: http.StatusBadRequest,
			Request: RequestResetPassword{
				Email: "testuser111@gmail.com",
			},
		},
	}

	for _, entry := range reset_password_entries {
		// Setup
		w := httptest.NewRecorder()
		jsonValue, _ := json.Marshal(entry.Request)
		req, _ := http.NewRequest("POST", "/users/reset_password", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		// Execute
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, entry.ResponseCode, w.Code, entry.Testname)
	}
}
