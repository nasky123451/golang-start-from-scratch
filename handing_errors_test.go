package main

import (
	"errors"
	"testing"
)

// 測試 noErrCanHappen 函數
func TestNoErrCanHappen(t *testing.T) {
	expected := 204
	result := noErrCanHappen()
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

// 測試 doOnErr 函數
func TestDoOnErr(t *testing.T) {
	// 模擬不失敗情況
	if err := doOnErr(func() bool { return false }); err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// 模擬失敗情況
	if err := doOnErr(func() bool { return true }); err == nil {
		t.Errorf("Expected an error, but got none")
	} else if err.Error() != "ups, XYZ failed" {
		t.Errorf("Expected error message 'ups, XYZ failed', but got %v", err)
	}
}

// 測試 intOrErr 函數
func TestIntOrErr(t *testing.T) {
	// 測試不會失敗的情況
	result, err := intOrErr(func() bool { return false })
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	expected := 204
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}

	// 測試會失敗的情況
	_, err = intOrErr(func() bool { return true })
	if err == nil {
		t.Errorf("Expected an error, but got none")
	} else if err.Error() != "ups, XYZ2 failed" {
		t.Errorf("Expected error message 'ups, XYZ2 failed', but got %v", err)
	}
}

func TestNestedDoOrErr(t *testing.T) {
	// 測試不會失敗的情況
	if err := nestedDoOrErr(func() bool { return false }); err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// 測試會失敗的情況
	if err := nestedDoOrErr(func() bool { return true }); err == nil {
		t.Errorf("Expected an error, but got none")
	} else {
		// 使用 errors.Is() 來檢查底層錯誤是否是全局錯誤 ErrXYZFailed
		if !errors.Is(err, ErrXYZFailed) {
			t.Errorf("Expected wrapped error containing '%v', but got %v", ErrXYZFailed, err)
		}
	}
}
