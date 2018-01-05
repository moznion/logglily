package logger

import "testing"

func TestValidateBulkByteSizeThresholdShouldBeSuccessfully(t *testing.T) {
	var err error

	err = validateBulkByteSizeThreshold(5 * 1024 * 1024)
	if err != nil {
		t.Error("err should be nil, but it exists", err)
	}

	err = validateBulkByteSizeThreshold(1)
	if err != nil {
		t.Error("err should be nil, but it exists", err)
	}
}

func TestValidateBulkByteSizeThresholdShouldBeFail(t *testing.T) {
	var err error

	err = validateBulkByteSizeThreshold(5*1024*1024 + 1)
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}

	err = validateBulkByteSizeThreshold(0)
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}
}
