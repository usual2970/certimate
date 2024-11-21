package deployer_test

import (
	"testing"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
)

/*
Shell command to run this test:

	go test -v logger_test.go
*/
func TestLogger(t *testing.T) {
	t.Run("Logger_Appendt", func(t *testing.T) {
		logger := deployer.NewDefaultLogger()

		logger.Appendt("test")
		logger.Appendt("test_nil", nil)
		logger.Appendt("test_int", 1024)
		logger.Appendt("test_string", "certimate")
		logger.Appendt("test_map", map[string]interface{}{"key": "value"})
		logger.Appendt("test_struct", struct{ Name string }{Name: "certimate"})
		logger.Appendt("test_slice", []string{"certimate"})
		t.Log(logger.GetRecords())
		if len(logger.GetRecords()) != 7 {
			t.Errorf("expected 7 records, got %d", len(logger.GetRecords()))
		}

		logger.Flush()
		if len(logger.GetRecords()) != 0 {
			t.Errorf("expected 0 records, got %d", len(logger.GetRecords()))
		}
	})

	t.Run("Logger_Appendf", func(t *testing.T) {
		logger := deployer.NewDefaultLogger()

		logger.Appendf("test")
		logger.Appendf("test_nil: %v", nil)
		logger.Appendf("test_int: %v", 1024)
		logger.Appendf("test_string: %v", "certimate")
		logger.Appendf("test_map: %v", map[string]interface{}{"key": "value"})
		logger.Appendf("test_struct: %v", struct{ Name string }{Name: "certimate"})
		logger.Appendf("test_slice: %v", []string{"certimate"})
		t.Log(logger.GetRecords())
		if len(logger.GetRecords()) != 7 {
			t.Errorf("expected 7 records, got %d", len(logger.GetRecords()))
		}

		logger.Flush()
		if len(logger.GetRecords()) != 0 {
			t.Errorf("expected 0 records, got %d", len(logger.GetRecords()))
		}
	})
}
