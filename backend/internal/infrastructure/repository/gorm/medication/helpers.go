package medication

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	domainmedication "github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	shared "github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// ToMedicationDB converts a domain Medication to a GORM MedicationDB, marshaling SideEffects as JSON.
func ToMedicationDB(m *domainmedication.Medication) (*domainmedication.MedicationDB, error) {
	sideEffectsJSON, err := json.Marshal(m.SideEffects)
	if err != nil {
		return nil, err
	}
	return &domainmedication.MedicationDB{
		ID:          uuid.MustParse(m.ID),
		UserID:      m.UserID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Name:        m.Name,
		GenericName: nilIfEmpty(m.GenericName),
		Brand:       nilIfEmpty(m.Brand),
		Category:    ptrToString(m.Category),
		Dosage:      m.Dosage,
		Form:        ptrToString(m.Form),
		Frequency:   m.Frequency,
		Route:       ptrToString(m.Route),
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
		TakenAt:     m.TakenAt,
		Purpose:     nilIfEmpty(m.Purpose),
		SideEffects: datatypes.JSON(sideEffectsJSON),
		Notes:       nilIfEmpty(m.Notes),
		PhotoURL:    nilIfEmpty(m.PhotoURL),
		IsActive:    m.IsActive,
		IsAsNeeded:  m.IsAsNeeded,
	}, nil
}

// ToMedication converts a GORM MedicationDB to a domain Medication, unmarshaling SideEffects from JSON.
func ToMedication(db *domainmedication.MedicationDB) (*domainmedication.Medication, error) {
	var sideEffects []string
	if len(db.SideEffects) > 0 {
		if err := json.Unmarshal(db.SideEffects, &sideEffects); err != nil {
			return nil, err
		}
	}
	return &domainmedication.Medication{
		ID:          db.ID.String(),
		UserID:      db.UserID,
		CreatedAt:   db.CreatedAt,
		UpdatedAt:   db.UpdatedAt,
		Name:        db.Name,
		GenericName: derefString(db.GenericName),
		Brand:       derefString(db.Brand),
		Category:    stringToMedicationCategory(db.Category),
		Dosage:      db.Dosage,
		Form:        stringToMedicationForm(db.Form),
		Frequency:   db.Frequency,
		Route:       stringToMedicationRoute(db.Route),
		StartDate:   db.StartDate,
		EndDate:     db.EndDate,
		TakenAt:     db.TakenAt,
		Purpose:     derefString(db.Purpose),
		SideEffects: sideEffects,
		Notes:       derefString(db.Notes),
		PhotoURL:    derefString(db.PhotoURL),
		IsActive:    db.IsActive,
		IsAsNeeded:  db.IsAsNeeded,
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

func stringToMedicationCategory(s *string) *shared.MedicationCategory {
	if s == nil {
		return nil
	}
	cat := shared.MedicationCategory(*s)
	return &cat
}

func stringToMedicationForm(s *string) *shared.MedicationForm {
	if s == nil {
		return nil
	}
	form := shared.MedicationForm(*s)
	return &form
}

func stringToMedicationRoute(s *string) *shared.MedicationRoute {
	if s == nil {
		return nil
	}
	route := shared.MedicationRoute(*s)
	return &route
}
