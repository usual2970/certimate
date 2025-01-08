package logger_test

import (
	"testing"

	"github.com/usual2970/certimate/internal/pkg/core/logger"
)

/*
Shell command to run this test:

	go test -v ./logger_test.go
*/
func TestLogger(t *testing.T) {
	t.Run("Logger_Appendt", func(t *testing.T) {
		logger := logger.NewDefaultLogger()

		logger.Logt("test")
		logger.Logt("test_nil", nil)
		logger.Logt("test_int", 1024)
		logger.Logt("test_string", "certimate")
		logger.Logt("test_map", map[string]interface{}{"key": "value"})
		logger.Logt("test_struct", struct{ Name string }{Name: "certimate"})
		logger.Logt("test_slice", []string{"certimate"})
		t.Log(logger.GetRecords())
		if len(logger.GetRecords()) != 7 {
			t.Errorf("expected 7 records, got %d", len(logger.GetRecords()))
		}

		logger.FlushRecords()
		if len(logger.GetRecords()) != 0 {
			t.Errorf("expected 0 records, got %d", len(logger.GetRecords()))
		}
	})

	t.Run("Logger_Appendf", func(t *testing.T) {
		logger := logger.NewDefaultLogger()

		logger.Logf("test")
		logger.Logf("test_nil: %v", nil)
		logger.Logf("test_int: %v", 1024)
		logger.Logf("test_string: %v", "certimate")
		logger.Logf("test_map: %v", map[string]interface{}{"key": "value"})
		logger.Logf("test_struct: %v", struct{ Name string }{Name: "certimate"})
		logger.Logf("test_slice: %v", []string{"certimate"})
		t.Log(logger.GetRecords())
		if len(logger.GetRecords()) != 7 {
			t.Errorf("expected 7 records, got %d", len(logger.GetRecords()))
		}

		logger.FlushRecords()
		if len(logger.GetRecords()) != 0 {
			t.Errorf("expected 0 records, got %d", len(logger.GetRecords()))
		}
	})
}
