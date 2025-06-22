package shared

import (
	"database/sql"        // For sql.Scanner interface
	"database/sql/driver" // For driver.Valuer interface
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

// ParseVolume converts a string to Volume with validation
func ParseVolume(s string) (Volume, error) {
	v := Volume(s)
	if !v.IsValid() {
		return "", fmt.Errorf("invalid volume: %s", s)
	}
	return v, nil
}

// Scan implements the sql.Scanner interface for database reading
func (v *Volume) Scan(value interface{}) error {
	if value == nil {
		*v = ""
		return nil
	}
	switch s := value.(type) {
	case string:
		parsed, err := ParseVolume(s)
		if err != nil {
			return err
		}
		*v = parsed
		return nil
	case []byte:
		parsed, err := ParseVolume(string(s))
		if err != nil {
			return err
		}
		*v = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Volume", value)
	}
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

// ParseColor converts a string to Color with validation
func ParseColor(s string) (Color, error) {
	c := Color(s)
	if !c.IsValid() {
		return "", fmt.Errorf("invalid color: %s", s)
	}
	return c, nil
}

// Scan implements the sql.Scanner interface for database reading
func (c *Color) Scan(value interface{}) error {
	if value == nil {
		*c = ""
		return nil
	}
	switch s := value.(type) {
	case string:
		parsed, err := ParseColor(s)
		if err != nil {
			return err
		}
		*c = parsed
		return nil
	case []byte:
		parsed, err := ParseColor(string(s))
		if err != nil {
			return err
		}
		*c = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Color", value)
	}
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

// ParseConsistency converts a string to Consistency with validation
func ParseConsistency(s string) (Consistency, error) {
	c := Consistency(s)
	if !c.IsValid() {
		return "", fmt.Errorf("invalid consistency: %s", s)
	}
	return c, nil
}

// Scan implements the sql.Scanner interface for database reading
func (c *Consistency) Scan(value interface{}) error {
	if value == nil {
		*c = ""
		return nil
	}
	switch s := value.(type) {
	case string:
		parsed, err := ParseConsistency(s)
		if err != nil {
			return err
		}
		*c = parsed
		return nil
	case []byte:
		parsed, err := ParseConsistency(string(s))
		if err != nil {
			return err
		}
		*c = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Consistency", value)
	}
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

// ParseSmellLevel converts a string to SmellLevel with validation
func ParseSmellLevel(s string) (SmellLevel, error) {
	sl := SmellLevel(s)
	if !sl.IsValid() {
		return "", fmt.Errorf("invalid smell level: %s", s)
	}
	return sl, nil
}

// Scan implements the sql.Scanner interface for database reading
func (s *SmellLevel) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	switch v := value.(type) {
	case string:
		parsed, err := ParseSmellLevel(v)
		if err != nil {
			return err
		}
		*s = parsed
		return nil
	case []byte:
		parsed, err := ParseSmellLevel(string(v))
		if err != nil {
			return err
		}
		*s = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into SmellLevel", value)
	}
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

// ParseMealCategory converts a string to MealCategory with validation
func ParseMealCategory(s string) (MealCategory, error) {
	mc := MealCategory(s)
	if !mc.IsValid() {
		return "", fmt.Errorf("invalid meal category: %s", s)
	}
	return mc, nil
}

// Scan implements the sql.Scanner interface for database reading
func (m *MealCategory) Scan(value interface{}) error {
	if value == nil {
		*m = ""
		return nil
	}
	switch s := value.(type) {
	case string:
		parsed, err := ParseMealCategory(s)
		if err != nil {
			return err
		}
		*m = parsed
		return nil
	case []byte:
		parsed, err := ParseMealCategory(string(s))
		if err != nil {
			return err
		}
		*m = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into MealCategory", value)
	}
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

// ParseSymptomCategory converts a string to SymptomCategory with validation
func ParseSymptomCategory(s string) (SymptomCategory, error) {
	sc := SymptomCategory(s)
	if !sc.IsValid() {
		return "", fmt.Errorf("invalid symptom category: %s", s)
	}
	return sc, nil
}

// Scan implements the sql.Scanner interface for database reading
func (sc *SymptomCategory) Scan(value interface{}) error {
	if value == nil {
		*sc = ""
		return nil
	}
	switch s := value.(type) {
	case string:
		parsed, err := ParseSymptomCategory(s)
		if err != nil {
			return err
		}
		*sc = parsed
		return nil
	case []byte:
		parsed, err := ParseSymptomCategory(string(s))
		if err != nil {
			return err
		}
		*sc = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into SymptomCategory", value)
	}
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

// ParseSymptomType converts a string to SymptomType with validation
func ParseSymptomType(s string) (SymptomType, error) {
	st := SymptomType(s)
	if !st.IsValid() {
		return "", fmt.Errorf("invalid symptom type: %s", s)
	}
	return st, nil
}

// Scan implements the sql.Scanner interface for database reading
func (s *SymptomType) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	switch v := value.(type) {
	case string:
		parsed, err := ParseSymptomType(v)
		if err != nil {
			return err
		}
		*s = parsed
		return nil
	case []byte:
		parsed, err := ParseSymptomType(string(v))
		if err != nil {
			return err
		}
		*s = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into SymptomType", value)
	}
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

// ParseMedicationCategory converts a string to MedicationCategory with validation
func ParseMedicationCategory(s string) (MedicationCategory, error) {
	mc := MedicationCategory(s)
	if !mc.IsValid() {
		return "", fmt.Errorf("invalid medication category: %s", s)
	}
	return mc, nil
}

// Scan implements the sql.Scanner interface for database reading
func (mc *MedicationCategory) Scan(value interface{}) error {
	if value == nil {
		*mc = ""
		return nil
	}
	switch s := value.(type) {
	case string:
		parsed, err := ParseMedicationCategory(s)
		if err != nil {
			return err
		}
		*mc = parsed
		return nil
	case []byte:
		parsed, err := ParseMedicationCategory(string(s))
		if err != nil {
			return err
		}
		*mc = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into MedicationCategory", value)
	}
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

// ParseMedicationForm converts a string to MedicationForm with validation
func ParseMedicationForm(s string) (MedicationForm, error) {
	mf := MedicationForm(s)
	if !mf.IsValid() {
		return "", fmt.Errorf("invalid medication form: %s", s)
	}
	return mf, nil
}

// Scan implements the sql.Scanner interface for database reading
func (mf *MedicationForm) Scan(value interface{}) error {
	if value == nil {
		*mf = ""
		return nil
	}
	switch s := value.(type) {
	case string:
		parsed, err := ParseMedicationForm(s)
		if err != nil {
			return err
		}
		*mf = parsed
		return nil
	case []byte:
		parsed, err := ParseMedicationForm(string(s))
		if err != nil {
			return err
		}
		*mf = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into MedicationForm", value)
	}
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

// ParseMedicationRoute converts a string to MedicationRoute with validation
func ParseMedicationRoute(s string) (MedicationRoute, error) {
	mr := MedicationRoute(s)
	if !mr.IsValid() {
		return "", fmt.Errorf("invalid medication route: %s", s)
	}
	return mr, nil
}

// Scan implements the sql.Scanner interface for database reading
func (mr *MedicationRoute) Scan(value interface{}) error {
	if value == nil {
		*mr = ""
		return nil
	}
	switch s := value.(type) {
	case string:
		parsed, err := ParseMedicationRoute(s)
		if err != nil {
			return err
		}
		*mr = parsed
		return nil
	case []byte:
		parsed, err := ParseMedicationRoute(string(s))
		if err != nil {
			return err
		}
		*mr = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into MedicationRoute", value)
	}
}

// Compile-time interface checks
var (
	_ driver.Valuer = (*Volume)(nil)
	_ sql.Scanner   = (*Volume)(nil)
	_ driver.Valuer = (*Color)(nil)
	_ sql.Scanner   = (*Color)(nil)
	_ driver.Valuer = (*Consistency)(nil)
	_ sql.Scanner   = (*Consistency)(nil)
	_ driver.Valuer = (*SmellLevel)(nil)
	_ sql.Scanner   = (*SmellLevel)(nil)
	_ driver.Valuer = (*MealCategory)(nil)
	_ sql.Scanner   = (*MealCategory)(nil)
	_ driver.Valuer = (*SymptomCategory)(nil)
	_ sql.Scanner   = (*SymptomCategory)(nil)
	_ driver.Valuer = (*SymptomType)(nil)
	_ sql.Scanner   = (*SymptomType)(nil)
	_ driver.Valuer = (*MedicationCategory)(nil)
	_ sql.Scanner   = (*MedicationCategory)(nil)
	_ driver.Valuer = (*MedicationForm)(nil)
	_ sql.Scanner   = (*MedicationForm)(nil)
	_ driver.Valuer = (*MedicationRoute)(nil)
	_ sql.Scanner   = (*MedicationRoute)(nil)
)
