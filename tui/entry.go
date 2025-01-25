package tui

import (
	"fmt"
	"math"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, textinput.Blink)
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
				if m.IsHeader {
					m.HeaderValues = m.HeaderBodyInput.Value()
				} else {
					m.BodyValues = m.HeaderBodyInput.Value()
				}
				m.Mode = NormalMode
				m.IsHeader = false //exiting the header
				m.MethodInput.Blur()
				m.URLInput.Blur()
				m.HeaderBodyInput.SetValue(m.BodyValues)
				m.HeaderBodyInput.Blur()
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
					m.HeaderBodyInput.Focus()
					m.HeaderBodyInput.SetValue(m.BodyValues)

				}
				return m, nil
			}
		case "h": //for header section
			if m.Mode == NormalMode {
				switch m.Focus {
				case 2:
					// m.HeaderBodyInput.SetValue(m.BodyValues)
					m.IsHeader = true
					m.HeaderBodyInput.Focus()
					m.HeaderBodyInput.SetValue(m.HeaderValues)
					m.Mode = EditMode
				}
				return m, nil
			}
		case "g": //change it to shift+enter (well shift+enter is not the right word guess : HELP NEEDED)
			if m.URLInput.Value() != "" && m.Mode == NormalMode {
				m.makeRequest()

			}

		}

		if m.Mode == NormalMode {
			switch msg.String() {
			case "j":
				if m.Focus < 3 {
					m.Focus++
					m.MethodInput.Blur()
					m.URLInput.Blur()
					m.HeaderBodyInput.Blur()
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
							m.HeaderBodyInput.Focus()
						}
					}
				}
			case "q":
				return m, tea.Quit
			}
		} else {
			if m.Focus != 3 {
				switch msg.String() {
				default:
					switch m.Focus {
					case 0:
						m.MethodInput, cmd = m.MethodInput.Update(msg)
					case 1:
						m.URLInput, cmd = m.URLInput.Update(msg)
					case 2:
						var textAreaCmd tea.Cmd
						m.HeaderBodyInput, textAreaCmd = m.HeaderBodyInput.Update(msg)
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

		topInputHeight := 3
		headerHeight := int(math.Max(float64(m.Dimension.height)/5, 6))
		statusHeight := 1
		spacerHeight := spacer * 3
		resultHeight := m.Dimension.height - topInputHeight - headerHeight - statusHeight - spacerHeight - 4

		m.HeaderBodyInput.SetHeight(headerHeight)
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
		NormalMode: "#bd93f9",
		EditMode:   "#50fa7b",
	}

	DataSection := map[int]string{
		0: "METHOD",
		1: "URL",
		2: "BODY",
		3: "RESULT",
	}[m.Focus]
	if m.IsHeader {
		DataSection = "HEADER"
	}
	statusBar := lipgloss.NewStyle().
		Background(lipgloss.Color(modeColor[m.Mode])).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 1).
		Render(fmt.Sprintf(" %s MODE | Focus: %s ",
			m.Mode,
			DataSection))

	topWidth := int(math.Floor(wSize * float64(m.Dimension.width)))
	methodBoxWidth := int(math.Floor(methodWidth * float64(topWidth)))
	urlBoxWidth := int(math.Floor(urlWidth * float64(topWidth)))

	methodStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(1, 1).
		Height(1).
		Width(methodBoxWidth)

	if m.Focus == 0 {
		methodStyle = methodStyle.BorderForeground(lipgloss.Color(focusColor))
	}

	urlStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(1, 1).
		Height(1).
		Width(urlBoxWidth)

	if m.Focus == 1 {
		urlStyle = urlStyle.BorderForeground(lipgloss.Color(focusColor))
	}

	m.HeaderBodyInput.SetWidth(topWidth)

	if m.Focus == 2 {
		m.HeaderBodyInput.Focus()
	} else {
		m.HeaderBodyInput.Blur()
	}

	viewportStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Width(topWidth)

	if m.Focus == 3 {
		viewportStyle = viewportStyle.BorderForeground(lipgloss.Color(focusColor))
	}

	hint := ""
	if m.Mode == NormalMode {
		hint = "\n[i: edit mode h: header, j/k: navigation]"
	} else {
		hint = "\n[esc: normal mode, g: make request]"
	}

	methodBox := methodStyle.Render(m.MethodInput.View())
	urlBox := urlStyle.Render(m.URLInput.View())
	headerBox := m.HeaderBodyInput.View()
	viewportContent := viewportStyle.Render(m.Result.View() + hint)

	topRow := lipgloss.JoinHorizontal(
		lipgloss.Left,
		methodBox,
		urlBox,
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		statusBar,
		topRow,
		headerBox,
		viewportContent,
	)
}
