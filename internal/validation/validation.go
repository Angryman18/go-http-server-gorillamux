package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	errs "go-server/internal/error"
	"net/http"

	"github.com/agrison/go-commons-lang/stringUtils"
)

type RequestBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateSignupInfo(r *http.Request) (RequestBody, error) {
	var requestBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Println(err.Error())
		return requestBody, errs.NewError("Invalid Request")
	}

	if stringUtils.IsBlank(requestBody.Name) {
		return requestBody, errs.NewError("Name is Required")
	}

	if stringUtils.IsBlank(requestBody.Email) {
		return requestBody, errs.NewError("Email is Required")
	}

	// re := regexp.MustCompile(constants.Regex)
	// if !re.MatchString(requestBody.Password) {
	// 	return requestBody, errs.NewError("Password must contain at least 8 letters contaning capital, small, special char and numbers")
	// }
	return requestBody, nil

}

type LoginReqestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateLoginInfo(r *http.Request) (LoginReqestBody, error) {
	var payload LoginReqestBody

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return payload, nil
	}

	if stringUtils.IsBlank(payload.Email) || stringUtils.IsBlank(payload.Password) {
		return payload, errors.New("email or password cannot be emprty")
	}
	return payload, nil
}
