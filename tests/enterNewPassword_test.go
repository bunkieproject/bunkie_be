package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bunkieproject/bunkie_be/controllers"
)

type RequestEnterNewPassword struct {
	Email              string `json:"email"`
	NewPassword        string `json:"new_password"`
	NewPasswordConfirm string `json:"new_password_confirm"`
}

type EnterNewPasswordInput struct {
	Request      RequestEnterNewPassword `json:"request"`
	Testname     string                  `json:"testname"`
	ResponseCode int                     `json:"response_code"`
}

func TestEnterNewPassword(t *testing.T) {
	// Setup
	router := gin.Default()
	router.POST("/users/enter_new_password", controllers.EnterNewPassword())
	enter_new_password_entries := []EnterNewPasswordInput{
		EnterNewPasswordInput{
			Testname:     "Correct enter new password",
			ResponseCode: http.StatusOK,
			Request: RequestEnterNewPassword{
				Email:              "testuser@gmail.com",
				NewPassword:        "testuser",
				NewPasswordConfirm: "testuser",
			},
		},
		EnterNewPasswordInput{
			Testname:     "No email",
			ResponseCode: http.StatusBadRequest,
			Request: RequestEnterNewPassword{
				Email:              "",
				NewPassword:        "testuser",
				NewPasswordConfirm: "testuser",
			},
		},
		EnterNewPasswordInput{
			Testname:     "Email is not valid",
			ResponseCode: http.StatusBadRequest,
			Request: RequestEnterNewPassword{
				Email:              "testuser",
				NewPassword:        "testuser",
				NewPasswordConfirm: "testuser",
			},
		},
		EnterNewPasswordInput{
			Testname:     "Email is not registered",
			ResponseCode: http.StatusBadRequest,
			Request: RequestEnterNewPassword{
				Email:              "testuser9999@gmail.com",
				NewPassword:        "testuser",
				NewPasswordConfirm: "testuser",
			},
		},
		EnterNewPasswordInput{
			Testname:     "No new password",
			ResponseCode: http.StatusBadRequest,
			Request: RequestEnterNewPassword{
				Email:              "testuser@gmail.com",
				NewPassword:        "",
				NewPasswordConfirm: "testuser",
			},
		},
		EnterNewPasswordInput{
			Testname:     "New password is too short",
			ResponseCode: http.StatusBadRequest,
			Request: RequestEnterNewPassword{
				Email:              "testuser@gmail.com",
				NewPassword:        "test",
				NewPasswordConfirm: "test",
			},
		},
		EnterNewPasswordInput{
			Testname:     "New password is too long",
			ResponseCode: http.StatusBadRequest,
			Request: RequestEnterNewPassword{
				Email:              "testuser@gmail.com",
				NewPassword:        "testtesttesttesttesttesttesttesttesttest",
				NewPasswordConfirm: "testtesttesttesttesttesttesttesttesttest",
			},
		},
		EnterNewPasswordInput{
			Testname:     "No new password confirm",
			ResponseCode: http.StatusBadRequest,
			Request: RequestEnterNewPassword{
				Email:              "testuser@gmail.com",
				NewPassword:        "testuser",
				NewPasswordConfirm: "",
			},
		},
	}

	for _, entry := range enter_new_password_entries {
		// Setup
		jsonValue, _ := json.Marshal(entry.Request)
		req, err := http.NewRequest("POST", "/users/enter_new_password", bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Execute
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, entry.ResponseCode, w.Code, entry.Testname)
	}
}
