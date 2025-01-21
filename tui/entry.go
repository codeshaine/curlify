package tui

import (
	"bufio"
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/codeshaine/curlify/utils/request"
	"github.com/codeshaine/curlify/utils/response"
)

const (
	spacer      = 1
	wSize       = 0.92
	methodWidth = 0.15
	urlWidth    = 0.83
)

type Mode string

const (
	NormalMode = "NORMAL"
	EditMode   = "EDIT"
)

type Dim struct {
	width  int
	height int
}

type Model struct {
	Focus        int
	message      string
	Input        string
	MethodInput  textinput.Model
	URLInput     textinput.Model
	HeaderInput  textarea.Model
	Output       string
	Result       viewport.Model
	Ready        bool
	Dimension    Dim
	Mode         Mode
	RequestType  string
	HeaderValues map[string][]string
}

func NewModel() Model {
	mi := textinput.New()
	mi.Placeholder = "GET"
	mi.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	mi.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	mi.Width = 10

	ui := textinput.New()
	ui.Placeholder = "Enter URL..."
	ui.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	ui.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

	hi := textarea.New()
	hi.Placeholder = "Enter headers (key: value)..."
	hi.ShowLineNumbers = false
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
		Focus:        0,
		message:      "",
		Input:        "",
		Output:       "",
		Ready:        false,
		MethodInput:  mi,
		URLInput:     ui,
		HeaderInput:  hi,
		Result:       vp,
		Mode:         NormalMode,
		RequestType:  "GET",
		HeaderValues: make(map[string][]string),
		Dimension:    Dim{height: 0, width: 0},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, textinput.Blink)
}

func makeRequest(m *Model) error {
	header := request.ParseHeaders(m.HeaderInput.Value())
	req := request.NewGet(m.URLInput.Value(), header)

	res, err := req.Do()
	if err != nil {
		m.Output = "Invalid URL"
		m.Result.SetContent(m.Output)
		return err
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
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Mode == EditMode {
				m.Mode = NormalMode
				m.MethodInput.Blur()
				m.URLInput.Blur()
				m.HeaderInput.Blur()
				return m, nil
			}
		case "i":
			if m.Mode == NormalMode {
				m.Mode = EditMode
				switch m.Focus {
				case 0:
					m.MethodInput.Focus()
				case 1:
					m.URLInput.Focus()
				case 2:
					m.HeaderInput.Focus()
				}
				return m, nil
			}
		}

		if m.Mode == NormalMode {
			switch msg.String() {
			case "j":
				if m.Focus < 3 {
					m.Focus++
					m.MethodInput.Blur()
					m.URLInput.Blur()
					m.HeaderInput.Blur()
				}
			case "k":
				if m.Focus > 0 {
					m.Focus--
					if m.Mode == EditMode {
						switch m.Focus {
						case 0:
							m.MethodInput.Focus()
						case 1:
							m.URLInput.Focus()
						case 2:
							m.HeaderInput.Focus()
						}
					}
				}
			case "q":
				return m, tea.Quit
			}
		} else {
			if m.Focus != 3 {
				switch msg.String() {
				case "enter":
					if m.Focus != 2 && m.URLInput.Value() != "" {
						makeRequest(&m)
					}
				default:
					switch m.Focus {
					case 0:
						m.MethodInput, cmd = m.MethodInput.Update(msg)
					case 1:
						m.URLInput, cmd = m.URLInput.Update(msg)
					case 2:
						var textAreaCmd tea.Cmd
						m.HeaderInput, textAreaCmd = m.HeaderInput.Update(msg)
						cmds = append(cmds, textAreaCmd)
					}
					cmds = append(cmds, cmd)
				}
			} else {
				m.Result, cmd = m.Result.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case tea.WindowSizeMsg:
		m.Dimension.width = msg.Width
		m.Dimension.height = msg.Height

		// Fixed height calculations
		topInputHeight := 3                                             // Fixed height for method/URL row
		headerHeight := int(math.Max(float64(m.Dimension.height)/5, 6)) // Min height of 6, max 20% of screen
		statusHeight := 1
		spacerHeight := spacer * 3
		resultHeight := m.Dimension.height - topInputHeight - headerHeight - statusHeight - spacerHeight - 4 // Account for borders

		// Update components with new dimensions
		m.HeaderInput.SetHeight(headerHeight - 2) // Account for borders
		if !m.Ready {
			m.Result = viewport.New(
				int(math.Floor(wSize*float64(m.Dimension.width))),
				resultHeight,
			)
			m.Ready = true
		} else {
			m.Result.Width = int(math.Floor(wSize * float64(m.Dimension.width)))
			m.Result.Height = resultHeight
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	modeColor := map[Mode]string{
		NormalMode: "69",
		EditMode:   "156",
	}

	focusColor := "39"
	borderColor := "244"

	// Status bar
	statusBar := lipgloss.NewStyle().
		Background(lipgloss.Color(modeColor[m.Mode])).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 1).
		Render(fmt.Sprintf(" %s MODE | Focus: %s ",
			m.Mode,
			map[int]string{
				0: "METHOD",
				1: "URL",
				2: "HEADERS",
				3: "RESULT",
			}[m.Focus]))

	// Calculate widths
	topWidth := int(math.Floor(wSize * float64(m.Dimension.width)))
	methodBoxWidth := int(math.Floor(methodWidth * float64(topWidth)))
	urlBoxWidth := int(math.Floor(urlWidth * float64(topWidth)))

	// Method input box with fixed height
	methodStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1).
		Height(3).
		Width(methodBoxWidth)

	if m.Focus == 0 {
		methodStyle = methodStyle.BorderForeground(lipgloss.Color(focusColor))
	}
	methodBox := methodStyle.Render(m.MethodInput.View())

	// URL input box with fixed height
	urlStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1).
		Height(3).
		Width(urlBoxWidth)

	if m.Focus == 1 {
		urlStyle = urlStyle.BorderForeground(lipgloss.Color(focusColor))
	}
	urlBox := urlStyle.Render(m.URLInput.View())

	// Set textarea width
	m.HeaderInput.SetWidth(topWidth - 4) // Account for borders and padding

	// Headers section
	headerBox := m.HeaderInput.View()
	if m.Focus == 2 {
		m.HeaderInput.Focus()
	} else {
		m.HeaderInput.Blur()
	}

	// Result viewport
	viewportStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Width(topWidth)

	if m.Focus == 3 {
		viewportStyle = viewportStyle.BorderForeground(lipgloss.Color(focusColor))
	}

	// Navigation hints
	hint := ""
	if m.Mode == NormalMode {
		hint = "\n[i: edit mode, j/k: navigation]"
	} else {
		hint = "\n[esc: normal mode, enter: make request]"
	}

	viewportContent := viewportStyle.Render(m.Result.View() + hint)

	// Top row with fixed height
	topRow := lipgloss.JoinHorizontal(
		lipgloss.Left,
		methodBox,
		strings.Repeat(" ", 2),
		urlBox,
	)

	// Final layout
	return lipgloss.JoinVertical(
		lipgloss.Top,
		statusBar,
		topRow,
		strings.Repeat(" ", spacer),
		headerBox,
		strings.Repeat(" ", spacer),
		viewportContent,
	)
}
