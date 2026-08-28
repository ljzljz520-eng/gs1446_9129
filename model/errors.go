package model

import "errors"

var ErrDuplicate = errors.New("duplicate record")
var ErrNotFound = errors.New("not found")
var ErrInvalidRecord = errors.New("invalid record")
var ErrInvalidProfile = errors.New("invalid profile")
var ErrArchived = errors.New("record archived")
