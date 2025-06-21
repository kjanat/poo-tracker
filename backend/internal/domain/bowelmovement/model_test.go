package bowelmovement

import "testing"

func TestNewBowelMovement(t *testing.T) {
	tests := []struct {
		name        string
		bristolType int
		wantErr     bool
	}{
		{"valid", 3, false},
		{"invalid low", 0, true},
		{"invalid high", 8, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm, err := NewBowelMovement("user1", tt.bristolType)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for bristolType %d", tt.bristolType)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if bm.BristolType != tt.bristolType {
				t.Errorf("expected BristolType %d, got %d", tt.bristolType, bm.BristolType)
			}
		})
	}
}
