package types

type ProfileId = int64

type NetworkControlInfo struct {
	ProfileID ProfileId `json:"profile_id"` // Id of the profile this network control is associated with. This is read-only, unless you use the POST api to add a network control.
	NextChange int `json:"next_change"` // UNIX timestamp of next rule change in seconds. 0 if no next change.
	OverrideMode RuleMode `json:"override_mode"` // Mode of current override.
	CurrentMode RuleMode `json:"current_mode"` // Mode in use. If override is true, it will be override_mode, otherwise it’s the mode from the rules attached to this NetworkControl.
	RuleMode RuleMode `json:"rule_mode"` // Mode that would be in use if there was no override. Depends only on rules, and is useful to determine what will happen when override is lifted.
	OverrideUntil int `json:"override_until"` // Unix timestamp in seconds when override ends. Relevant when override is true. Set at 0 for unlimited.
	Override bool `json:"override"` // Whether there’s an override at the moment.
	Macs []string `json:"macs"` // List of mac adresses associated with this profile’s network control.
	Hosts []LanInterfaceHost `json:"hosts"` // List of Lan Host objects associated with this profile’s network control. Derived from the macs array.
	Resolution int `json:"resolution"` // Control resolution per day of this network control.
	CustomDayRanges[]DayRange `json:"cdayranges"` // List of custom day range, each custom day range represents a group of days for which you want to use a different planning than other week days.
}

type DayRange = string

const (
	DayRangeFrenchBankHolidays DayRange = ":fr_bank_holidays" // French bank holidays
	DayRangeFrenchSchoolHolidaysA DayRange = ":fr_school_holidays_a" // French school holidays - Zone A
	DayRangeFrenchSchoolHolidaysB DayRange = ":fr_school_holidays_b" // French school holidays - Zone B
	DayRangeFrenchSchoolHolidaysC DayRange = ":fr_school_holidays_c" // French school holidays - Zone C
	DayRangeFrenchSchoolHolidaysCorse DayRange = ":fr_school_holidays_corse" // French school holidays - Corse
)

type RuleMode = string

const (
	RuleModeAllowed RuleMode = "allowed" // Allowed
	RuleModeDenied RuleMode = "denied" // Denied
	RuleModeWebOnly RuleMode = "webonly" // Web only
)

type NetworkControlPayload struct {
	ProfileID ProfileId `json:"profile_id,omitempty"` // Id of the profile this network control is associated with. This is read-only, unless you use the POST api to add a network control.
	NextChange int `json:"next_change,omitempty"` // UNIX timestamp of next rule change in seconds. 0 if no next change.
	OverrideMode RuleMode `json:"override_mode"` // Mode of current override.
	CurrentMode RuleMode `json:"current_mode,omitempty"` // Mode in use. If override is true, it will be override_mode, otherwise it’s the mode from the rules attached to this NetworkControl.
	RuleMode RuleMode `json:"rule_mode,omitempty"` // Mode that would be in use if there was no override. Depends only on rules, and is useful to determine what will happen when override is lifted.
	OverrideUntil int `json:"override_until"` // Unix timestamp in seconds when override ends. Relevant when override is true. Set at 0 for unlimited.
	Override bool `json:"override"` // Whether there’s an override at the moment.
	Macs []string `json:"macs"` // List of mac adresses associated with this profile’s network control.
	Hosts []LanHost `json:"hosts,omitempty"` // List of Lan Host objects associated with this profile’s network control. Derived from the macs array.
	Resolution int `json:"resolution,omitempty"` // Control resolution per day of this network control.
	CustomDayRanges[]DayRange `json:"cdayranges"` // list of custom day range, each custom day range represents a group of days for which you want to use a different planning than other week days.
}
