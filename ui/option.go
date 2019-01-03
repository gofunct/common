package ui

type Opt struct {
	ID        int
	Text      string
	Value     interface{}
	function  func(Opt) error
	isDefault bool
}

func NewOption(id int, text string, value interface{}, def bool, function func(Opt) error) *Opt {
	return &Opt{
		ID:        id,
		Text:      text,
		Value:     value,
		isDefault: def,
		function:  function,
	}
}
