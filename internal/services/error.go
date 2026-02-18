package services

import "errors"

var (
	ErrUserExists = errors.New("user exists")
	ErrNotFound   = errors.New("subscription not found")
)
