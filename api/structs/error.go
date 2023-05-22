// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
	Msg  string
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

func (se StatusError) Message() string {
	return se.Msg
}
