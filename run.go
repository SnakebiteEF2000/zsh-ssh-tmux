package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SnakebiteEF2000/zsh-ssh-tmux/internal/setup"
	"gopkg.in/yaml.v2"
)

func execute() StatusCode {
	var cfg = new(setup.ExecuteConfig)
	err := cfg.Init()
	if err != nil {
		return StatusErrUserHomeDir
	}

	inventoryData, err := os.ReadFile(*cfg.InventoryPath)
	if err != nil {
		fmt.Printf("Error reading inventory file: %v\n", err)
		return StatusErrReadingInventory
	}

	var inventory InventoryData
	err = yaml.Unmarshal(inventoryData, &inventory)
	if err != nil {
		fmt.Printf("Error parsing inventory YAML: %v\n", err)
		return StatusErrParsingInventory
	}

	var sshConfig strings.Builder

	sshConfig.WriteString("# SSH Config generated from inventory\n\n")
	if *cfg.AltSSHUserRegex != "" && *cfg.AltSSHUser != "" {
		sshConfig.WriteString(fmt.Sprintf("Match Host %s\n", *cfg.AltSSHUserRegex))
		sshConfig.WriteString(fmt.Sprintf("  User %s\n\n", *cfg.AltSSHUser))
	}

	sshConfig.WriteString("Match Host *\n")
	sshConfig.WriteString(fmt.Sprintf("  User %s\n\n", *cfg.SSHUser))

	for groupName, hosts := range inventory.All.Children {
		fmt.Printf("Processing group: %s\n", groupName)
		for hostname, host := range hosts.Hosts {
			if len(host.AnsibleHost) > 1 {
				host.writeHostConfig(&sshConfig, hostname)
			}
		}
	}

	fmt.Println("Generated SSH config:") // missing format?

	sshConfigPath := filepath.Join(cfg.HomeDir, ".ssh", "ansible-inventory")
	err = os.WriteFile(sshConfigPath, []byte(sshConfig.String()), 0600)
	if err != nil {
		fmt.Printf("Error writing SSH config file: %v\n", err)
		return StatusErrWritingSSHConfig
	}

	fmt.Printf("SSH config file generated at: %s\n", sshConfigPath)

	return StatusOK
}

func (h *Host) writeHostConfig(sshConfig *strings.Builder, hostname string) {
	sshConfig.WriteString(fmt.Sprintf("Host %s %s\n", hostname, h.getDescription()))
	sshConfig.WriteString(fmt.Sprintf("  HostName %s\n", h.getHostname(hostname)))

	sshConfig.WriteString("\n")
}

func (h *Host) getDescription() string {
	var description strings.Builder
	if role, ok := h.CustomFields["HW_role_services_description"].(string); ok && role != "" {
		description.WriteString(fmt.Sprintf(" %s", role))
	}

	if len(h.Tags) > 0 {
		description.WriteString(fmt.Sprintf(" %s", strings.Join(h.Tags, " ")))
	}

	return description.String()
}

func (h *Host) getHostname(defaultHost string) string {
	if h.AnsibleHost != "" {
		return h.AnsibleHost
	}
	return defaultHost
}
