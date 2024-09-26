package api

import "testing"

// 測試 Add 函數
func TestAdd(t *testing.T) {
	result := Add(3, 2)
	expected := 5
	if result != expected {
		t.Errorf("Add(3, 2) = %d; want %d", result, expected)
	}
}

// 測試 Subtract 函數
func TestSubtract(t *testing.T) {
	result := Subtract(5, 3)
	expected := 2
	if result != expected {
		t.Errorf("Subtract(5, 3) = %d; want %d", result, expected)
	}
}

// 測試 Multiply 函數
func TestMultiply(t *testing.T) {
	result := Multiply(4, 2)
	expected := 8
	if result != expected {
		t.Errorf("Multiply(4, 2) = %d; want %d", result, expected)
	}
}

// 測試 Divide 函數
func TestDivide(t *testing.T) {
	// 測試正常除法
	result, err := Divide(10, 2)
	expected := 5
	if err != nil || result != expected {
		t.Errorf("Divide(10, 2) = %d, %v; want %d, nil", result, err, expected)
	}

	// 測試除以零
	_, err = Divide(10, 0)
	if err == nil {
		t.Error("Expected error when dividing by zero, but got none")
	}
}
