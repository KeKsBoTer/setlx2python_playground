package main

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

// TranspilerResponse is the result of a code transpiling
type TranspilerResponse struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func transpile(code []byte) (*TranspilerResponse, error) {
	if len(code) == 0 {
		return nil, errors.New("Empty code")
	}
	out, err := transpileCode(string(code))
	if err != nil {
		return nil, err
	}
	return &TranspilerResponse{
		Code:    string(out),
		Success: true,
	}, nil
}

const maxTranspileTime = 1 * time.Minute

func transpileCode(code string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxTranspileTime)
	defer cancel()
	cmd := exec.CommandContext(ctx, "setlx2python", "-c", code)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, errors.New("Programm exceeded maximum execution time")
	}
	return out, nil
}
