package main

// A clever way of implementing error constants,
// courtesy of Dave Cheney: https://dave.cheney.net/2016/04/07/constant-errors

type ErrorWithLevel interface {
	IsCritical() bool
}

type Error string

func (e Error) Error() string    { return string(e) }
func (e Error) IsCritical() bool { return true }

type Warning string

func (e Warning) Error() string    { return string(e) }
func (e Warning) IsCritical() bool { return false }
