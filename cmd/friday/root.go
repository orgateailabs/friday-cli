package friday

import (
	"fmt"
	"os"
	"strings"

	"github.com/Orgate-AI/friday-cli/utils"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	apiKey         = os.Getenv("ORGATEAI_API_KEY")
	schemaFilePath = os.Getenv("DB_SCHEMA_FILE_PATH")
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	senderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	botStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
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
	spinner     spinner.Model
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
	spin := spinner.New(spinner.WithSpinner(spinner.Points))
	ti.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ti,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		spinner:     spin,
		err:         nil,
	}
}

// take the JSON file and form a string
// func getDBSchema(schemaFilePath string) string {
// 	finalSchema := ""
// 	contents, err := os.ReadFile(schemaFilePath)
// 	if err != nil {
// 		fmt.Println("File reading error:", err)
// 		return "File reading error. Check log."
// 	}
// 	// schema := string(contents)
// 	fmt.Println(contents)

// 	return finalSchema

// }

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
		footerStyle = lipgloss.NewStyle().
			Height(1).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("8")).
			Faint(true)
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			query := strings.TrimSpace(m.textarea.Value())
			m.messages = append(m.messages, query)
			m.viewport.SetContent(m.RenderConversation(m.viewport.Width))
			m.textarea.Reset()
			m.textarea.Blur()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) RenderConversation(maxWidth int) string {
	var sb strings.Builder

	dbSchema := "Employee(id, name, department_id)\n# Department(id, name, address)\n# Salary_Payments(id, employee_id, amount, date)"
	apiKey := "orai-e423b59f-e915-41d8-a173-6411ea9b4c88"
	query := m.messages[0]

	if query == "" {
		return ""
	}

	renderYou := func(content string) {
		sb.WriteString(senderStyle.Render("You: "))
		content, _ = glamour.Render(content, "dark")
		sb.WriteString(content)
	}
	renderBot := func(content string) {
		if content == "" {
			return
		}
		sb.WriteString(botStyle.Render("Friday: "))
		content, _ = glamour.Render(content, "dark")
		sb.WriteString(content)
	}
	renderYou(query)
	sqlQuery, err := utils.RunQuery(query, apiKey, dbSchema)
	if err != nil {
		fmt.Println(err)
	}
	renderBot(sqlQuery)

	return sb.String()
}

func (m model) RenderFooter() string {
	var columns []string

	columns = append(columns, fmt.Sprintf("To quit: %s", "ctrl + c"))

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
		fmt.Println("Something is wrong!!", err)
		os.Exit(1)
	}
}
