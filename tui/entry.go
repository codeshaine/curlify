package tui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/codeshaine/curlify/pkg/request"
)

var (
	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Foreground(lipgloss.Color("205")).
			BorderForeground(lipgloss.Color("63")).Width(100).Height(15)
	normalStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 2).
			Foreground(lipgloss.Color("120")).Width(100).Height(13)
	spacer = 1
)

type Model struct {
	Focus       int
	InputBuffer string
	Output      string
}

func NewModel() Model {
	return Model{
		Focus:       0,
		InputBuffer: "",
		Output:      "",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func formatJSONResponse(rawJSON string) string {
	var jsonResponse map[string]interface{}

	// Unmarshal the raw JSON string into a Go map
	if err := json.Unmarshal([]byte(rawJSON), &jsonResponse); err != nil {
		return fmt.Sprintf("Error decoding JSON: %v", err)
	}

	// Marshal the JSON with indentation
	prettyJSON, err := json.MarshalIndent(jsonResponse, "\t", "\t") // Using tabs for indentation
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}

	// Return the formatted JSON string
	return string(prettyJSON)
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.Focus > 0 {
				m.Focus--
			}
		case "down":
			if m.Focus < 1 {
				m.Focus++
			}
		case "q":
			return m, tea.Quit
		case "enter":
			if m.Focus == 0 {
				header := make(map[string][]string)
				header["Content-Type"] = []string{"application/json"}

				getReq := request.NewGet(m.InputBuffer, header)
				res, err := getReq.Do()
				if err != nil {
					panic(err)
				}
				defer res.Body.Close()

				data := bufio.NewScanner(res.Body)
				for data.Scan() {
					m.Output += data.Text() + "\n"
				}
				// m.Output = formatJSONResponse(m.Output)
				m.InputBuffer = ""
			}
		case "backspace":
			if m.Focus == 0 && len(m.InputBuffer) > 0 {
				m.InputBuffer = m.InputBuffer[:len(m.InputBuffer)-1]
			}
		case "ctrl+d":
			if m.Focus == 0 {
				m.InputBuffer += "https://dummyjson.com/recipes/1"
			}
		default:
			if m.Focus == 0 {
				m.InputBuffer += msg.String()
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	topContent := "Input: " + m.InputBuffer + "\n\nUse 'up'/'down' to switch focus."
	bottomContent := "Output: " + m.Output + "\n\nPress 'q' to quit."

	topPane := normalStyle.Render(topContent)
	bottomPane := normalStyle.Render(bottomContent)

	if m.Focus == 0 {
		topPane = focusedStyle.Render(topContent)
	} else {
		bottomPane = focusedStyle.Render(bottomContent)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		topPane,
		strings.Repeat(" ", spacer),
		bottomPane,
	)
}
