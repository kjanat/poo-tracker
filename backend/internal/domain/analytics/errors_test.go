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
