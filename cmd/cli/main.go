package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"zoubun/internal/routes"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	cli "github.com/urfave/cli/v3"
)

var endpoint = "http://localhost:3000"

func main() {
	motdCmd := &cli.Command{
		Name:   "motd",
		Usage:  "Display the server's message of the day",
		Action: motd,
	}

	incCmd := &cli.Command{
		Name:   "increment",
		Usage:  "increment your count by 1",
		Action: increment,
	}
	countCmd := &cli.Command{
		Name:   "count",
		Usage:  "Displays your current count",
		Action: count,
	}

	tuiCmd := &cli.Command{
		Name:   "tui",
		Usage:  "Displays a terminal ui",
		Action: tui,
	}
	cmd := &cli.Command{
		Name:    "zbcli",
		Version: "v0.0.1",
		Usage:   "zbcli [command]",
		Commands: []*cli.Command{
			motdCmd, incCmd, countCmd, tuiCmd,
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func motd(ctx context.Context, cmd *cli.Command) error {
	resp, err := http.Get(endpoint + "/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var jsonOutput routes.Motd

	err = json.NewDecoder(resp.Body).Decode(&jsonOutput)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.Root().Writer, "Message of the day: %v", jsonOutput.Message)
	return nil
}

func count(ctx context.Context, cmd *cli.Command) error {
	resp, err := http.Get(endpoint + "/count")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var jsonOutput routes.Counter

	err = json.NewDecoder(resp.Body).Decode(&jsonOutput)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.Root().Writer, "Count %v", jsonOutput.Count)
	return nil
}

func increment(ctx context.Context, cmd *cli.Command) error {
	jsonData, err := json.Marshal("{}")
	if err != nil {
		return err
	}
	resp, err := http.Post(endpoint+"/increment", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var jsonOutput routes.Counter

	err = json.NewDecoder(resp.Body).Decode(&jsonOutput)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.Root().Writer, "Incremented to %v", jsonOutput.Count)
	return nil
}

// Bubbletea based tui just to test and get the package.
// TODO: Implement fully once the commands are refactored.
type (
	errMsg error
	model  struct {
		spinner  spinner.Model
		quitting bool
		err      error
	}
)

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n	%s Loading forever...press q to quit\n\n", m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	return str
}

func tui(ctx context.Context, cmd *cli.Command) error {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
