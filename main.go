package main

import (
	"flag"
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
	inventoryPath := flag.String("inventory", "hosts.yml", "Path to the inventory file")
	sshUser := flag.String("user", "", "SSH user for hosts")
	altsshUser := flag.String("altuser", "", "Alternative user for connects")
	altsshUserRegex := flag.String("altuserregex", "", "Use Match Host <altuserregex> for user matching. See man ssh_config.")
	flag.Parse()

	if flag.NArg() > 0 {
		*inventoryPath = flag.Arg(0)
	}

	inventoryData, err := os.ReadFile(*inventoryPath)
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

	if *sshUser == "" {
		fmt.Print("Enter your admin user used for SSH on non DMZ-Hosts: ")
		fmt.Scanln(sshUser)
	}

	var sshConfig strings.Builder

	sshConfig.WriteString("# SSH Config generated from inventory\n\n")
	if *altsshUserRegex != "" && *altsshUser != "" {
		sshConfig.WriteString(fmt.Sprintf("Match Host %s\n", *altsshUserRegex))
		sshConfig.WriteString(fmt.Sprintf("  User %s\n\n", *altsshUser))
	}

	sshConfig.WriteString(fmt.Sprintf("Match Host *\n"))
	sshConfig.WriteString(fmt.Sprintf("  User %s\n\n", *sshUser))

	for groupName, hosts := range inventory.All.Children {
		fmt.Printf("Processing group: %s\n", groupName)
		for hostname, hostData := range hosts.Hosts {
			if len(hostData.AnsibleHost) > 1 {
				writeHostConfig(&sshConfig, hostname, hostData, groupName)
			}
		}	
	}

	fmt.Println("Generated SSH config:")

	homeDir, _ := os.UserHomeDir()
	sshConfigPath := filepath.Join(homeDir, ".ssh", "ansible-inventory")
	err = os.WriteFile(sshConfigPath, []byte(sshConfig.String()), 0600)
	if err != nil {
		fmt.Printf("Error writing SSH config file: %v\n", err)
		return
	}

	fmt.Printf("SSH config file generated at: %s\n", sshConfigPath)
}

func writeHostConfig(sshConfig *strings.Builder, hostname string, hostData Host, groupName string) {
	description := getDescription(hostData, groupName)

	sshConfig.WriteString(fmt.Sprintf("Host %s %s\n", hostname, description))
	sshConfig.WriteString(fmt.Sprintf("  HostName %s\n", getHostname(hostData.AnsibleHost, hostname)))

	sshConfig.WriteString("\n")
}

func getDescription(hostData Host, groupName string) string {
	var description strings.Builder
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
