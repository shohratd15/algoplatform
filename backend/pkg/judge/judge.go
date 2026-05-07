package judge

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Константы для статусов (оставим для совместимости)
const (
	StatusInQueue       = 1
	StatusProcessing    = 2
	StatusAccepted      = 3
	StatusWrongAnswer   = 4
	StatusTimeLimit     = 5
	StatusCompileError  = 6
	StatusInternalError = 13
)

type Client struct {
	results sync.Map // map[string]*ResultResponse
}

func NewClient() *Client {
	return &Client{}
}

type SubmissionRequest struct {
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
	Stdin      string `json:"stdin"`
	Expected   string `json:"expected_output"`
}

type SubmissionResponse struct {
	Token string `json:"token"`
}

type ResultResponse struct {
	Status struct {
		ID   int    `json:"id"`
		Name string `json:"description"`
	} `json:"status"`
	Stdout         string `json:"stdout"`
	ExpectedOutput string `json:"expected_output"`
	CompileOutput  string `json:"compile_output"`
	Stderr         string `json:"stderr"`
	Time           string `json:"time"`
	Memory         int    `json:"memory"`
	Message        string `json:"message"`
}

func (c *Client) Submit(ctx context.Context, req SubmissionRequest) (string, error) {
	token := generateToken()
	result := c.executeCode(req)
	c.results.Store(token, result)
	return token, nil
}

func (c *Client) GetResult(ctx context.Context, token string) (*ResultResponse, error) {
	if val, ok := c.results.Load(token); ok {
		return val.(*ResultResponse), nil
	}
	return nil, fmt.Errorf("result not found for token %s", token)
}

func generateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (c *Client) executeCode(req SubmissionRequest) *ResultResponse {
	result := &ResultResponse{
		Status: struct {
			ID   int    `json:"id"`
			Name string `json:"description"`
		}{ID: StatusProcessing, Name: "Processing"},
	}

	result.ExpectedOutput = req.Expected

	// Определяем язык
	var lang string
	switch req.LanguageID {
	case 71: // Python
		lang = "python"
	case 60: // Go
		lang = "go"
	case 63: // JavaScript
		lang = "js"
	case 54: // C++
		lang = "cpp"
	case 62: // Java
		lang = "java"
	default:
		result.Status.ID = StatusInternalError
		result.Status.Name = "Unsupported language"
		result.Message = "Unsupported language ID"
		return result
	}

	// Создаем временную директорию
	tmpDir, err := ioutil.TempDir("", "judge")
	if err != nil {
		result.Status.ID = StatusInternalError
		result.Status.Name = "Internal error"
		result.Message = "Failed to create temp dir"
		return result
	}
	defer os.RemoveAll(tmpDir)

	// Пишем код в файл
	codeFile := filepath.Join(tmpDir, "code")
	switch lang {
	case "python":
		codeFile += ".py"
	case "go":
		codeFile += ".go"
	case "js":
		codeFile += ".js"
	case "cpp":
		codeFile += ".cpp"
	case "java":
		codeFile += ".java"
	}

	if err := ioutil.WriteFile(codeFile, []byte(req.SourceCode), 0644); err != nil {
		result.Status.ID = StatusInternalError
		result.Message = "Failed to write code file"
		return result
	}

	// Пишем input в файл
	inputFile := filepath.Join(tmpDir, "input.txt")
	if err := ioutil.WriteFile(inputFile, []byte(req.Stdin), 0644); err != nil {
		result.Status.ID = StatusInternalError
		result.Message = "Failed to write input file"
		return result
	}

	// Компилируем если нужно
	var execCmd string
	var compileCmd []string
	switch lang {
	case "python":
		execCmd = fmt.Sprintf("python3 %s", codeFile)
	case "go":
		execCmd = fmt.Sprintf("go run %s", codeFile)
	case "js":
		execCmd = fmt.Sprintf("node %s", codeFile)
	case "cpp":
		compileCmd = []string{"g++", "-o", filepath.Join(tmpDir, "program"), codeFile}
		execCmd = filepath.Join(tmpDir, "program")
	case "java":
		// Для Java, код должен быть в классе, но для простоты предполагаем main
		compileCmd = []string{"javac", codeFile}
		execCmd = "java -cp " + tmpDir + " Main"
	}

	// Компиляция
	if len(compileCmd) > 0 {
		cmd := exec.Command(compileCmd[0], compileCmd[1:]...)
		cmd.Dir = tmpDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			result.Status.ID = StatusCompileError
			result.Status.Name = "Compilation error"
			result.CompileOutput = string(output)
			result.Message = err.Error()
			return result
		}
	}

	// Запуск напрямую
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("%s < %s", execCmd, inputFile))
	cmd.Dir = tmpDir

	start := time.Now()
	output, err := cmd.CombinedOutput()
	duration := time.Since(start)

	result.Time = fmt.Sprintf("%.3f", duration.Seconds())

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 137 { // SIGKILL, вероятно time limit
				result.Status.ID = StatusTimeLimit
				result.Status.Name = "Time limit exceeded"
			} else {
				result.Status.ID = StatusInternalError
				result.Status.Name = "Runtime error"
				result.Stderr = string(output)
				result.Message = err.Error()
			}
		} else {
			result.Status.ID = StatusInternalError
			result.Stderr = string(output)
			result.Message = err.Error()
		}
		return result
	}

	result.Stdout = strings.TrimSpace(string(output))

	// Сравниваем с expected
	if strings.TrimSpace(result.Stdout) == strings.TrimSpace(req.Expected) {
		result.Status.ID = StatusAccepted
		result.Status.Name = "Accepted"
	} else {
		result.Status.ID = StatusWrongAnswer
		result.Status.Name = "Wrong answer"
	}

	return result
}
