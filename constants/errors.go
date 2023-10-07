package errorMessages

import "errors"

var RecordDoesntExist = errors.New("record doesnt exist")
var RecordAlreadyExist = errors.New("record already exist")
