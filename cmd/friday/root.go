package friday

import (
	"fmt"
	"os"
	"strings"

	"github.com/Orgate-AI/friday-cli/utils"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle = lipgloss.NewStyle().Foreground(hotPink)
)

var (
	senderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	// botStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	// errorStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1"))
	footerStyle = lipgloss.NewStyle().
			Height(1).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("8")).
			Faint(true)
)

type (
	errMsg error
)

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	width       int
	height      int
	err         error
}

func initialModel() model {
	ti := textarea.New()
	ti.Placeholder = "write a query ..."
	ti.Focus()

	ti.Prompt = "┃ "
	ti.CharLimit = 280

	ti.SetWidth(50)
	ti.SetHeight(1)

	// Remove cursor line styling
	ti.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ti.ShowLineNumbers = false

	vp := viewport.New(50, 5)
	ti.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ti,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(m.RenderFooter())
		m.textarea.SetWidth(msg.Width)
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			input := strings.TrimSpace(m.textarea.Value())

			if input == "" {
				break
			}
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			resp := utils.RunQuery(input, "somehign")
			m.messages = append(m.messages, m.senderStyle.Render("Friday: ")+string(resp))
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) RenderFooter() string {
	var columns []string

	columns = append(columns, fmt.Sprintf("%s ctrl+h"))

	totalWidth := lipgloss.Width(strings.Join(columns, ""))
	padding := 2

	if totalWidth+(len(columns)-1)*padding > m.width {
		remainingSpace := 5
		columns[len(columns)-1] = columns[len(columns)-1][:remainingSpace] + "..."
	}

	footer := strings.Join(columns, strings.Repeat(" ", padding))
	footer = footerStyle.Render(footer)
	return footer
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.viewport.View(),
		m.textarea.View(),
		m.RenderFooter(),
	)
}

var rootCmd = &cobra.Command{
	Use:   "Friday",
	Short: "Get SQL query in everyday language",
	Long:  "This is Long Description",
	Run: func(cmd *cobra.Command, args []string) {

		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("something is wrong")
		}

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Something is wrong!!", err)
		os.Exit(1)
	}
}
