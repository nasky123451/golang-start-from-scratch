// Package calculator 提供簡單的數學操作。
package calculator

import (
	"errors"
)

// Add 返回兩個整數的和。
func Add(a, b int) int {
	return a + b
}

// Subtract 返回第一個整數減去第二個整數的結果。
func Subtract(a, b int) int {
	return a - b
}

// Multiply 返回兩個整數的乘積。
func Multiply(a, b int) int {
	return a * b
}

// Divide 返回第一個整數除以第二個整數的結果。
// 如果第二個數為 0，會返回一個錯誤。
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}
