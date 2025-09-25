// Package term contains all the functionality required of the cli program
package term

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"zoubun/internal/routes"

	"github.com/urfave/cli/v3"
)

var endpoint = "http://localhost"

type CLIConfig struct {
	Username string `json:"username"`
	APIKey1  string `json:"apikey1"`
	APIKey2  string `json:"apikey2"`
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(home, ".config", "zoubun.cfg")
}

func (c *CLIConfig) Load() error {
	cfgpath := getConfigPath()

	log.Printf("Path to config is: %v", cfgpath)
	if _, err := os.Stat(cfgpath); errors.Is(err, os.ErrNotExist) {
		log.Printf("File not found creating new file at %v", cfgpath)
		file, err := os.Create(cfgpath)
		if err != nil {
			return err
		}
		rawjson, err := json.Marshal(c)
		if err != nil {
			return err
		}
		_, err = file.Write(rawjson)
		if err != nil {
			return err
		}
	} else {
		file, err := os.ReadFile(cfgpath)
		if err != nil {
			return err
		}
		return json.Unmarshal(file, c)
	}
	return nil
}

func (c *CLIConfig) Save() error {
	cfgpath := getConfigPath()
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(cfgpath, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	return err
}

func CLICommands() *cli.Command {
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
	register := &cli.Command{
		Name:  "register",
		Usage: "Registers an account with a given `username`",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "username",
				UsageText: "Should be a valid username of 3 or more characters. This also must be unique to the service.",
			},
		},
		Action: register,
	}
	tuiCmd := &cli.Command{
		Name:   "tui",
		Usage:  "Displays a terminal ui",
		Action: tui,
	}
	return &cli.Command{
		Name:    "zbcli",
		Version: "v0.0.1",
		Usage:   "zbcli [command]",
		Commands: []*cli.Command{
			motdCmd, register, incCmd, countCmd, tuiCmd,
		},
	}
}

func register(ctx context.Context, cmd *cli.Command) error {
	cfg := &CLIConfig{}
	cfg.Load()
	if cfg.Username != "" || cfg.APIKey1 != "" || cfg.APIKey2 != "" {
		return fmt.Errorf("you are already registered as %v", cfg.Username)
	}
	reqbody := routes.RegisterRequestBody{
		Username: cmd.StringArg("username"),
	}
	log.Printf("Attempting to register user %v", reqbody.Username)
	jsonvalue, err := json.Marshal(reqbody)
	if err != nil {
		return err
	}
	resp, err := http.Post(endpoint+"/register", "application/json", bytes.NewReader(jsonvalue))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var output routes.RegisterResponseBody
	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return err
	}
	config := CLIConfig{
		Username: output.Username,
		APIKey1:  output.APIKey1,
		APIKey2:  output.APIKey2,
	}
	log.Printf("Created username and saved keys locally %v", output)
	return config.Save()
}

func motd(ctx context.Context, cmd *cli.Command) error {
	resp, err := http.Get(endpoint + "/motd")
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
	cfg := &CLIConfig{}
	err := cfg.Load()
	if err != nil {
		return err
	}
	if cfg.APIKey1 == "" || cfg.APIKey2 == "" {
		return errors.New("no valid API keys in config")
	}
	req, err := http.NewRequest(http.MethodGet, endpoint+"/count", nil)
	if err != nil {
		return err
	}
	req.Header.Add("zoubun-api-key", cfg.APIKey1)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New("resp there was a server error")
	}

	defer resp.Body.Close()
	var jsonOutput routes.Counter

	err = json.NewDecoder(resp.Body).Decode(&jsonOutput)
	if err != nil {
		return err
	}
	log.Printf("Your current count is %v", jsonOutput.Count)
	return nil
}

func increment(ctx context.Context, cmd *cli.Command) error {
	cfg := &CLIConfig{}
	err := cfg.Load()
	if err != nil {
		return err
	}

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

	req.Header.Add("zoubun-api-key", cfg.APIKey1)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New("resp there was a server error")
	}

	defer resp.Body.Close()
	var jsonOutput routes.Counter

	err = json.NewDecoder(resp.Body).Decode(&jsonOutput)
	if err != nil {
		return err
	}

	log.Printf("Incremented to %v", jsonOutput.Count)
	return nil
}
