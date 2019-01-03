package ui

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"
	wlog "gopkg.in/dixonwille/wlog.v2"
)

//DefaultYN is used to specify what the default answer is to a yes/no Question.
type DefaultYN int

const (
	//DefY defaults yes/no Question to use yes.
	DefY DefaultYN = iota + 1
	//DefN defaults yes/no Question to use no.
	DefN
)

var (
	noColor = os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
)

//Menu is used to display Options to a user.
//A user can then select Options and Menu will validate the response and perform the correct action.
type Menu struct {
	Question       string
	Function       func([]Opt) error
	Options        []Opt
	UI             *Messenger
	MultiSeparator string
	AllowMultipleQ bool
	LoopOnInvalidA bool
	Clear          bool
	Tries          int
	DefIcon        string
	IsYN           bool
	YNDef          DefaultYN
}

//NewMenu creates a menu with a wlog.UI as the writer.
func NewMenu(Question string) *Menu {

	return &Menu{
		Question:       Question,
		UI:             NewMessenger(),
		MultiSeparator: " ",
		Tries:          3,
		DefIcon:        "?",
		YNDef:          0,
	}
}

//ClearOnMenuRun will clear the screen when a menu is ran.
//This is checked when LoopOnInvalid is activated.
//Meaning if an error occurred then it will clear the screen before asking again.
func (m *Menu) ClearOnMenuRun() {
	m.Clear = true
}

//SetSeparator sets the separator to use when multiple Options are valid responses.
//Default value is a space.
func (m *Menu) SetSeparator(sep string) {
	m.MultiSeparator = sep
}

//SetTries sets the number of tries on the loop before failing out.
//Default is 3.
//Negative values act like 0.
func (m *Menu) SetTries(i int) {
	m.Tries = i
}

//LoopOnInvalid is used if an invalid option was given then it will prompt the user again.
func (m *Menu) LoopOnInvalid() {
	m.LoopOnInvalidA = true
}

//SetDefaultIcon sets the icon used to identify which Options will be selected by default
func (m *Menu) SetDefaultIcon(icon string) {
	m.DefIcon = icon
}

//IsYesNo sets the menu to a yes/no state.
//Does not show Options but does ask Question.
//Will also parse the answer to allow for all variants of yes/no (IE Y yes No ...)
//Both will call the Action function you specified.
// Opt{ID: 1, Text: "y"} for yes and Opt{ID: 2, Text: "n"} for no will be passed to the function.
func (m *Menu) IsYesNo(def DefaultYN) {
	m.IsYN = true
	m.YNDef = def
}

//Option adds an option to the menu for the user to select from.
//value is an empty interface that can be used to pass anything through to the function.
//title is the string the user will select
//isDefault is whether this option is a default option (IE when no Options are selected).
//function is what is called when only this option is selected.
//If function is nil then it will default to the menu's Action.
func (m *Menu) Option(title string, value interface{}, isDefault bool, function func(Opt) error) {
	option := NewOption(len(m.Options)+1, title, value, isDefault, function)
	m.Options = append(m.Options, *option)
}

//Action adds a default action to use in certain scenarios.
//If the selected option (by default or user selected) does not have a function applied to it this will be called.
//If there are no default Options and no option was selected this will be called with an option that has an ID of -1.
func (m *Menu) Action(function func([]Opt) error) {
	m.Function = function
}

//AllowMultiple will tell the menu to allow multiple selections.
//The menu will fail if this is not called and mulple selections were selected.
func (m *Menu) AllowMultiple() {
	m.AllowMultipleQ = true
}

//ChangeReaderWriter changes where the menu listens and writes to.
//reader is where user input is collected.
//writer and errorWriter is where the menu should write to.
func (m *Menu) ChangeReaderWriter(reader io.Reader, writer, errorWriter io.Writer) {
	var UI wlog.UI
	UI = wlog.New(reader, writer, errorWriter)
	m.UI = UI
}

//Run is used to execute the menu.
//It will print to Options and Question to the screen.
//It will only clear the screen if ClearOnMenuRun is activated.
//This will validate all responses.
//Errors are of type MenuError.
func (m *Menu) Run() error {
	if m.Clear {
		Clear()
	}
	valid := false
	var Options []Opt
	//Loop and on error check if loopOnInvalid is enabled.
	//If it is Clear the screen and write error.
	//Then ask again
	for !valid {
		//step 1 print Options to screen
		m.print()
		//step 2 ask Question, get and validate response
		opt, err := m.ask()
		if err != nil {
			m.Tries = m.Tries - 1
			if !IsMenuErr(err) {
				err = newMenuError(err, "", m.triesLeft())
			}
			if m.LoopOnInvalidA && m.Tries > 0 {
				if m.Clear {
					Clear()
				}
				m.UI.UI.Error(err.Error())
			} else {
				return err
			}
		} else {
			Options = opt
			valid = true
		}
	}
	//step 3 call appropriate action with the responses
	return m.callAppropriate(Options)
}

func (m *Menu) callAppropriate(Options []Opt) (err error) {
	if len(Options) == 0 {
		return m.callAppropriateNoOptions()
	}
	if len(Options) == 1 && Options[0].function != nil {
		return Options[0].function(Options[0])
	}
	return m.Function(Options)
}

func (m *Menu) callAppropriateNoOptions() (err error) {
	Options := m.getDefault()
	if len(Options) == 0 {
		return m.Function([]Opt{{ID: -1}})
	}
	if len(Options) == 1 && Options[0].function != nil {
		return Options[0].function(Options[0])
	}
	return m.Function(Options)
}

//hide Options when this is a yes or no
func (m *Menu) print() {
	if !m.IsYN {
		for _, opt := range m.Options {
			icon := m.DefIcon
			if !opt.isDefault {
				icon = ""
			}
			m.UI.UI.Output(fmt.Sprintf("%d) %s%s", opt.ID, icon, opt.Text))
		}
	} else {
		//TODO Allow user to specify what to use as value for YN Options
		m.Options = []Opt{}
		m.Option("y", "yes", m.YNDef == DefY, nil)
		m.Option("n", "no", m.YNDef == DefN, nil)
	}
}

func (m *Menu) ask() ([]Opt, error) {
	if m.IsYN {
		if m.YNDef == DefY {
			m.Question += " (Y/n)"
		} else {
			m.Question += " (y/N)"
		}
	}
	var trim string
	if m.MultiSeparator == " " {
		trim = m.MultiSeparator
	} else {
		trim = m.MultiSeparator + " "
	}
	res, err := m.UI.UI.Ask(m.Question, trim)
	if err != nil {
		return nil, err
	}
	//Validate responses
	//Check if no responses are returned and no action to call
	if res == "" {
		//get default Options
		opt := m.getDefault()
		if !m.validOptAndFunc(opt) {
			return nil, newMenuError(ErrNoResponse, "", m.triesLeft())
		}
		return nil, nil
	}

	var responses []int
	if !m.IsYN {
		responses, err = m.resToInt(res)
		if err != nil {
			return nil, err
		}

		err = m.validateResponses(responses)
		if err != nil {
			return nil, err
		}
	} else {
		responses, err = m.ynResParse(res)
		if err != nil {
			return nil, err
		}
	}

	//Parse responses and return them as Options
	var finalOptions []Opt
	for _, response := range responses {
		finalOptions = append(finalOptions, m.Options[response-1])
	}

	return finalOptions, nil
}

//Converts the response string to a slice of ints, also validates along the way.
func (m *Menu) resToInt(res string) ([]int, error) {
	resStrings := strings.Split(res, m.MultiSeparator)
	//Check if we don't want multiple responses
	if !m.AllowMultipleQ && len(resStrings) > 1 {
		return nil, newMenuError(ErrTooMany, "", m.triesLeft())
	}

	//Convert responses to intigers
	var responses []int
	for _, response := range resStrings {
		//Check if it is an intiger
		response = strings.Trim(response, " ")
		r, err := strconv.Atoi(response)
		if err != nil {
			return nil, newMenuError(ErrInvalid, response, m.triesLeft())
		}
		responses = append(responses, r)
	}
	return responses, nil
}

func (m *Menu) ynResParse(res string) ([]int, error) {
	resStrings := strings.Split(res, m.MultiSeparator)
	if len(resStrings) > 1 {
		return nil, newMenuError(ErrTooMany, "", m.triesLeft())
	}
	re := regexp.MustCompile("^\\s*(?:([Yy])(?:es|ES)?|([Nn])(?:o|O)?)\\s*$")
	matches := re.FindStringSubmatch(res)
	if len(matches) < 2 {
		return nil, newMenuError(ErrInvalid, res, m.triesLeft())
	}
	if strings.ToLower(matches[1]) == "y" {
		return []int{int(DefY)}, nil
	}
	return []int{int(DefN)}, nil
}

//Check if response is in the range of Options
//If it is make sure it is not duplicated
func (m *Menu) validateResponses(responses []int) error {
	var tmp []int
	for _, response := range responses {
		if response < 1 || len(m.Options) < response {
			return newMenuError(ErrInvalid, strconv.Itoa(response), m.triesLeft())
		}

		if exist(tmp, response) {
			return newMenuError(ErrDuplicate, strconv.Itoa(response), m.triesLeft())
		}

		tmp = append(tmp, response)
	}
	return nil
}

//Simply checks if number exists in the slice
func exist(slice []int, number int) bool {
	for _, s := range slice {
		if number == s {
			return true
		}
	}
	return false
}

//gets a list of default Options
func (m *Menu) getDefault() []Opt {
	var opt []Opt
	for _, o := range m.Options {
		if o.isDefault {
			opt = append(opt, o)
		}
	}
	return opt
}

//make sure that there is an action available to be called in certain cases
//returns false if it chould not find an action for the number Options available
func (m *Menu) validOptAndFunc(opt []Opt) bool {
	if m.Function == nil {
		if len(opt) == 1 && opt[0].function != nil {
			return true
		}
		return false
	}
	return true
}

func (m *Menu) triesLeft() int {
	if m.LoopOnInvalidA && m.Tries > 0 {
		return m.Tries
	}
	return 0
}
