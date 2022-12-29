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

type RequestLogin struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type LoginInput struct {
	Request      RequestLogin `json:"request"`
	Testname     string       `json:"testname"`
	ResponseCode int          `json:"response_code"`
}

func TestLogin(t *testing.T) {
	// Setup
	router := gin.Default()
	router.POST("/users/login", controllers.Login())
	login_entries := []LoginInput{
		LoginInput{
			Testname:     "Correct login",
			ResponseCode: http.StatusOK,
			Request: RequestLogin{
				UsernameOrEmail: "testuser",
				Password:        "testuser",
			},
		},
		LoginInput{
			Testname:     "No username",
			ResponseCode: http.StatusBadRequest,
			Request: RequestLogin{
				Password: "testuser",
			},
		},
		LoginInput{
			Testname:     "No password",
			ResponseCode: http.StatusBadRequest,
			Request: RequestLogin{
				UsernameOrEmail: "testuser",
			},
		},
		LoginInput{
			Testname:     "Wrong password",
			ResponseCode: http.StatusBadRequest,
			Request: RequestLogin{
				UsernameOrEmail: "testuser",
				Password:        "testuser1",
			},
		},
		LoginInput{
			Testname:     "Wrong username",
			ResponseCode: http.StatusBadRequest,
			Request: RequestLogin{
				UsernameOrEmail: "testuser1",
				Password:        "testuser",
			},
		},
	}

	for _, entry := range login_entries {
		jsonValue, _ := json.Marshal(entry.Request)
		req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, entry.ResponseCode, w.Code, entry.Testname)
		log.Println(entry.Testname, w.Body.String())
	}
}
