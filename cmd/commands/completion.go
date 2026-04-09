package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	root "vault/cmd"
)

var completionShells = []string{"bash", "zsh", "fish", "powershell", "cmd"}

var CompletionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell|cmd]",
	Short: "Generate shell completion script",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			autoInstall()
			return
		}

		shell := args[0]
		switch shell {
		case "bash":
			genBashCompletion()
		case "zsh":
			genZshCompletion()
		case "fish":
			genFishCompletion()
		case "powershell":
			genPowerShellCompletion()
		case "cmd":
			genCmdCompletion()
		default:
			fmt.Printf("Error: unsupported shell '%s'\n", shell)
			os.Exit(1)
		}
	},
}

func autoInstall() {
	shell := detectShell()
	switch shell {
	case "bash":
		installBashCompletion()
	case "zsh":
		installZshCompletion()
	case "fish":
		installFishCompletion()
	default:
		fmt.Println("Error: unsupported shell. Use: vault completion [bash|zsh|fish]")
		os.Exit(1)
	}
}

func detectShell() string {
	pid := os.Getpid()
	for i := 0; i < 10; i++ {
		cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", pid), "-o", "comm=")
		output, _ := cmd.Output()
		name := string(output)
		name = name[:len(name)-1]

		if name == "bash" || name == "zsh" || name == "fish" {
			return name
		}

		cmd = exec.Command("ps", "-o", "ppid=", "-p", fmt.Sprintf("%d", pid))
		output, _ = cmd.Output()
		var newPid int
		fmt.Sscanf(string(output), "%d", &newPid)
		if newPid == 0 || newPid == 1 {
			break
		}
		pid = newPid
	}
	return ""
}

func genBashCompletion() {
	root.RootCmd.GenBashCompletion(os.Stdout)
}

func genZshCompletion() {
	root.RootCmd.GenZshCompletion(os.Stdout)
}

func genFishCompletion() {
	root.RootCmd.GenFishCompletion(os.Stdout, true)
}

func genPowerShellCompletion() {
	root.RootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
}

func genCmdCompletion() {
	fmt.Println("CMD completion: Use 'vault completion powershell' in PowerShell")
}

func installBashCompletion() {
	fmt.Println("Installing Bash completion...")
	home, _ := os.UserHomeDir()

	exec.Command("mkdir", "-p", home+"/.bash_completion.d").Run()

	out, _ := os.Create(home + "/.bash_completion.d/vault")
	root.RootCmd.GenBashCompletion(out)
	out.Close()

	fmt.Println("Installed to ~/.bash_completion.d/vault")
	fmt.Println("Add to ~/.bashrc: source ~/.bash_completion.d/vault")
}

func installZshCompletion() {
	fmt.Println("Installing Zsh completion...")
	home, _ := os.UserHomeDir()

	exec.Command("mkdir", "-p", home+"/.zfunc").Run()

	out, _ := os.Create(home + "/.zfunc/_vault")
	root.RootCmd.GenZshCompletion(out)
	out.Close()

	fmt.Println("Installed to ~/.zfunc/_vault")
	fmt.Println("Add to ~/.zshrc: fpath+=(~/.zfunc) && autoload -Uz compinit && compinit")
}

func installFishCompletion() {
	fmt.Println("Installing Fish completion...")
	home, _ := os.UserHomeDir()

	exec.Command("mkdir", "-p", home+"/.config/fish/completions").Run()

	out, _ := os.Create(home + "/.config/fish/completions/vault.fish")
	root.RootCmd.GenFishCompletion(out, true)
	out.Close()

	fmt.Println("Installed to ~/.config/fish/completions/vault.fish")
}
