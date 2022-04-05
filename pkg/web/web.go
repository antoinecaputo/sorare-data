package web

import (
	"bytes"
	"errors"
	"github.com/bamboutech/golog"
	"io/ioutil"
	"net/http"
	"sorare-data/constants"
)

func FctHTTPGet(_url string) ([]byte, error) {
	// Send request
	constants.Logger.FctLog(golog.LogLvl_Debug, "Sending HTTP request")
	res, err := http.Get(_url)
	if err != nil {
		constants.Logger.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		return nil, errors.New("http request failed")
	}
	defer res.Body.Close()
	constants.Logger.FctLog(golog.LogLvl_Debug, "   = OK")

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, errors.New("http request failed with status code " + res.Status)
	}

	// Reading response body
	constants.Logger.FctLog(golog.LogLvl_Debug, "Reading response body")
	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		constants.Logger.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		return nil, errors.New("can't read response body")
	}

	return resBodyBytes, nil
}

func FctJWT(_email string, _hashedPassword string) ([]byte, error) {

	/*
		var jsonBytes = []byte(`{
		"operationName": "SignInMutation"
		"variables": { "input": { "email": "` + _email + `", "password": "` + _hashedPassword + `" } },
		"query": "mutation SignInMutation($input: signInInput!) { signIn(input: $input) { currentUser { slug jwtToken(aud: "sorare-data") { token expiredAt } } errors { message } } }"
		}`)
	*/

	var jsonBytes = []byte(`{}`)

	// Send request
	constants.Logger.FctLog(golog.LogLvl_Debug, "Sending HTTP request")
	res, err := http.Post("https://api.sorare.com/graphql", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		constants.Logger.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		return nil, errors.New("http request failed")
	}
	defer res.Body.Close()
	constants.Logger.FctLog(golog.LogLvl_Debug, "   = OK")

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, errors.New("http request failed with status code " + res.Status)
	}

	// Reading response body
	constants.Logger.FctLog(golog.LogLvl_Debug, "Reading response body")
	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		constants.Logger.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		return nil, errors.New("can't read response body")
	}

	return resBodyBytes, nil
}
