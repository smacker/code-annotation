package handler

import (
	"errors"
	"net/http"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/service"
)

// Me handler returns information about current user
func Me(usersRepo *repository.Users) http.HandlerFunc {
	return render(func(r *http.Request) response {
		uID := service.GetUserID(r.Context())
		if uID == 0 {
			return resp500(errors.New("no user id in context"))
		}
		u, err := usersRepo.Get(uID)
		if err == repository.NotFoundErr {
			return respStringErr(http.StatusNotFound, "user not found")
		}
		if err != nil {
			return resp500f("me: %s", err)
		}
		return respOK(u)
	})
}
