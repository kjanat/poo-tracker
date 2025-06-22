package repository

import "errors"

var ErrNotFound = errors.New("not found")
var ErrRelationAlreadyExists = errors.New("relation already exists")
