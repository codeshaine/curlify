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
	mi.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	mi.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	mi.Width = 6

	ui := textinput.New()
	ui.Placeholder = "Enter URL..."
	ui.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	ui.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

	// mi.Width = 44

	hi := textarea.New()
	hi.Placeholder = "Enter headers (key: value)..."
	hi.ShowLineNumbers = true
	hi.Prompt = ""
	hi.FocusedStyle.Base = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("39")).
		Border(lipgloss.RoundedBorder())
	hi.BlurredStyle.Base = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("244")).
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
