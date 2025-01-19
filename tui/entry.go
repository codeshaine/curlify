package tui

import (
	"bufio"
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/codeshaine/curlify/utils/request"
	"github.com/codeshaine/curlify/utils/response"
)

const (
	// Existing constants
	spacer      = 1
	wSize       = 0.98
	botPanHSize = 0.50
	topPanHSize = 0.40

	// Mode constants
	NormalMode = "NORMAL"
	EditMode   = "EDIT"
)

type Mode string

type Dim struct {
	width  int
	height int
}

type Model struct {
	Focus     int
	message   string
	Input     string
	TextInput textinput.Model
	Output    string
	Result    viewport.Model
	Ready     bool
	Dimension Dim
	Mode      Mode
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter URL..."
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	vp := viewport.New(0, 0)
	vp.SetContent("")
	return Model{
		Focus:     0,
		message:   "",
		Input:     "",
		Output:    "",
		Ready:     false,
		TextInput: ti,
		Result:    vp,         //view port
		Mode:      NormalMode, // Start in normal mode
		Dimension: Dim{
			height: 0,
			width:  0,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.WindowSize(),
		textinput.Blink,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle mode switching and global commands
		switch msg.String() {
		case "esc":
			if m.Mode == EditMode {
				m.Mode = NormalMode
				m.TextInput.Blur()
				return m, nil
			}
		case "i":
			if m.Mode == NormalMode {
				m.Mode = EditMode
				if m.Focus == 0 {
					m.TextInput.Focus()
				}
				return m, nil
			}
		}

		// Handle mode-specific commands
		if m.Mode == NormalMode {
			switch msg.String() {
			case "j": // Move focus down
				if m.Focus < 1 {
					m.Focus++
					m.TextInput.Blur()
				}
			case "k": // Move focus up
				if m.Focus > 0 {
					m.Focus--
					if m.Mode == EditMode {
						m.TextInput.Focus()
					}
				}
			case "q":
				return m, tea.Quit

				// case "pagedown": // Page down
				// 	if m.Focus == 1 {
				// 		m.Result.ViewDown()
				// 	}
				// case "pageup": // Page up
				// 	if m.Focus == 1 {
				// 	}
			}
		} else { // Edit mode
			if m.Focus == 0 {
				// Handle input field editing
				switch msg.String() {
				case "enter":
					header := make(map[string][]string)
					header["Content-Type"] = []string{"application/json"}
					req := request.NewGet(m.TextInput.Value(), header)

					res, err := req.Do()
					if err != nil {
						m.Output = "Invalid URL"
						m.Result.SetContent(m.Output)
						break
					}

					defer res.Body.Close()
					var output strings.Builder
					data := bufio.NewScanner(res.Body)
					for data.Scan() {
						output.WriteString(data.Text() + "\n")
					}
					m.Output = response.FormatJSONResponse(output.String())
					m.Result.SetContent(m.Output)
					m.Result.GotoTop()

				default:
					m.TextInput, cmd = m.TextInput.Update(msg)
					cmds = append(cmds, cmd)
				}
			} else if m.Focus == 1 {
				m.Result, cmd = m.Result.Update(msg)
				cmds = append(cmds, cmd)

			}
		}

	case tea.WindowSizeMsg:
		m.Dimension.width = msg.Width
		m.Dimension.height = msg.Height
		headerHeight := int(math.Floor(topPanHSize * float64(m.Dimension.height)))
		viewportHeight := m.Dimension.height - headerHeight - spacer - 10

		if !m.Ready {
			m.Result = viewport.New(
				int(math.Floor(wSize*float64(m.Dimension.width))),
				viewportHeight,
			)
			m.Ready = true
		} else {
			// Subsequent resize
			m.Result.Width = int(math.Floor(wSize * float64(m.Dimension.width)))
			m.Result.Height = viewportHeight
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	modeColor := map[Mode]string{
		NormalMode: "212",
		EditMode:   "114",
	}

	statusBar := lipgloss.NewStyle().
		Background(lipgloss.Color(modeColor[m.Mode])).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 1).
		Render(fmt.Sprintf(" %s MODE | Focus: %s ",
			m.Mode,
			map[int]string{
				0: "INPUT",
				1: "RESULT",
			}[m.Focus]))

	// Input section with hint text
	inputHint := ""
	if m.Mode == NormalMode {
		inputHint = "\n[i: edit mode, j/k: navigation]"
	} else {
		inputHint = "\n[esc: normal mode]"
	}

	topContent := m.TextInput.View() + inputHint

	// Style the panes based on focus and mode
	topPane := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Foreground(lipgloss.Color("104")).
		Width(int(math.Floor(wSize * float64(m.Dimension.width)))).
		Height(int(math.Floor(topPanHSize * float64(m.Dimension.height)))).
		Render(topContent)

	if m.Focus == 0 {
		topPane = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(1, 2).
			Foreground(lipgloss.Color("205")).
			Width(int(math.Floor(wSize * float64(m.Dimension.width)))).
			Height(int(math.Floor(topPanHSize * float64(m.Dimension.height)))).
			Render(topContent)
	}

	viewportStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(int(math.Floor(wSize * float64(m.Dimension.width))))

	if m.Focus == 1 {
		viewportStyle = viewportStyle.BorderForeground(lipgloss.Color("205"))
	}

	scrollHint := ""
	if m.Focus == 1 {
		if m.Mode == NormalMode {
			scrollHint = "\n[h/l: scroll, ctrl+f/ctrl+b: page up/down]"
		} else {
			scrollHint = "\n[↑/↓: scroll, PgUp/PgDn: page up/down]"
		}
	}

	viewportContent := viewportStyle.Render(m.Result.View() + scrollHint)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		statusBar,
		topPane,
		strings.Repeat(" ", spacer),
		viewportContent,
	)
}
