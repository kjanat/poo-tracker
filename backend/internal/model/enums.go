package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// Volume represents the size/volume of a bowel movement
type Volume string

const (
	VolumeSmall   Volume = "SMALL"
	VolumeMedium  Volume = "MEDIUM"
	VolumeLarge   Volume = "LARGE"
	VolumeMassive Volume = "MASSIVE"
)

// AllVolumes returns all valid Volume values
func AllVolumes() []Volume {
	return []Volume{VolumeSmall, VolumeMedium, VolumeLarge, VolumeMassive}
}

// IsValid checks if the Volume value is valid
func (v Volume) IsValid() bool {
	for _, valid := range AllVolumes() {
		if v == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (v Volume) String() string {
	return string(v)
}

// Value implements the driver.Valuer interface for database storage
func (v Volume) Value() (driver.Value, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("invalid volume: %s", v)
	}
	return string(v), nil
}

// Color represents the color of a bowel movement
type Color string

const (
	ColorBrown      Color = "BROWN"
	ColorDarkBrown  Color = "DARK_BROWN"
	ColorLightBrown Color = "LIGHT_BROWN"
	ColorYellow     Color = "YELLOW"
	ColorGreen      Color = "GREEN"
	ColorRed        Color = "RED"
	ColorBlack      Color = "BLACK"
)

// AllColors returns all valid Color values
func AllColors() []Color {
	return []Color{
		ColorBrown, ColorDarkBrown, ColorLightBrown,
		ColorYellow, ColorGreen, ColorRed, ColorBlack,
	}
}

// IsValid checks if the Color value is valid
func (c Color) IsValid() bool {
	for _, valid := range AllColors() {
		if c == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (c Color) String() string {
	return string(c)
}

// Value implements the driver.Valuer interface for database storage
func (c Color) Value() (driver.Value, error) {
	if !c.IsValid() {
		return nil, fmt.Errorf("invalid color: %s", c)
	}
	return string(c), nil
}

// Consistency represents the consistency of a bowel movement
type Consistency string

const (
	ConsistencySolid  Consistency = "SOLID"
	ConsistencySoft   Consistency = "SOFT"
	ConsistencyLoose  Consistency = "LOOSE"
	ConsistencyWatery Consistency = "WATERY"
)

// AllConsistencies returns all valid Consistency values
func AllConsistencies() []Consistency {
	return []Consistency{ConsistencySolid, ConsistencySoft, ConsistencyLoose, ConsistencyWatery}
}

// IsValid checks if the Consistency value is valid
func (c Consistency) IsValid() bool {
	for _, valid := range AllConsistencies() {
		if c == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (c Consistency) String() string {
	return string(c)
}

// Value implements the driver.Valuer interface for database storage
func (c Consistency) Value() (driver.Value, error) {
	if !c.IsValid() {
		return nil, fmt.Errorf("invalid consistency: %s", c)
	}
	return string(c), nil
}

// SmellLevel represents the intensity of smell
type SmellLevel string

const (
	SmellNone     SmellLevel = "NONE"
	SmellMild     SmellLevel = "MILD"
	SmellModerate SmellLevel = "MODERATE"
	SmellStrong   SmellLevel = "STRONG"
	SmellToxic    SmellLevel = "TOXIC"
)

// AllSmellLevels returns all valid SmellLevel values
func AllSmellLevels() []SmellLevel {
	return []SmellLevel{SmellNone, SmellMild, SmellModerate, SmellStrong, SmellToxic}
}

// IsValid checks if the SmellLevel value is valid
func (s SmellLevel) IsValid() bool {
	for _, valid := range AllSmellLevels() {
		if s == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (s SmellLevel) String() string {
	return string(s)
}

// Value implements the driver.Valuer interface for database storage
func (s SmellLevel) Value() (driver.Value, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid smell level: %s", s)
	}
	return string(s), nil
}

// MealCategory represents the type/category of meal
type MealCategory string

const (
	MealBreakfast MealCategory = "BREAKFAST"
	MealLunch     MealCategory = "LUNCH"
	MealDinner    MealCategory = "DINNER"
	MealSnack     MealCategory = "SNACK"
)

// AllMealCategories returns all valid MealCategory values
func AllMealCategories() []MealCategory {
	return []MealCategory{MealBreakfast, MealLunch, MealDinner, MealSnack}
}

// IsValid checks if the MealCategory value is valid
func (m MealCategory) IsValid() bool {
	for _, valid := range AllMealCategories() {
		if m == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (m MealCategory) String() string {
	return string(m)
}

// Value implements the driver.Valuer interface for database storage
func (m MealCategory) Value() (driver.Value, error) {
	if !m.IsValid() {
		return nil, fmt.Errorf("invalid meal category: %s", m)
	}
	return string(m), nil
}

// SymptomCategory represents categories for symptoms
type SymptomCategory string

const (
	SymptomCategoryDigestive    SymptomCategory = "DIGESTIVE"
	SymptomCategoryAbdominal    SymptomCategory = "ABDOMINAL"
	SymptomCategorySystemic     SymptomCategory = "SYSTEMIC"
	SymptomCategoryNeurological SymptomCategory = "NEUROLOGICAL"
	SymptomCategoryOther        SymptomCategory = "OTHER"
)

// AllSymptomCategories returns all valid SymptomCategory values
func AllSymptomCategories() []SymptomCategory {
	return []SymptomCategory{
		SymptomCategoryDigestive,
		SymptomCategoryAbdominal,
		SymptomCategorySystemic,
		SymptomCategoryNeurological,
		SymptomCategoryOther,
	}
}

// IsValid checks if the SymptomCategory value is valid
func (sc SymptomCategory) IsValid() bool {
	for _, valid := range AllSymptomCategories() {
		if sc == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (sc SymptomCategory) String() string {
	return string(sc)
}

// Value implements the driver.Valuer interface for database storage
func (sc SymptomCategory) Value() (driver.Value, error) {
	if !sc.IsValid() {
		return nil, fmt.Errorf("invalid symptom category: %s", sc)
	}
	return string(sc), nil
}

// SymptomType represents different types of symptoms
type SymptomType string

const (
	SymptomBloating     SymptomType = "BLOATING"
	SymptomCramps       SymptomType = "CRAMPS"
	SymptomNausea       SymptomType = "NAUSEA"
	SymptomHeartburn    SymptomType = "HEARTBURN"
	SymptomConstipation SymptomType = "CONSTIPATION"
	SymptomDiarrhea     SymptomType = "DIARRHEA"
	SymptomGas          SymptomType = "GAS"
	SymptomFatigue      SymptomType = "FATIGUE"
	SymptomOther        SymptomType = "OTHER"
)

// AllSymptomTypes returns all valid SymptomType values
func AllSymptomTypes() []SymptomType {
	return []SymptomType{
		SymptomBloating, SymptomCramps, SymptomNausea, SymptomHeartburn,
		SymptomConstipation, SymptomDiarrhea, SymptomGas, SymptomFatigue, SymptomOther,
	}
}

// IsValid checks if the SymptomType value is valid
func (s SymptomType) IsValid() bool {
	for _, valid := range AllSymptomTypes() {
		if s == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (s SymptomType) String() string {
	return string(s)
}

// Value implements the driver.Valuer interface for database storage
func (s SymptomType) Value() (driver.Value, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid symptom type: %s", s)
	}
	return string(s), nil
}

// MedicationCategory represents categories for medications
type MedicationCategory string

const (
	MedicationCategoryGastrointestinal MedicationCategory = "GASTROINTESTINAL"
	MedicationCategoryPainRelief       MedicationCategory = "PAIN_RELIEF"
	MedicationCategoryAntibiotic       MedicationCategory = "ANTIBIOTIC"
	MedicationCategoryProbiotics       MedicationCategory = "PROBIOTICS"
	MedicationCategorySupplements      MedicationCategory = "SUPPLEMENTS"
	MedicationCategoryAntiInflammatory MedicationCategory = "ANTI_INFLAMMATORY"
	MedicationCategoryOtherMedication  MedicationCategory = "OTHER"
)

// AllMedicationCategories returns all valid MedicationCategory values
func AllMedicationCategories() []MedicationCategory {
	return []MedicationCategory{
		MedicationCategoryGastrointestinal,
		MedicationCategoryPainRelief,
		MedicationCategoryAntibiotic,
		MedicationCategoryProbiotics,
		MedicationCategorySupplements,
		MedicationCategoryAntiInflammatory,
		MedicationCategoryOtherMedication,
	}
}

// IsValid checks if the MedicationCategory value is valid
func (mc MedicationCategory) IsValid() bool {
	for _, valid := range AllMedicationCategories() {
		if mc == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (mc MedicationCategory) String() string {
	return string(mc)
}

// Value implements the driver.Valuer interface for database storage
func (mc MedicationCategory) Value() (driver.Value, error) {
	if !mc.IsValid() {
		return nil, fmt.Errorf("invalid medication category: %s", mc)
	}
	return string(mc), nil
}

// MedicationForm represents forms of medications
type MedicationForm string

const (
	MedicationFormTablet      MedicationForm = "TABLET"
	MedicationFormCapsule     MedicationForm = "CAPSULE"
	MedicationFormLiquid      MedicationForm = "LIQUID"
	MedicationFormCream       MedicationForm = "CREAM"
	MedicationFormPowder      MedicationForm = "POWDER"
	MedicationFormInjection   MedicationForm = "INJECTION"
	MedicationFormSuppository MedicationForm = "SUPPOSITORY"
	MedicationFormOtherForm   MedicationForm = "OTHER"
)

// AllMedicationForms returns all valid MedicationForm values
func AllMedicationForms() []MedicationForm {
	return []MedicationForm{
		MedicationFormTablet,
		MedicationFormCapsule,
		MedicationFormLiquid,
		MedicationFormCream,
		MedicationFormPowder,
		MedicationFormInjection,
		MedicationFormSuppository,
		MedicationFormOtherForm,
	}
}

// IsValid checks if the MedicationForm value is valid
func (mf MedicationForm) IsValid() bool {
	for _, valid := range AllMedicationForms() {
		if mf == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (mf MedicationForm) String() string {
	return string(mf)
}

// Value implements the driver.Valuer interface for database storage
func (mf MedicationForm) Value() (driver.Value, error) {
	if !mf.IsValid() {
		return nil, fmt.Errorf("invalid medication form: %s", mf)
	}
	return string(mf), nil
}

// MedicationRoute represents routes of administration for medications
type MedicationRoute string

const (
	MedicationRouteOral       MedicationRoute = "ORAL"
	MedicationRouteTopical    MedicationRoute = "TOPICAL"
	MedicationRouteRectal     MedicationRoute = "RECTAL"
	MedicationRouteInjection  MedicationRoute = "INJECTION"
	MedicationRouteInhalation MedicationRoute = "INHALATION"
	MedicationRouteOtherRoute MedicationRoute = "OTHER"
)

// AllMedicationRoutes returns all valid MedicationRoute values
func AllMedicationRoutes() []MedicationRoute {
	return []MedicationRoute{
		MedicationRouteOral,
		MedicationRouteTopical,
		MedicationRouteRectal,
		MedicationRouteInjection,
		MedicationRouteInhalation,
		MedicationRouteOtherRoute,
	}
}

// IsValid checks if the MedicationRoute value is valid
func (mr MedicationRoute) IsValid() bool {
	for _, valid := range AllMedicationRoutes() {
		if mr == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (mr MedicationRoute) String() string {
	return string(mr)
}

// Value implements the driver.Valuer interface for database storage
func (mr MedicationRoute) Value() (driver.Value, error) {
	if !mr.IsValid() {
		return nil, fmt.Errorf("invalid medication route: %s", mr)
	}
	return string(mr), nil
}

// CorrelationType represents the type of correlation between entities
type CorrelationType string

const (
	CorrelationPositive CorrelationType = "POSITIVE" // Meal causes/triggers the outcome
	CorrelationNegative CorrelationType = "NEGATIVE" // Meal prevents/improves the outcome
	CorrelationNeutral  CorrelationType = "NEUTRAL"  // No clear correlation
	CorrelationUnknown  CorrelationType = "UNKNOWN"  // Not determined yet
)

// AllCorrelationTypes returns all valid CorrelationType values
func AllCorrelationTypes() []CorrelationType {
	return []CorrelationType{
		CorrelationPositive,
		CorrelationNegative,
		CorrelationNeutral,
		CorrelationUnknown,
	}
}

// IsValid checks if the CorrelationType value is valid
func (ct CorrelationType) IsValid() bool {
	for _, valid := range AllCorrelationTypes() {
		if ct == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (ct CorrelationType) String() string {
	return string(ct)
}

// Value implements the driver.Valuer interface for database storage
func (ct CorrelationType) Value() (driver.Value, error) {
	if !ct.IsValid() {
		return nil, fmt.Errorf("invalid correlation type: %s", ct)
	}
	return string(ct), nil
}

// AuditAction represents the type of action performed for audit logging
type AuditAction string

const (
	AuditCreate AuditAction = "CREATE"
	AuditUpdate AuditAction = "UPDATE"
	AuditDelete AuditAction = "DELETE"
)

// AllAuditActions returns all valid AuditAction values
func AllAuditActions() []AuditAction {
	return []AuditAction{AuditCreate, AuditUpdate, AuditDelete}
}

// IsValid checks if the AuditAction value is valid
func (a AuditAction) IsValid() bool {
	for _, valid := range AllAuditActions() {
		if a == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (a AuditAction) String() string {
	return string(a)
}

// Value implements the driver.Valuer interface for database storage
func (a AuditAction) Value() (driver.Value, error) {
	if !a.IsValid() {
		return nil, fmt.Errorf("invalid audit action: %s", a)
	}
	return string(a), nil
}

// ParseEnum is a generic helper function to parse string values to enums
func ParseEnum[T ~string](value string, allValues []T) (T, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	for _, valid := range allValues {
		if string(valid) == normalized {
			return valid, nil
		}
	}
	var zero T
	return zero, fmt.Errorf("invalid enum value: %s", value)
}

// ParseVolume parses a string to Volume enum
func ParseVolume(value string) (Volume, error) {
	return ParseEnum(value, AllVolumes())
}

// ParseColor parses a string to Color enum
func ParseColor(value string) (Color, error) {
	return ParseEnum(value, AllColors())
}

// ParseConsistency parses a string to Consistency enum
func ParseConsistency(value string) (Consistency, error) {
	return ParseEnum(value, AllConsistencies())
}

// ParseSmellLevel parses a string to SmellLevel enum
func ParseSmellLevel(value string) (SmellLevel, error) {
	return ParseEnum(value, AllSmellLevels())
}

// ParseMealCategory parses a string to MealCategory enum
func ParseMealCategory(value string) (MealCategory, error) {
	return ParseEnum(value, AllMealCategories())
}

// ParseSymptomType parses a string to SymptomType enum
func ParseSymptomType(value string) (SymptomType, error) {
	return ParseEnum(value, AllSymptomTypes())
}

// ParseAuditAction parses a string to AuditAction enum
func ParseAuditAction(value string) (AuditAction, error) {
	return ParseEnum(value, AllAuditActions())
}

// ParseSymptomCategory parses a string to SymptomCategory enum
func ParseSymptomCategory(value string) (SymptomCategory, error) {
	return ParseEnum(value, AllSymptomCategories())
}

// ParseMedicationCategory parses a string to MedicationCategory enum
func ParseMedicationCategory(value string) (MedicationCategory, error) {
	return ParseEnum(value, AllMedicationCategories())
}

// ParseMedicationForm parses a string to MedicationForm enum
func ParseMedicationForm(value string) (MedicationForm, error) {
	return ParseEnum(value, AllMedicationForms())
}

// ParseMedicationRoute parses a string to MedicationRoute enum
func ParseMedicationRoute(value string) (MedicationRoute, error) {
	return ParseEnum(value, AllMedicationRoutes())
}

// ParseCorrelationType parses a string to CorrelationType enum
func ParseCorrelationType(value string) (CorrelationType, error) {
	return ParseEnum(value, AllCorrelationTypes())
}
