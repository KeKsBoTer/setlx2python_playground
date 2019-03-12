package setlx2python_playground

import (
	"bytes"
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
	stdOut, errOut, err := transpileCode(string(code))
	if err != nil {
		return nil, err
	}
	success := stdOut != nil
	return &TranspilerResponse{
		Code:    string(stdOut),
		Success: success,
		Error:   string(errOut),
	}, nil
}

const maxTranspileTime = 1 * time.Minute

func transpileCode(code string) ([]byte, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxTranspileTime)
	defer cancel()
	cmd := exec.CommandContext(ctx, "setlx2python", "-c", code)

	var bStd, bErr bytes.Buffer

	cmd.Stdout = &bStd
	cmd.Stderr = &bErr

	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// if exit code != 0 return error output
			return nil, bErr.Bytes(), nil
		}
		return nil, nil, err
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, nil, errors.New("Programm exceeded maximum execution time")
	}
	return bStd.Bytes(), bErr.Bytes(), nil
}
