package controllers

import (
	"encoding/json"
	"errors"
	"os"
)

type Stack struct {
	Stack []string `json:"stack"`
}

func (c *Stack) Push(checksum string) {
	c.Stack = append(c.Stack, checksum)
}

func (c *Stack) Pop() (string, error) {
	if c.isEmpty() {
		return "", errors.New("stack is empty; cannot pop")
	}
	top := c.Stack[len(c.Stack)-1]
	c.Stack = c.Stack[:len(c.Stack)-1]
	return top, nil
}

func (c *Stack) isEmpty() bool {
	return len(c.Stack) == 0
}

func (c *Stack) WriteToJSON(filePath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return errors.New("failed to marshal stack to JSON: " + err.Error())
	}
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return errors.New("failed to write stack to file: " + err.Error())
	}
	return nil
}

func (c *Stack) ReadFromJSON(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return errors.New("failed to read file: " + err.Error())
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return errors.New("failed to unmarshal JSON to stack: " + err.Error())
	}
	return nil
}
