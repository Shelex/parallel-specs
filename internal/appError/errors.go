package appError

import "errors"

var AccessDenied = errors.New("access denied")                      //nolint
var InvalidEmailOrPassord = errors.New("invalid email or password") //nolint
var ProjectNotFound = errors.New("project not found")               //nolint
var SessionNotFound = errors.New("session not found")               //nolint
var SpecNotFound = errors.New("spec not found")                     //nolint
var SessionFinished = errors.New("session already finished")        //nolint
var ApiKeyNotFound = errors.New("api key not found")                //nolint
