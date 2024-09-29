package setup

import (
	"flag"
	"fmt"
	"os"
)

type ExecuteConfig struct {
	InventoryPath   *string
	SSHUser         *string
	AltSSHUser      *string
	AltSSHUserRegex *string
	HomeDir         string
}

func (c *ExecuteConfig) Init() error {
	c.InventoryPath = flag.String("inventory", "hosts.yml", "Path to the inventory file")
	c.SSHUser = flag.String("user", "", "SSH user for hosts")
	c.AltSSHUser = flag.String("altuser", "", "Alternative user for connects")
	c.AltSSHUserRegex = flag.String("altuserregex", "", "Use Match Host <altuserregex> for user matching. See man ssh_config.")
	flag.Parse()

	if flag.NArg() > 0 {
		*c.InventoryPath = flag.Arg(0)
	}

	if *c.SSHUser == "" {
		fmt.Print("Enter your admin user used for SSH on non DMZ-Hosts: ")
		fmt.Scanln(c.SSHUser) // fail needs handle
	}

	var err error
	c.HomeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Printf("cant get user home dir, err: %W", err)
		return err
	}
	return nil
}
