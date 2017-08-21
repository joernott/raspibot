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
// limitations under the License.package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/joernott/go-sbc-motorshield/motor"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type Action struct {
	Duration int `json:"duration"`
}

type ErrorResponse struct {
	Action  string `json:"action"`
	Context string `json:"context"`
	Message string `json:"message"`
}

func SendErrorResponse(w http.ResponseWriter, response ErrorResponse) {
	log.WithFields(log.Fields{"Action": response.Action, "Context": response.Context}).Error(response.Message)
	w.WriteHeader(http.StatusInternalServerError)
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(&response)
	if err == nil {
		w.Write(b.Bytes())
	}
}

func decodeAction(body io.ReadCloser) (Action, error) {
	var action Action
	logger := log.WithField("action", "DecodeAction").Logger
	logger.Debug("decode json")
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&action)
	return action, err
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := log.WithField("action", "Index").Logger
	logger.Debug("Read index")
	f, err := statikFS.Open("/index.html")
	if err != nil {
		SendErrorResponse(w, ErrorResponse{Action: "Index", Context: "open index file", Message: err.Error()})
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		SendErrorResponse(w, ErrorResponse{Action: "Index", Context: "stat index file", Message: err.Error()})
		return
	}
	buffer := make([]byte, fi.Size())
	n, err := f.Read(buffer)
	if err != nil {
		SendErrorResponse(w, ErrorResponse{Action: "Index", Context: "read index file", Message: err.Error()})
		return
	}
	if n != int(fi.Size()) {
		SendErrorResponse(w, ErrorResponse{Action: "Index", Context: "read index file", Message: "read " + strconv.Itoa(n) + " instead of " + strconv.Itoa(int(fi.Size()))})
		return
	}
	logger.Debug("Send index")
	w.Write(buffer)
}

func DriveForward(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	logger := log.WithField("action", "DriveForward").Logger
	action, err := decodeAction(r.Body)
	if err != nil {
		SendErrorResponse(w, ErrorResponse{Action: "DriveForward", Context: "json decoding", Message: err.Error()})
		return
	}
	logger.Debug("Forward 2+3, A 3")
	Motor.ArrowOn("MOTOR3")
	Motor.Forward("MOTOR2")
	Motor.Forward("MOTOR3")
	logger.Debug("Sleep " + strconv.Itoa(action.Duration))
	if action.Duration >= 0 {
		time.Sleep(time.Second * time.Duration(action.Duration))
		logger.Debug("Off 2+3, A 3")
		Motor.Stop("MOTOR2")
		Motor.Stop("MOTOR3")
		Motor.ArrowOff("MOTOR3")
	}
	w.WriteHeader(http.StatusNoContent)
}

func DriveReverse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	logger := log.WithField("action", "DriveReverse").Logger
	action, err := decodeAction(r.Body)
	if err != nil {
		SendErrorResponse(w, ErrorResponse{Action: "DriveReverse", Context: "json decoding", Message: err.Error()})
		return
	}
	logger.Debug("Reverse 2+3, A 3")
	Motor.ArrowOn("MOTOR1")
	Motor.Reverse("MOTOR2")
	Motor.Reverse("MOTOR3")
	logger.Debug("Sleep " + strconv.Itoa(action.Duration))
	if action.Duration >= 0 {
		time.Sleep(time.Second * time.Duration(action.Duration))
		logger.Debug("Off 2+3, A 1")
		Motor.Stop("MOTOR2")
		Motor.Stop("MOTOR3")
		Motor.ArrowOff("MOTOR1")
	}
	w.WriteHeader(http.StatusNoContent)
}

func TurnLeft(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	logger := log.WithField("action", "TurnLeft").Logger
	action, err := decodeAction(r.Body)
	if err != nil {
		SendErrorResponse(w, ErrorResponse{Action: "TurnLeft", Context: "json decoding", Message: err.Error()})
		return
	}
	logger.Debug("Forward 2, Reverse 3, A 4")
	Motor.ArrowOn("MOTOR2")
	Motor.Reverse("MOTOR2")
	Motor.Forward("MOTOR3")
	logger.Debug("Sleep " + strconv.Itoa(action.Duration))
	if action.Duration >= 0 {
		time.Sleep(time.Second * time.Duration(action.Duration))
		logger.Debug("Off 2+3, A 4")
		Motor.Stop("MOTOR2")
		Motor.Stop("MOTOR3")
		Motor.ArrowOff("MOTOR2")
	}
	w.WriteHeader(http.StatusNoContent)
}

func TurnRight(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	logger := log.WithField("action", "TurnRight").Logger
	action, err := decodeAction(r.Body)
	if err != nil {
		SendErrorResponse(w, ErrorResponse{Action: "TurnRight", Context: "json decoding", Message: err.Error()})
		return
	}
	logger.Debug("Forward 2, Reverse 3, A 2")
	Motor.ArrowOn("MOTOR2")
	Motor.Forward("MOTOR2")
	Motor.Reverse("MOTOR3")
	logger.Debug("Sleep " + strconv.Itoa(action.Duration))
	if action.Duration >= 0 {
		time.Sleep(time.Second * time.Duration(action.Duration))
		logger.Debug("Off 2+3, A 2")
		Motor.Stop("MOTOR2")
		Motor.Stop("MOTOR3")
		Motor.ArrowOff("MOTOR2")
	}
	w.WriteHeader(http.StatusNoContent)
}

func CameraUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	logger := log.WithField("action", "CameraUp").Logger
	action, err := decodeAction(r.Body)
	if err != nil {
		SendErrorResponse(w, ErrorResponse{Action: "CameraUp", Context: "json decoding", Message: err.Error()})
		return
	}
	logger.Debug("Forward 1")
	Motor.Forward("MOTOR1")
	logger.Debug("Sleep " + strconv.Itoa(action.Duration))
	if action.Duration >= 0 {
		time.Sleep(time.Second * time.Duration(action.Duration))
		logger.Debug("Off 1")
		Motor.Stop("MOTOR1")
	}
	w.WriteHeader(http.StatusNoContent)
}

func CameraDown(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	logger := log.WithField("action", "CameraDown").Logger
	action, err := decodeAction(r.Body)
	if err != nil {
		SendErrorResponse(w, ErrorResponse{Action: "CameraDown", Context: "json decoding", Message: err.Error()})
		return
	}
	logger.Debug("Reverse 1")
	Motor.Reverse("MOTOR1")
	logger.Debug("Sleep " + strconv.Itoa(action.Duration))
	if action.Duration >= 0 {
		time.Sleep(time.Second * time.Duration(action.Duration))
		logger.Debug("Off 1")
		Motor.Stop("MOTOR1")
	}
	w.WriteHeader(http.StatusNoContent)
}

func AllStop(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	logger := log.WithField("action", "AllStop").Logger
	logger.Debug("All Stop")
	for _, m := range motor.Motors {
		logger.Debug("Motor+Arrow Stop " + m)
		Motor.Stop(m)
		Motor.ArrowOff(m)
	}
	w.WriteHeader(http.StatusNoContent)
}
