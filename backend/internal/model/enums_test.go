package model

import (
	"testing"
)

func TestVolumeEnum(t *testing.T) {
	tests := []struct {
		name    string
		volume  Volume
		isValid bool
	}{
		{"Valid SMALL", VolumeSmall, true},
		{"Valid MEDIUM", VolumeMedium, true},
		{"Valid LARGE", VolumeLarge, true},
		{"Valid MASSIVE", VolumeMassive, true},
		{"Invalid empty", Volume(""), false},
		{"Invalid random", Volume("RANDOM"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.volume.IsValid(); got != tt.isValid {
				t.Errorf("Volume.IsValid() = %v, want %v", got, tt.isValid)
			}
		})
	}
}

func TestColorEnum(t *testing.T) {
	tests := []struct {
		name    string
		color   Color
		isValid bool
	}{
		{"Valid BROWN", ColorBrown, true},
		{"Valid DARK_BROWN", ColorDarkBrown, true},
		{"Valid LIGHT_BROWN", ColorLightBrown, true},
		{"Valid YELLOW", ColorYellow, true},
		{"Valid GREEN", ColorGreen, true},
		{"Valid RED", ColorRed, true},
		{"Valid BLACK", ColorBlack, true},
		{"Invalid empty", Color(""), false},
		{"Invalid random", Color("PURPLE"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.IsValid(); got != tt.isValid {
				t.Errorf("Color.IsValid() = %v, want %v", got, tt.isValid)
			}
		})
	}
}

func TestConsistencyEnum(t *testing.T) {
	tests := []struct {
		name        string
		consistency Consistency
		isValid     bool
	}{
		{"Valid SOLID", ConsistencySolid, true},
		{"Valid SOFT", ConsistencySoft, true},
		{"Valid LOOSE", ConsistencyLoose, true},
		{"Valid WATERY", ConsistencyWatery, true},
		{"Invalid empty", Consistency(""), false},
		{"Invalid random", Consistency("MUSHY"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.consistency.IsValid(); got != tt.isValid {
				t.Errorf("Consistency.IsValid() = %v, want %v", got, tt.isValid)
			}
		})
	}
}

func TestSmellLevelEnum(t *testing.T) {
	tests := []struct {
		name       string
		smellLevel SmellLevel
		isValid    bool
	}{
		{"Valid NONE", SmellNone, true},
		{"Valid MILD", SmellMild, true},
		{"Valid MODERATE", SmellModerate, true},
		{"Valid STRONG", SmellStrong, true},
		{"Valid TOXIC", SmellToxic, true},
		{"Invalid empty", SmellLevel(""), false},
		{"Invalid random", SmellLevel("HORRIBLE"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.smellLevel.IsValid(); got != tt.isValid {
				t.Errorf("SmellLevel.IsValid() = %v, want %v", got, tt.isValid)
			}
		})
	}
}

func TestMealCategoryEnum(t *testing.T) {
	tests := []struct {
		name     string
		category MealCategory
		isValid  bool
	}{
		{"Valid BREAKFAST", MealBreakfast, true},
		{"Valid LUNCH", MealLunch, true},
		{"Valid DINNER", MealDinner, true},
		{"Valid SNACK", MealSnack, true},
		{"Invalid empty", MealCategory(""), false},
		{"Invalid random", MealCategory("DESSERT"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.IsValid(); got != tt.isValid {
				t.Errorf("MealCategory.IsValid() = %v, want %v", got, tt.isValid)
			}
		})
	}
}

func TestSymptomTypeEnum(t *testing.T) {
	tests := []struct {
		name        string
		symptomType SymptomType
		isValid     bool
	}{
		{"Valid BLOATING", SymptomBloating, true},
		{"Valid CRAMPS", SymptomCramps, true},
		{"Valid NAUSEA", SymptomNausea, true},
		{"Valid HEARTBURN", SymptomHeartburn, true},
		{"Valid CONSTIPATION", SymptomConstipation, true},
		{"Valid DIARRHEA", SymptomDiarrhea, true},
		{"Valid GAS", SymptomGas, true},
		{"Valid FATIGUE", SymptomFatigue, true},
		{"Valid OTHER", SymptomOther, true},
		{"Invalid empty", SymptomType(""), false},
		{"Invalid random", SymptomType("HEADACHE"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.symptomType.IsValid(); got != tt.isValid {
				t.Errorf("SymptomType.IsValid() = %v, want %v", got, tt.isValid)
			}
		})
	}
}

func TestAuditActionEnum(t *testing.T) {
	tests := []struct {
		name    string
		action  AuditAction
		isValid bool
	}{
		{"Valid CREATE", AuditCreate, true},
		{"Valid UPDATE", AuditUpdate, true},
		{"Valid DELETE", AuditDelete, true},
		{"Invalid empty", AuditAction(""), false},
		{"Invalid random", AuditAction("DESTROY"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.action.IsValid(); got != tt.isValid {
				t.Errorf("AuditAction.IsValid() = %v, want %v", got, tt.isValid)
			}
		})
	}
}

func TestParseEnum(t *testing.T) {
	t.Run("ParseVolume", func(t *testing.T) {
		tests := []struct {
			input   string
			want    Volume
			wantErr bool
		}{
			{"SMALL", VolumeSmall, false},
			{"small", VolumeSmall, false},
			{"  MEDIUM  ", VolumeMedium, false},
			{"INVALID", Volume(""), true},
			{"", Volume(""), true},
		}

		for _, tt := range tests {
			got, err := ParseVolume(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVolume(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				continue
			}
			if got != tt.want {
				t.Errorf("ParseVolume(%q) = %v, want %v", tt.input, got, tt.want)
			}
		}
	})

	t.Run("ParseColor", func(t *testing.T) {
		tests := []struct {
			input   string
			want    Color
			wantErr bool
		}{
			{"BROWN", ColorBrown, false},
			{"brown", ColorBrown, false},
			{"  DARK_BROWN  ", ColorDarkBrown, false},
			{"INVALID", Color(""), true},
		}

		for _, tt := range tests {
			got, err := ParseColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				continue
			}
			if got != tt.want {
				t.Errorf("ParseColor(%q) = %v, want %v", tt.input, got, tt.want)
			}
		}
	})
}

func TestEnumValue(t *testing.T) {
	t.Run("Volume Value", func(t *testing.T) {
		v := VolumeSmall
		val, err := v.Value()
		if err != nil {
			t.Errorf("Volume.Value() error = %v", err)
		}
		if val != "SMALL" {
			t.Errorf("Volume.Value() = %v, want SMALL", val)
		}
	})

	t.Run("Invalid Volume Value", func(t *testing.T) {
		v := Volume("INVALID")
		_, err := v.Value()
		if err == nil {
			t.Error("Volume.Value() should return error for invalid volume")
		}
	})
}

func TestEnumString(t *testing.T) {
	tests := []struct {
		name string
		enum interface{ String() string }
		want string
	}{
		{"Volume", VolumeSmall, "SMALL"},
		{"Color", ColorBrown, "BROWN"},
		{"Consistency", ConsistencySoft, "SOFT"},
		{"SmellLevel", SmellMild, "MILD"},
		{"MealCategory", MealBreakfast, "BREAKFAST"},
		{"SymptomType", SymptomBloating, "BLOATING"},
		{"AuditAction", AuditCreate, "CREATE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enum.String(); got != tt.want {
				t.Errorf("%s.String() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestAllEnumValues(t *testing.T) {
	t.Run("AllVolumes", func(t *testing.T) {
		volumes := AllVolumes()
		expected := []Volume{VolumeSmall, VolumeMedium, VolumeLarge, VolumeMassive}
		if len(volumes) != len(expected) {
			t.Errorf("AllVolumes() length = %d, want %d", len(volumes), len(expected))
		}
		for i, v := range expected {
			if volumes[i] != v {
				t.Errorf("AllVolumes()[%d] = %v, want %v", i, volumes[i], v)
			}
		}
	})

	t.Run("AllColors", func(t *testing.T) {
		colors := AllColors()
		if len(colors) != 7 {
			t.Errorf("AllColors() length = %d, want 7", len(colors))
		}
		// Check that all expected colors are present
		colorMap := make(map[Color]bool)
		for _, c := range colors {
			colorMap[c] = true
		}
		expectedColors := []Color{ColorBrown, ColorDarkBrown, ColorLightBrown, ColorYellow, ColorGreen, ColorRed, ColorBlack}
		for _, c := range expectedColors {
			if !colorMap[c] {
				t.Errorf("AllColors() missing color: %v", c)
			}
		}
	})
}
