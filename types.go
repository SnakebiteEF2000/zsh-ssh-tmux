package main

type StatusCode = int

const (
	StatusOK = iota
	StatusErrReadingInventory
	StatusErrParsingInventory
	StatusErrWritingSSHConfig
	StatusErrUserHomeDir
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
