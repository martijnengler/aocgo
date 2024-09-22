package resources

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"golang.org/x/term"
)

type LeaderboardModel struct {
	content  string
	viewport viewport.Model
	title    string
}

type ViewableLB interface {
	GetTitle() string
	GetContent() string
}

func NewLeaderboardViewport(content, title string) {
	m := LeaderboardModel{
		content: content,
		title:   title,
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Couldn't run viewport:", err)
		os.Exit(1)
	}
}

func (m LeaderboardModel) Init() tea.Cmd {
	log.Debug("'Init' function")

	return func() tea.Msg { return initMsg(0) }
}

func (m LeaderboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case initMsg:
		width, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return m, nil
		}

		headerHeight := lipgloss.Height(m.headerView())

		m.viewport = viewport.New(min(ViewportWidth, width), height-headerHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.HighPerformanceRendering = UseHighPerformanceRenderer
		m.viewport.SetContent(m.content)
		m.viewport.YPosition = headerHeight + 1

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())

		m.viewport.Width = min(ViewportWidth, msg.Width)
		m.viewport.Height = msg.Height - headerHeight

		if UseHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m LeaderboardModel) View() string {
	return fmt.Sprintf("%s\n%s", m.headerView(), m.viewport.View())
}

func (m LeaderboardModel) headerView() string {
	title := titleStyle.Render(m.title)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}