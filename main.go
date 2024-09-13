package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type InventoryData struct {
	All struct {
		Children map[string]Group `yaml:"children"`
	} `yaml:"all"`
}

type Group struct {
	Hosts map[string]Host `yaml:"hosts"`
}

type Host struct {
	AnsibleHost  string                 `yaml:"ansible_host"`
	CustomFields map[string]interface{} `yaml:"custom_fields"`
	Tags         []string               `yaml:"tags"`
}

func main() {
	inventoryPath := "hosts.yml"
	inventoryData, err := os.ReadFile(inventoryPath)
	if err != nil {
		fmt.Printf("Error reading inventory file: %v\n", err)
		return
	}

	var inventory InventoryData
	err = yaml.Unmarshal(inventoryData, &inventory)
	if err != nil {
		fmt.Printf("Error parsing inventory YAML: %v\n", err)
		return
	}

	fmt.Print("Enter your admin user used for SSH on non DMZ-Hosts: ")
	var sshUser string
	fmt.Scanln(&sshUser)

	var sshConfig strings.Builder
	sshConfig.WriteString("# SSH Config generated from inventory\n\n")

	for groupName, group := range inventory.All.Children {
		fmt.Printf("Processing group: %s\n", groupName)
		for hostname, hostData := range group.Hosts {
			writeHostConfig(&sshConfig, hostname, hostData, sshUser, groupName)
		}
	}

	fmt.Println("Generated SSH config:")
	fmt.Println(sshConfig.String())

	homeDir, _ := os.UserHomeDir()
	sshConfigPath := filepath.Join(homeDir, ".ssh", "config")
	err = os.WriteFile(sshConfigPath, []byte(sshConfig.String()), 0600)
	if err != nil {
		fmt.Printf("Error writing SSH config file: %v\n", err)
		return
	}

	fmt.Printf("SSH config file generated at: %s\n", sshConfigPath)
}

func writeHostConfig(sshConfig *strings.Builder, hostname string, hostData Host, sshUser, groupName string) {
	description := getDescription(hostData, groupName)
	sshConfig.WriteString(fmt.Sprintf("Host %s %s\n", hostname, description))
	sshConfig.WriteString(fmt.Sprintf("  HostName %s\n", getHostname(hostData.AnsibleHost, hostname)))

	if strings.Contains(strings.ToLower(hostname), "dmz") {
		sshConfig.WriteString("  User dummy\n")
	} else {
		sshConfig.WriteString(fmt.Sprintf("  User %s\n", sshUser))
	}

	sshConfig.WriteString("  IdentityFile ~/.ssh/id_rsa\n")
	sshConfig.WriteString("  AddKeysToAgent yes\n")
	sshConfig.WriteString("  SendEnv LANG LC_*\n")
	sshConfig.WriteString("  HashKnownHosts yes\n")
	sshConfig.WriteString("  GSSAPIAuthentication yes\n")
	sshConfig.WriteString("  GSSAPIDelegateCredentials no\n")
	sshConfig.WriteString("\n")
}

func getDescription(hostData Host, groupName string) string {
	var description strings.Builder
	description.WriteString(groupName)

	if role, ok := hostData.CustomFields["HW_role_services_description"].(string); ok && role != "" {
		description.WriteString(fmt.Sprintf(" %s", role))
	}

	if len(hostData.Tags) > 0 {
		description.WriteString(fmt.Sprintf(" %s", strings.Join(hostData.Tags, " ")))
	}

	return description.String()
}

func getHostname(ansibleHost, defaultHost string) string {
	if ansibleHost != "" {
		return ansibleHost
	}
	return defaultHost
}
