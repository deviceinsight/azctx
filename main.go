package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Account struct {
	ID        string `json:"id"`
	IsDefault bool   `json:"isDefault"`
	Name      string `json:"name"`
}

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AzConfig struct {
	Defaults []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"defaults"`
}

var (
	azCli       = "az"
	GreenColor  = "\033[1;32m%s\033[0m"
	YellowColor = "\033[1;33m%s\033[0m"
)

func main() {
	selectedAccount := withFilter("Subscription", func(in io.WriteCloser) {
		accounts, err := getAllAccounts()
		if err != nil {
			in.Close()
			return
		}
		for _, account := range accounts {
			if account.IsDefault {
				fmt.Fprintf(in, GreenColor+"\n", account.Name)
			} else {
				fmt.Fprintln(in, account.Name)
			}
		}
	})

	if selectedAccount == "" {
		return
	}

	selectedGroup := withFilter("Resource group (Cancel for none)", func(in io.WriteCloser) {
		setAccount(selectedAccount)
		groups, err := getGroups()
		if err != nil {
			in.Close()
			return
		}
		defaultGroup := getDefaultGroup()
		for _, group := range groups {
			if group.Name == defaultGroup {
				fmt.Fprintf(in, GreenColor+"\n", group.Name)
			} else {
				fmt.Fprintln(in, group.Name)
			}
		}
		in.Close()
	})
	setGroup(selectedGroup)

	conclusion := "Switched to subscription " + fmt.Sprintf(GreenColor, selectedAccount)
	if selectedGroup != "" {
		conclusion += " and resource group " + fmt.Sprintf(GreenColor, selectedGroup)
	} else {
		conclusion += " and " + fmt.Sprintf(YellowColor, "no") + " resource group"
	}
	fmt.Println(conclusion)
}

func getAllAccounts() ([]Account, error) {
	output, execErr := exec.Command(azCli, "account", "list").Output()
	if execErr != nil {
		return nil, execErr
	}
	var accounts []Account
	jsonErr := json.Unmarshal(output, &accounts)
	if jsonErr != nil {
		return nil, execErr
	}
	return accounts, nil
}

func setAccount(accountName string) error {
	execErr := exec.Command(azCli, "account", "set", "--subscription", accountName).Run()
	if execErr != nil {
		return execErr
	}
	return nil
}

func getDefaultGroup() string {
	output, execErr := exec.Command(azCli, "config", "get").Output()
	if execErr != nil {
		return ""
	}
	var config AzConfig
	jsonErr := json.Unmarshal(output, &config)
	if jsonErr != nil {
		return ""
	}
	defaultGroup := ""
	for _, option := range config.Defaults {
		if option.Name == "group" {
			return option.Value
		}
	}
	return defaultGroup
}

func getGroups() ([]Group, error) {
	output, execErr := exec.Command(azCli, "group", "list").Output()
	if execErr != nil {
		return nil, execErr
	}
	var groups []Group
	jsonErr := json.Unmarshal(output, &groups)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return groups, nil
}

func setGroup(groupName string) {
	exec.Command(azCli, "config", "set", "defaults.group="+groupName).Run()
}

// https://junegunn.kr/2016/02/using-fzf-in-your-program/
func withFilter(header string, input func(in io.WriteCloser)) string {
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		shell = "sh"
	}
	cmd := exec.Command(shell, "-c", "fzf --ansi --layout=reverse --header=\""+header+"\"")
	cmd.Stderr = os.Stderr
	cmdIn, _ := cmd.StdinPipe()
	go func() {
		input(cmdIn)
		cmdIn.Close()
	}()
	result, _ := cmd.Output()
	return strings.ReplaceAll(string(result), "\n", "")
}
