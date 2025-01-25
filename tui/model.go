package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Focus           int
	message         string //haven't decided what to do with this
	MethodInput     textinput.Model
	URLInput        textinput.Model
	HeaderBodyInput textarea.Model
	Result          viewport.Model
	Ready           bool //for viewport
	Dimension       Dim
	Mode            Mode
	RequestType     string
	HeaderValues    string
	IsHeader        bool
	BodyValues      string
}

type Dim struct {
	width  int
	height int
}

type Mode string

func NewModel() Model {
	mi := textinput.New()
	mi.Placeholder = "GET"
	mi.SetValue("GET")
	mi.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(textColor))
	mi.Width = 6

	ui := textinput.New()
	ui.Placeholder = "https://example.com"
	ui.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(textColor))

	hi := textarea.New()
	hi.Placeholder = `{"key":"value"} or Content-Type: application/json`
	hi.ShowLineNumbers = true
	hi.Prompt = ""
	hi.FocusedStyle.Base = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(focusColor)).
		Border(lipgloss.RoundedBorder()).Foreground(lipgloss.Color(textColor))
	hi.BlurredStyle.Base = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(borderColor)).
		Border(lipgloss.RoundedBorder())

	vp := viewport.New(0, 0)
	vp.SetContent("")
	return Model{
		Focus:           0,
		message:         "",
		Ready:           false,
		MethodInput:     mi,
		URLInput:        ui,
		HeaderBodyInput: hi,
		Result:          vp,
		Mode:            NormalMode,
		RequestType:     "GET",
		HeaderValues:    "",
		BodyValues:      "",
		IsHeader:        false,
		Dimension:       Dim{height: 0, width: 0},
	}
}
