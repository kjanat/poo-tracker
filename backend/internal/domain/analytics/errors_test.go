package analytics

import "testing"

func TestNewAnalysisError(t *testing.T) {
	e := NewAnalysisError("TREND", "failed", "ctx")
	if e.Type != "TREND" {
		t.Errorf("expected Type TREND, got %s", e.Type)
	}
	if e.Message != "failed" {
		t.Errorf("expected Message failed, got %s", e.Message)
	}
	if e.Context != "ctx" {
		t.Errorf("expected Context ctx, got %s", e.Context)
	}
	expectedErr := "analysis error 'TREND': failed (context: ctx)"
	if e.Error() != expectedErr {
		t.Errorf("unexpected Error string: %s", e.Error())
	}
}

func TestNewAnalysisErrorEdgeCases(t *testing.T) {
	e := NewAnalysisError("", "", "")
	if e.Type != "" || e.Message != "" || e.Context != "" {
		t.Error("empty strings not handled correctly")
	}

	e2 := NewAnalysisError("TYPE/SPECIAL", "message with 'quotes'", "ctx:value")
	expectedErr := "analysis error 'TYPE/SPECIAL': message with 'quotes' (context: ctx:value)"
	if e2.Error() != expectedErr {
		t.Errorf("special characters not handled correctly: %s", e2.Error())
	}
}
