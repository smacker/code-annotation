package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/src-d/code-annotation/server/model"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/service"
)

func Experiment(repo *repository.Experiments) http.HandlerFunc {
	return render(func(r *http.Request) response {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			return respErr(http.StatusBadRequest, err)
		}

		exp, err := repo.Get(id)
		if err == repository.NotFoundErr {
			return respStringErr(http.StatusNotFound, "repository not found")
		}
		if err != nil {
			return resp500f("get experiment: %s", err)
		}

		return respOK(exp)
	})
}

func Assignments(repo *repository.Experiments) http.HandlerFunc {
	return render(func(r *http.Request) response {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			return respErr(http.StatusBadRequest, err)
		}
		uID := service.GetUserID(r.Context())
		if uID == 0 {
			return resp500(errors.New("no user id in context"))
		}

		yes, err := repo.HasAssignments(id, uID)
		if err != nil {
			return resp500f("has assignments: %s", err)
		}
		if !yes {
			if err := repo.CreateAssignments(id, uID); err != nil {
				return resp500f("create assignments: %s", err)
			}
		}
		as, err := repo.GetAssignments(id, uID)
		if err != nil {
			return resp500f("get assignments error: %s", err)
		}

		return respOK(as)
	})
}

type assignmentInput struct {
	// TODO Add validation here
	Answer   string `json:"answer"`
	Duration int    `json:"duration"`
}

func UpdateAssignment(repo *repository.Experiments) http.HandlerFunc {
	return render(func(r *http.Request) response {
		uID := service.GetUserID(r.Context())
		if uID == 0 {
			return resp500(errors.New("no user id in context"))
		}
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			return respErr(http.StatusBadRequest, err)
		}
		var input assignmentInput
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			return respErr(http.StatusBadRequest, err)
		}
		// TODO actually we need to check experiment id also
		a, err := repo.GetAssignment(id)
		if err == repository.NotFoundErr {
			return respStringErr(http.StatusNotFound, "assignment not found")
		}
		if err != nil {
			return resp500f("get assignment: %s", err)
		}

		a.Answer = model.ToNullString(input.Answer)
		a.Duration = model.ToNullInt64(input.Duration)
		err = repo.UpdateAssignments(a)
		if err != nil {
			return resp500f("update assignments: %s", err)
		}

		return response{statusCode: http.StatusNoContent}
	})
}

func FilePair(repo *repository.Experiments) http.HandlerFunc {
	return render(func(r *http.Request) response {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			return respErr(http.StatusBadRequest, err)
		}

		// TODO actually we need to check experiment id also
		fp, err := repo.GetFilePair(id)
		if err == repository.NotFoundErr {
			return respStringErr(http.StatusNotFound, "file pair not found")
		}
		if err != nil {
			return resp500f("get file pair: %s", err)
		}
		return respOK(fp)
	})
}
