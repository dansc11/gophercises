package main

import "io"

type gameConfig struct {
	ProblemsPath string
	Reader       io.Reader
	Writer       io.Writer
}
