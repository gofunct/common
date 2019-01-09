package runtime

type FlagValue struct{}

func (f FlagValue) HasChanged() bool    { return false }
func (f FlagValue) Name() string        { return "my-flag-name" }
func (f FlagValue) ValueString() string { return "my-flag-value" }
func (f FlagValue) ValueType() string   { return "string" }
