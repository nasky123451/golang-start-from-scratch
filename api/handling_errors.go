package api

import (
	"errors"
	"fmt"
)

var ErrXYZFailed = errors.New("ups, XYZ failed")   // 定義全局錯誤
var ErrXYZFailed2 = errors.New("ups, XYZ2 failed") // 定義全局錯誤

func noErrCanHappen() int {
	return 204
}

func doOnErr(shouldFail func() bool) error {
	if shouldFail() {
		return ErrXYZFailed
	}
	return nil
}

func intOrErr(shouldFail func() bool) (int, error) {
	if shouldFail() {
		return 0, ErrXYZFailed2
	}
	return noErrCanHappen(), nil
}

func nestedDoOrErr(shouldFail func() bool) error {
	if err := doOnErr(shouldFail); err != nil {
		return fmt.Errorf("od: %w", err)
	}
	return nil
}
