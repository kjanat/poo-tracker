package shared

import (
	"database/sql/driver"
	"fmt"
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

// SmellLevel represents the smell intensity level
type SmellLevel string

const (
	SmellLevelNone         SmellLevel = "NONE"
	SmellLevelMild         SmellLevel = "MILD"
	SmellLevelModerate     SmellLevel = "MODERATE"
	SmellLevelStrong       SmellLevel = "STRONG"
	SmellLevelOverwhelming SmellLevel = "OVERWHELMING"
)

// AllSmellLevels returns all valid SmellLevel values
func AllSmellLevels() []SmellLevel {
	return []SmellLevel{SmellLevelNone, SmellLevelMild, SmellLevelModerate, SmellLevelStrong, SmellLevelOverwhelming}
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

// MealCategory represents the category/type of a meal
type MealCategory string

const (
	MealCategoryBreakfast MealCategory = "BREAKFAST"
	MealCategoryLunch     MealCategory = "LUNCH"
	MealCategoryDinner    MealCategory = "DINNER"
	MealCategorySnack     MealCategory = "SNACK"
	MealCategoryOther     MealCategory = "OTHER"
)

// AllMealCategories returns all valid MealCategory values
func AllMealCategories() []MealCategory {
	return []MealCategory{MealCategoryBreakfast, MealCategoryLunch, MealCategoryDinner, MealCategorySnack, MealCategoryOther}
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
