package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/kirsle/configdir"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type ConfigDb struct {
	Url   string `yaml:"sonar_url" env:"SONAR_URL"`
	Token string `yaml:"sonar_token" env:"SONAR_TOKEN"`
}

func main() {
	// Make sure config dir exists
	configPath := configdir.LocalConfig("sonar_request")
	err := configdir.MakePath(configPath)
	if err != nil {
		panic(err)
	}

	// Set config file path
	configFile := filepath.Join(configPath, "config.yml")

	// Check if exists
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Error: no config.yml file at " + configPath)
		return
	}

	var cfg ConfigDb

	// Read config
	err = cleanenv.ReadConfig(configFile, &cfg)
	if err != nil {
		panic(err)
	}

	// Set variables from config
	url := cfg.Url + "/api/graphql"
	token := cfg.Token

	// Check for graphql file argument
	if len(os.Args) < 2 {
		fmt.Println("Missing graphql file")
		return
	}

	// Read graphql file
	graphql, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Create POST body
	data := map[string]string{"query": string(graphql)}

	// Make JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// Create POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// Set custom headers for auth & JSON
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	// Create client
	client := &http.Client{}

	// Make request
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Close on finish
	defer res.Body.Close()

	// Read response body
	body, _ := io.ReadAll(res.Body)

	// Output response
	fmt.Println(string(body))
}
