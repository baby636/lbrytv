package users

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lbryio/lbrytv/models"
	"github.com/stretchr/testify/assert"
)

type DummyRetriever struct{}

func (r *DummyRetriever) Retrieve(t string) (*models.User, error) {
	if t == "XyZ" {
		return &models.User{SDKAccountID: "aBc"}, nil
	}
	return nil, errors.New("cannot authenticate")
}

func AuthenticatedHandler(w http.ResponseWriter, r *AuthenticatedRequest) {
	if r.IsAuthenticated() {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(r.AccountID))
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(r.AuthError.Error()))
	}
}

func TestAuthenticator(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/proxy", nil)
	r.Header.Set(TokenHeader, "XyZ")
	rr := httptest.NewRecorder()
	authenticator := NewAuthenticator(&DummyRetriever{})

	http.HandlerFunc(authenticator.Wrap(AuthenticatedHandler)).ServeHTTP(rr, r)

	response := rr.Result()
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "aBc", string(body))
}

func TestAuthenticatorFailure(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/proxy", nil)
	r.Header.Set(TokenHeader, "ALSDJ")
	rr := httptest.NewRecorder()

	authenticator := NewAuthenticator(&DummyRetriever{})

	http.HandlerFunc(authenticator.Wrap(AuthenticatedHandler)).ServeHTTP(rr, r)
	response := rr.Result()
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "cannot authenticate", string(body))
	assert.Equal(t, http.StatusForbidden, response.StatusCode)
}
