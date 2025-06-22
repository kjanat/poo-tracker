package symptom

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	shared "github.com/kjanat/poo-tracker/backend/internal/domain/shared"
	domainsymptom "github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

// ToSymptomDB converts a domain Symptom to a GORM SymptomDB, marshaling Triggers as JSON.
func ToSymptomDB(s *domainsymptom.Symptom) (*domainsymptom.SymptomDB, error) {
	triggersJSON, err := json.Marshal(s.Triggers)
	if err != nil {
		return nil, err
	}
	return &domainsymptom.SymptomDB{
		ID:          uuid.MustParse(s.ID),
		UserID:      s.UserID,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		RecordedAt:  s.RecordedAt,
		Name:        s.Name,
		Description: nilIfEmpty(s.Description),
		Category:    ptrToString(s.Category),
		Severity:    s.Severity,
		Duration:    s.Duration,
		BodyPart:    nilIfEmpty(s.BodyPart),
		Type:        ptrToString(s.Type),
		Triggers:    datatypes.JSON(triggersJSON),
		Notes:       nilIfEmpty(s.Notes),
		PhotoURL:    nilIfEmpty(s.PhotoURL),
	}, nil
}

// ToSymptom converts a GORM SymptomDB to a domain Symptom, unmarshaling Triggers from JSON.
func ToSymptom(db *domainsymptom.SymptomDB) (*domainsymptom.Symptom, error) {
	var triggers []string
	if len(db.Triggers) > 0 {
		if err := json.Unmarshal(db.Triggers, &triggers); err != nil {
			return nil, err
		}
	}
	return &domainsymptom.Symptom{
		ID:          db.ID.String(),
		UserID:      db.UserID,
		CreatedAt:   db.CreatedAt,
		UpdatedAt:   db.UpdatedAt,
		RecordedAt:  db.RecordedAt,
		Name:        db.Name,
		Description: derefString(db.Description),
		Category:    stringToSymptomCategory(db.Category),
		Severity:    db.Severity,
		Duration:    db.Duration,
		BodyPart:    derefString(db.BodyPart),
		Type:        stringToSymptomType(db.Type),
		Triggers:    triggers,
		Notes:       derefString(db.Notes),
		PhotoURL:    derefString(db.PhotoURL),
	}, nil
}

// Helper functions for pointer/string conversions
func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ptrToString[T ~string](v *T) *string {
	if v == nil {
		return nil
	}
	str := string(*v)
	return &str
}

func stringToSymptomCategory(s *string) *shared.SymptomCategory {
	if s == nil {
		return nil
	}
	cat := shared.SymptomCategory(*s)
	return &cat
}

func stringToSymptomType(s *string) *shared.SymptomType {
	if s == nil {
		return nil
	}
	t := shared.SymptomType(*s)
	return &t
}
