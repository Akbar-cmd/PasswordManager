package main

import "errors"

var ErrPassManagerNotInit = errors.New("password manager not initialized")
var ErrPassExists = errors.New("password already exists")
var ErrPassNotFound = errors.New("password not found")
var ErrPassWeak = errors.New("password is too weak")
