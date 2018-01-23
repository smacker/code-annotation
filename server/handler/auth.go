package handler

import (
	"fmt"
	"net/http"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/service"
)

// Login handler redirects user to oauth provider
func Login(oAuth *service.OAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := oAuth.MakeAuthURL()
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// OAuthCallback makes exchange with oauth provider, gets&creates user and redirects to index page with JWT token
func OAuthCallback(oAuth *service.OAuth, jwt *service.JWT, userRepo *repository.Users, uiDomain string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := oAuth.ValidateState(r.FormValue("state")); err != nil {
			writeResponse(w, respErr(http.StatusBadRequest, err))
			return
		}

		code := r.FormValue("code")
		user, err := oAuth.GetUser(r.Context(), code)
		if err != nil {
			// FIXME can it be not server error? for wrong code
			writeResponse(w, resp500f("oauth get user error: %s", err))
			return
		}

		_, err = userRepo.Get(user.ID)
		if err != nil && err != repository.NotFoundErr {
			writeResponse(w, resp500f("can't get user: %s", err))
			return
		}
		if err == repository.NotFoundErr {
			if err := userRepo.Create(user); err != nil {
				writeResponse(w, resp500f("can't create user: %s", err))
				return
			}
		}

		token, err := jwt.MakeToken(user)
		if err != nil {
			writeResponse(w, resp500f("make jwt token error: %s", err))
			return
		}
		url := fmt.Sprintf("%s/?token=%s", uiDomain, token)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
