package friday

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	textInput textinput.Model
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "What do you wanna know about your data?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		inputStyle.Width(30).Render("You: "),
		m.textInput.View(),
		// "(esc to quit)",
	)
}

var rootCmd = &cobra.Command{
	Use: "Friday",
	Short: "Get SQL query in everyday language",
	Long: "This is Long Description",
	Run: func(cmd *cobra.Command, args []string){

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