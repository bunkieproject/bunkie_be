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

type RequestSignup struct {
	Username         string `json:"username"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Password_confirm string `json:"password_confirm"`
}

type SignupInput struct {
	Request      RequestSignup `json:"request"`
	Testname     string        `json:"testname"`
	ResponseCode int           `json:"response_code"`
}

func TestSignUp(t *testing.T) {
	// Setup
	router := gin.Default()
	router.POST("/users/signup", controllers.SignUp())
	signup_entries := []SignupInput{
		SignupInput{
			Testname:     "Correct signup",
			ResponseCode: http.StatusOK,
			Request: RequestSignup{
				Username:         "testuser",
				Email:            "testuser@gmail.com",
				Password:         "testuser",
				Password_confirm: "testuser",
			},
		},
		SignupInput{
			Testname:     "Short username",
			ResponseCode: http.StatusBadRequest,
			Request: RequestSignup{
				Username:         "tes",
				Email:            "testuser1@gmail.com",
				Password:         "testuser1",
				Password_confirm: "testuser1",
			},
		},
		SignupInput{
			Testname:     "Long username",
			ResponseCode: http.StatusBadRequest,
			Request: RequestSignup{
				Username:         "testtesttesttesttesttesttesttesttesttest",
				Email:            "testuser2@gmail.com",
				Password:         "testuser2",
				Password_confirm: "testuser2",
			},
		},
		SignupInput{
			Testname:     "No username",
			ResponseCode: http.StatusBadRequest,
			Request: RequestSignup{
				Email:            "testuser3@gmail.com",
				Password:         "testuser3",
				Password_confirm: "testuser3",
			},
		},
		SignupInput{
			Testname:     "No email",
			ResponseCode: http.StatusBadRequest,
			Request: RequestSignup{
				Username:         "testuser4",
				Password:         "testuser4",
				Password_confirm: "testuser4",
			},
		},
		SignupInput{
			Testname:     "Short password",
			ResponseCode: http.StatusBadRequest,
			Request: RequestSignup{

				Username:         "testuser5",
				Email:            "testuser5@gmail.com",
				Password:         "tes",
				Password_confirm: "tes",
			},
		},
		SignupInput{
			Testname:     "Long password",
			ResponseCode: http.StatusBadRequest,
			Request: RequestSignup{
				Username:         "testuser6",
				Email:            "testuser6@gmail.com",
				Password:         "testtesttesttesttesttesttesttesttesttest",
				Password_confirm: "testtesttesttesttesttesttesttesttesttest",
			},
		},
		SignupInput{
			Testname:     "No password",
			ResponseCode: http.StatusBadRequest,
			Request: RequestSignup{
				Username:         "testuser7",
				Email:            "testuser7@gmail.com",
				Password_confirm: "testuser7",
			},
		},
		SignupInput{
			Testname:     "No password confirm",
			ResponseCode: http.StatusBadRequest,
			Request: RequestSignup{
				Username: "testuser8",
				Email:    "testuser8@gmail.com",
				Password: "testuser8",
			},
		},
		SignupInput{
			Testname:     "Password and password confirm don't match",
			ResponseCode: http.StatusBadRequest,
			Request: RequestSignup{
				Username:         "testuser9",
				Email:            "testuser9@gmail.com",
				Password:         "testuser9",
				Password_confirm: "testuser999",
			},
		},
		SignupInput{
			Testname:     "Email already exists",
			ResponseCode: http.StatusBadRequest,
			Request: RequestSignup{
				Username:         "testuser10",
				Email:            "testuser10@gmail.com",
				Password:         "testuser10",
				Password_confirm: "testuser10",
			},
		},
		SignupInput{
			Testname:     "Username already exists",
			ResponseCode: http.StatusBadRequest,
			Request: RequestSignup{
				Username:         "testuser11",
				Email:            "testuser11@gmail.com",
				Password:         "testuser11",
				Password_confirm: "testuser11",
			},
		},
	}

	for _, signup_entry := range signup_entries {
		// Execute
		jsonValue, _ := json.Marshal(signup_entry.Request)
		req, _ := http.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, signup_entry.ResponseCode, w.Code, signup_entry.Testname)
		log.Println(w.Body)
	}
}
