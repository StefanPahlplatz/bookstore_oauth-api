package rest

import (
	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid rest response when trying to log in user", err.Message)
}

func TestInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid credentials", "status": "404", "error": "not_found"}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to log in user", err.Message)
}

func TestInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid credentials", "status": 404, "error": "not_found"}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid credentials", err.Message)
}

func TestInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": "11","first_name": "Stefan","last_name": "Pahlplatz","email": "stefan@live.nl","created_on": "2020-03-11 18:56:30","status": "active"}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying to unmarshall user login response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 1,"first_name": "Stefan","last_name": "Pahlplatz","email": "stefan@live.nl","created_on": "2020-03-11 18:56:30","status": "active"}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "Stefan", user.FirstName)
	assert.EqualValues(t, "Pahlplatz", user.LastName)
	assert.EqualValues(t, "stefan@live.nl", user.Email)
}