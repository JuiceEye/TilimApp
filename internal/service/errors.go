package service

import (
	"fmt"
)

var fieldGenders = map[string]string{
	"логин":            "m",
	"почта":            "f",
	"имя пользователя": "n",
}

type ErrTaken struct {
	Field string
}

func (e *ErrTaken) Error() string {
	gender, ok := fieldGenders[e.Field]
	if !ok {
		gender = "m"
	}

	var suffix string
	switch gender {
	case "m":
		suffix = "занят"
	case "f":
		suffix = "занята"
	case "n":
		suffix = "занято"
	}

	return fmt.Sprintf("%s уже %s", e.Field, suffix)
}

type BadRequestError struct {
	Msg string
}

func (e *BadRequestError) Error() string {
	return e.Msg
}
