package term

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"zoubun/internal/routes"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v3"
)

// Bubbletea based tui just to test and get the package.
// TODO: Implement fully once the commands are refactored.
type (
	countMsg struct {
		statusCode int
		count      int
	}
	errMsg error
	model  struct {
		motd        string
		latestCount int
		quitting    bool
		err         error
	}
)

var style = lg.NewStyle().
	Width(80).
	MaxHeight(40).
	AlignHorizontal(lg.Center).
	Border(lg.NormalBorder())

var titlestyle = lg.NewStyle().
	AlignHorizontal(lg.Center).
	AlignVertical(lg.Top).
	Height(2).
	Border(lg.NormalBorder(), false, false, true, false).
	Bold(true).
	Inherit(style)

var counterstyle = lg.NewStyle().
	AlignHorizontal(lg.Center).
	Width(80)

var errorstyle = lg.NewStyle().
	AlignHorizontal(lg.Left).
	Foreground(lg.Color("#f00000")).
	UnsetHeight().
	UnsetBorderStyle()

var helpstyle = lg.NewStyle().
	Inherit(style).
	AlignHorizontal(lg.Left).
	UnsetHeight().
	Background(lg.Color("#000080")).
	UnsetBorderStyle()

var config = CLIConfig{}

func newModel() model {
	err := config.Load()
	if err != nil {
		panic(err)
	}
	return model{
		motd:        "",
		latestCount: 0,
		quitting:    false,
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return countRequest
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			return m, incrementRequest
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	case errMsg:
		m.err = msg
		return m, nil
	case countMsg:
		m.latestCount = msg.count
		return m, nil
	default:
		return m, nil
	}
}

func getCurrentHotkeys(configIssue bool) string {
	if configIssue {
		return "[q]uit"
	} else {
		return "[q]uit [i]ncrement"
	}
}

func (m model) View() string {
	var header string
	title := titlestyle.Render("Zoubun TUI (experimental)")
	noUsername := config.Username == ""
	noKey := config.APIKey1 == "" || config.APIKey2 == ""

	// General Message
	if noUsername || noKey {
		header = "Welcome! Please exit and run the command `zoubun register [your desired username]` in your terminal to start"
	} else {
		header = fmt.Sprintf("Welcome to Zoubun %v! Please send the commands in the help menu as you like", config.Username)
	}

	counter := fmt.Sprintf("Count: %v\n", m.latestCount)
	var errormsg string
	if m.err != nil {
		header = fmt.Sprintf("ERROR: %v\n", m.err.Error())
	}

	hotkeys := getCurrentHotkeys(noUsername || noKey)
	output := lg.JoinVertical(
		lg.Top,
		titlestyle.Render(title),
		header,
		counterstyle.Render(counter),
		errorstyle.Render(errormsg),
		helpstyle.Render(hotkeys),
	)
	return style.Render(output)
}

// TODO: Add a mini-runtime compatibility for the requests
// Need to include a queue to check the status of the requests.

func countRequest() tea.Msg {
	req, err := http.NewRequest(http.MethodGet, endpoint+"/count", nil)
	if err != nil {
		return errMsg(err)
	}
	req.Header.Add("zoubun-api-key", config.APIKey1)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errMsg(err)
	}

	if resp.StatusCode >= 400 {
		return errMsg(errors.New("resp there was a server error"))
	}

	defer resp.Body.Close()
	var jsonOutput routes.Counter

	err = json.NewDecoder(resp.Body).Decode(&jsonOutput)
	if err != nil {
		return errMsg(err)
	}
	return countMsg{statusCode: resp.StatusCode, count: int(jsonOutput.Count)}
}

func incrementRequest() tea.Msg {
	jsonData, err := json.Marshal("{}")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		endpoint+"/increment",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	req.Header.Add("zoubun-api-key", config.APIKey1)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errMsg(err)
	}

	if resp.StatusCode >= 400 {
		return errMsg(errors.New("resp there was a server error"))
	}

	defer resp.Body.Close()
	var jsonOutput routes.Counter

	err = json.NewDecoder(resp.Body).Decode(&jsonOutput)
	if err != nil {
		return errMsg(err)
	}

	return countMsg{
		statusCode: resp.StatusCode,
		count:      int(jsonOutput.Count),
	}
}

func tui(ctx context.Context, cmd *cli.Command) error {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
