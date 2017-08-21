// Copyright © 2017 Ott-Consult UG (haftungsbeschränkt), Jörn Ott <raspibot@ott-consult.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	//"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/joernott/go-sbc-motorshield/motor"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/joernott/raspibot/statik"
	"github.com/julienschmidt/httprouter"
	"github.com/rakyll/statik/fs"
	"golang.org/x/crypto/bcrypt"
)

type UserList map[string]string

var Motor motor.Motor
var statikFS http.FileSystem
var Users UserList

func main() {
	var err error
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.Debug("Initializing")
	Motor, err = motor.NewMotor()
	if err != nil {
		panic(err)
	}
	defer Motor.CloseMotor()

	Users = initBasicAuth()

	router := httprouter.New()
	router.POST("/api/v1/drive/forward", BasicAuth(DriveForward))
	router.POST("/api/v1/drive/reverse", BasicAuth(DriveReverse))
	router.POST("/api/v1/drive/turnleft", BasicAuth(TurnLeft))
	router.POST("/api/v1/drive/turnright", BasicAuth(TurnRight))
	router.POST("/api/v1/camera/up", BasicAuth(CameraUp))
	router.POST("/api/v1/camera/down", BasicAuth(CameraDown))
	router.POST("/api/v1/stop", BasicAuth(AllStop))
	router.GET("/api/v1/stop", BasicAuth(AllStop))
	statikFS, err = fs.New()
	if err != nil {
		log.WithField("context", "initializing").Error("Could not initialize http filesystem.")
		panic("Could not get absolute path of the application.")

	}
	//d, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//if err != nil {
	//	log.WithField("context", "initializing").Error("Could not get absolute path of the application.")
	//	panic("Could not get absolute path of the application.")
	//}

	router.ServeFiles("/static/*filepath", statikFS) //http.Dir(d+string(os.PathSeparator)+"static")
	router.GET("/", BasicAuth(Index))

	log.Fatal(http.ListenAndServe(":80", router))
}

func initBasicAuth() UserList {
	var u UserList
	u = make(map[string]string)
	f, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.WithField("context", "init basic auth").Fatal("Could not open users.json.")
		return u
	}
	err = json.Unmarshal(f, &u)
	if err != nil {
		log.WithField("context", "init basic auth").Fatal(err)
		return u
	}
	spew.Dump(u)
	return u
}

func BasicAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		const basicAuthPrefix string = "Basic "
		logger := log.WithField("context", "basic auth").Logger
		// Get the Basic Authentication credentials
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
			if err != nil {
				logger.Error(err)
			} else {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) != 2 {
					logger.Error("Could not split credentials")
				} else {
					pass := pair[1]
					username := string(pair[0][:])
					if opass, ok := Users[username]; ok == false {
						logger.Warn("Could not find user " + username)
					} else {
						err = bcrypt.CompareHashAndPassword([]byte(opass), pass)
						if err != nil {
							logger.Warn("Password mismatch for user " + username)
							logger.Debug(err)
						} else {
							// Delegate request to the given handle
							logger.Debug("Successfully authenticated " + username)
							h(w, r, ps)
							return
						}
					}
				}
			}
		}

		// Request Basic Authentication otherwise
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}
