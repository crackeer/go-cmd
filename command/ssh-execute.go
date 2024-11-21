package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/crackeer/go-cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type SSHExecuteConfig struct {
	Servers  map[string]string   `json:"servers"`
	Commands map[string][]string `json:"commands"`
}

func NewSSHExecute(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run:   doSSHExecute,
		Args:  cobra.MinimumNArgs(0),
	}
	cmd.SetUsageTemplate(`got ssh-execute host_aliais command
`)

	return cmd
}

func doSSHExecute(cmd *cobra.Command, args []string) {
	viper.AutomaticEnv()
	sshConfigFile := "/usr/local/etc/got/ssh-config.json"
	if value := viper.Get("SSH_EXEC_CONFIG"); value != nil {
		sshConfigFile = value.(string)
	}
	sshConfig := SSHExecuteConfig{}
	if err := util.ReadFileAs(sshConfigFile, &sshConfig); err != nil {
		fmt.Println("read ssh config file error: ", err)
		return
	}
	// 1. list hosts
	if len(args) < 1 {
		fmt.Println("available hosts:")
		for name, host := range sshConfig.Servers {
			fmt.Println("host:", name, " -> ", host)
			fmt.Println("")
		}
		return
	}

	// 2. verify host config
	hostConfig, ok := sshConfig.Servers[args[0]]
	if !ok {
		fmt.Println("host not found")
		return
	}
	parts := strings.Split(hostConfig, ":")
	if len(parts) < 4 {
		fmt.Println("host config error")
		return
	}

	// 3. list commands
	if len(args) < 2 {
		fmt.Println("available commands:")
		for name, cmds := range sshConfig.Commands {
			fmt.Println("command_alias:", name)
			fmt.Println("commands:", strings.Join(cmds, "   "))
			fmt.Println("")
		}
		return
	}

	// 4. connect remote server
	user := parts[0]
	password := parts[1]
	host := parts[2]
	port := parts[3]
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshClient, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		fmt.Println("dial error: ", err)
		return
	}

	// 5. execute commands

	for _, cmdAlias := range args[1:] {
		if _, ok := sshConfig.Commands[cmdAlias]; !ok {
			fmt.Println("command not found")
			return
		}

		fmt.Println("----- execute command alias: ", cmdAlias, " -----")

		commands := sshConfig.Commands[cmdAlias]

		for _, command := range commands {
			fmt.Println("command: ", command)
			session, err := sshClient.NewSession()
			if err != nil {
				fmt.Println("new session error: ", err)
				return
			}
			//session.RequestPty("bash", 80, 40, ssh.TerminalModes{})
			session.Stdout = os.Stdout
			session.Stderr = os.Stderr
			if err := session.Run(command); err != nil {
				session.Close()
				fmt.Println("run command error: ", err)
				return
			}
			session.Close()
		}

		fmt.Println("----- end command alias: ", cmdAlias, " -----")
		fmt.Println("")

	}

	fmt.Println("all commands done!!!")

}
