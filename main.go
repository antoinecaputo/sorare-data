package main

import (
	"encoding/json"
	"fmt"
	"github.com/bamboutech/golog"
	"log"
	"sorare-data/constants"
	"sorare-data/pkg/crypt"
	"sorare-data/pkg/web"
)

func main() {
	logger, err := golog.FctCreateLogger(golog.TrcMth_Buffer, golog.LogLvl_Debug, "", "")
	if err != nil {
		log.Fatal(err)
		return
	}
	constants.Logger = logger

	res, err := web.FctHTTPGet("https://api.sorare.com/api/v1/users/" + constants.Email)
	if err != nil {
		fctFatal(err.Error())
		return
	}

	// Decoding response
	constants.Logger.FctLog(golog.LogLvl_Debug, "Decoding response")
	var resMap map[string]interface{}
	err = json.Unmarshal(res, &resMap)
	if err != nil {
		constants.Logger.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		fctFatal("can't decode response")
		return
	}
	salt := resMap["salt"].(string)
	constants.Logger.FctLog(golog.LogLvl_Debug, "   = OK, got %s", salt)

	// Hashing password
	constants.Logger.FctLog(golog.LogLvl_Debug, "Hashing password")
	hashStr, err := crypt.FctHashSecret(constants.UnhashedPassword, salt)
	if err != nil {
		constants.Logger.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		fctFatal("can't hash password")
		return
	}
	constants.Logger.FctLog(golog.LogLvl_Debug, "   = OK, got %s", hashStr)

	// Requesting JWT token
	constants.Logger.FctLog(golog.LogLvl_Debug, "Requesting JWT token")
	res, err = web.FctJWT(constants.Email, hashStr)
	if err != nil {
		constants.Logger.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		fctFatal("can't get JWT token")
		return
	}
	constants.Logger.FctLog(golog.LogLvl_Debug, "   = OK, got %s", res)

	buffStr, err := constants.Logger.FctGetBufferContent()
	if err != nil {
		panic(err)
		return
	}

	fmt.Println(buffStr)
}

func fctFatal(error string) int {
	constants.Logger.FctLog(golog.LogLvl_Emerg, "FATAL ERROR = %s", error)

	buffStr, err := constants.Logger.FctGetBufferContent()
	if err != nil {
		panic(err)
		return -1
	}

	fmt.Println(buffStr)
	return 0
}
