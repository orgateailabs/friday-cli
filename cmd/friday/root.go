package friday

import (
	"fmt"
	"os"
	"strings"

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

type (
	errMsg error
)

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
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
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
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

func (m model) View() string {
	// if m.width == 0 || m.height == 0 {
	// 	return "Initializing..."
	// }

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.viewport.View(),
		m.textarea.View(),
		// m.RenderFooter(),
	)
	// return fmt.Sprintf(
	// 	// inputStyle.Width(30).Render("You: "),
	// 	"%s\n\n%s",
	// 	m.viewport.View(),
	// 	m.textarea.View(),
	// 	// "(esc to quit)",
	// )
}

var rootCmd = &cobra.Command{
	Use:   "Friday",
	Short: "Get SQL query in everyday language",
	Long:  "This is Long Description",
	Run: func(cmd *cobra.Command, args []string) {

		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			// log.Fatal(err)
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
