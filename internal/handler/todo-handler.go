package handler

import (
	"encoding/json"
	"fmt"
	constants "go-server/internal/const"
	utils "go-server/pkg/helper"
	"net/http"
	"strings"
	"time"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/google/uuid"
)

type CreateTodoPayload struct {
	Todo        string `json:"todo"`
	IsCompleted bool   `json:"is_completed"`
}

type TodoResponseData struct {
	CreateTodoPayload
	Id        string     `json:"id"`
	UserId    string     `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
}

type GetAllTodoResp struct {
	Data  []*TodoResponseData `json:"data"`
	Count int                 `json:"count"`
}

// @title Get Todo API
// @Consume json
// @Tags Todo
// @Success 200 {object} GetAllTodoResp
// @Produce json
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Router /api/v1/get-all-todo [get]
func (a *AuthHandler) GetAllTodo(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	claims := utils.GetClaimFromCtx(r)

	query := `select
		t.id,
		t.user_id,
		t.todo,
		t.is_completed,
		t.created_at,
		t.updated_at,
		u.name,
		u.email
		from todo t join users u on u.id=t.user_id where t.user_id=$1`

	if !stringUtils.IsBlank(limit) {
		query += ` limit ` + limit
	}

	if !stringUtils.IsBlank(offset) {
		query += ` offset ` + offset
	}

	rows, err := a.Conn.Query(r.Context(), query, claims.Id)

	if err != nil {
		writeResponse(w, http.StatusNotFound, NewError("No Data Found", err.Error()))
		return
	}
	var resp []*TodoResponseData
	for rows.Next() {
		var Todo TodoResponseData
		rows.Scan(
			&Todo.Id,
			&Todo.UserId,
			&Todo.Todo,
			&Todo.IsCompleted,
			&Todo.CreatedAt,
			&Todo.UpdatedAt,
			&Todo.Name,
			&Todo.Email,
		)
		resp = append(resp, &Todo)
	}

	response := GetAllTodoResp{Data: resp, Count: len(resp)}
	writeResponse(w, http.StatusOK, response)

}

// @title Create Todo
// @Consume json
// @Tags Todo
// @Success 200 {object} Response
// @Produce json
// @Param request body CreateTodoPayload true "Request Body"
// @Router /api/v1/create-todo [post]
func (a *AuthHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var payload CreateTodoPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeResponse(w, http.StatusBadRequest, NewError("Invalid Json", err.Error()))
		return
	}

	claim := utils.GetClaimFromCtx(r)
	_, sqlErr := a.Conn.Exec(r.Context(),
		"insert into todo (id, user_id, todo, is_completed, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)",
		uuid.New(),
		claim.Id,
		payload.Todo,
		payload.IsCompleted,
		time.Now(),
		time.Now(),
	)
	if sqlErr != nil {
		writeResponse(w, http.StatusBadRequest, NewError("failed to create todo", sqlErr.Error()))
		return
	}
	resp := Response{Status: http.StatusOK, Message: "New Todo is Created"}
	writeResponse(w, http.StatusOK, resp)

}

type UpdateTodoPayload struct {
	Id         string  `json:"id"`
	Todo       *string `json:"todo"`
	IsComplted *bool   `json:"is_completed"`
}

// @title Update Todo
// @Consume json
// @Tags Todo
// @Success 200 {object} Response
// @Produce json
// @Param request body UpdateTodoPayload true "Request Body"
// @Router /api/v1/update-todo [post]
func (a *AuthHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var payload UpdateTodoPayload

	if payloadErr := json.NewDecoder(r.Body).Decode(&payload); payloadErr != nil {
		writeResponse(w, http.StatusBadRequest, NewError(constants.PAYLOAD_ERROR, payloadErr.Error()))
		return
	}

	fmt.Println("---> Todo ", *payload.Todo)
	// fmt.Println("---> Is Completed ", *payload.IsComplted)

	if stringUtils.IsBlank(payload.Id) {
		writeResponse(w, http.StatusBadRequest, NewError("Todo Id is required", constants.PAYLOAD_ERROR))
		return
	}

	if payload.Todo == nil && payload.IsComplted == nil {
		writeResponse(w, http.StatusBadRequest, NewError("No data to update", constants.PAYLOAD_ERROR))
		return
	}

	claim := utils.GetClaimFromCtx(r)

	var conditions []string

	if payload.Todo != nil {
		conditions = append(conditions, fmt.Sprintf("todo='%s'", *payload.Todo))
	}

	if payload.IsComplted != nil {
		conditions = append(conditions, fmt.Sprintf("is_completed='%t'", *payload.IsComplted))
	}

	query := fmt.Sprintf(`update todo set %s, updated_at=$1 where id=$2 and user_id=$3`, strings.Join(conditions, ","))

	fmt.Println(query)

	cmd, sqlErr := a.Conn.Exec(r.Context(), query, time.Now(), payload.Id, claim.Id)

	if sqlErr != nil {
		writeResponse(w, http.StatusBadRequest, NewError("Failed to update", sqlErr.Error()))
		return
	}

	if cmd.RowsAffected() == 0 {
		writeResponse(w, http.StatusBadRequest, NewError("Failed to update", "0 Rows Effected"))
		return
	}
	res := Response{Status: http.StatusOK, Message: "Todo has been updated"}
	writeResponse(w, http.StatusOK, res)
}

type DeletePayload struct {
	Id string `json:"id"`
}

// @title Delete Todo
// @Consume json
// @Tags Todo
// @Success 200 {object} Response
// @Produce json
// @Param request body DeletePayload true "Request Body"
// @Router /api/v1/delete-todo [post]
func (a *AuthHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	var payload DeletePayload

	payloadErr := json.NewDecoder(r.Body).Decode(&payload)
	if payloadErr != nil || stringUtils.IsBlank(payload.Id) {
		writeResponse(w, http.StatusBadRequest, constants.PAYLOAD_ERROR)
		return
	}
	claims := utils.GetClaimFromCtx(r)
	query := `delete from todo where id=$1 and user_id=$2`

	cmdTag, _ := a.Conn.Exec(r.Context(), query, payload.Id, claims.Id)
	if cmdTag.RowsAffected() > 0 {
		writeResponse(w, http.StatusOK, Response{Status: http.StatusOK, Message: "Deleted Successfully"})
		return
	}
	writeResponse(w, http.StatusBadRequest, NewError("Failed to delete", "Failed to Delete"))
}
