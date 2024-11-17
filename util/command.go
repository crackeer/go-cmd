package util

import (
	"context"
	"os"
	"os/exec"
	"time"
)

// RunCommand
//
//	@param ctx
//	@param cmd
//	@return error
func RunCommand(ctx context.Context, cmd *exec.Cmd) error {

	// start the command; if it fails to start, report error immediately
	err := cmd.Start()
	if err != nil {
		return err
	}

	// wait for the command in a goroutine; the reason for this is
	// very subtle: if, in our select, we do `case cmdErr := <-cmd.Wait()`,
	// then that case would be chosen immediately, because cmd.Wait() is
	// immediately available (even though it blocks for potentially a long
	// time, it can be evaluated immediately). So we have to remove that
	// evaluation from the `case` statement.
	cmdErrChan := make(chan error)
	go func() {
		cmdErrChan <- cmd.Wait()
	}()

	// unblock either when the command finishes, or when the done
	// channel is closed -- whichever comes first
	select {
	case cmdErr := <-cmdErrChan:
		// process ended; report any error immediately
		return cmdErr
	case <-ctx.Done():
		// context was canceled, either due to timeout or
		// maybe a signal from higher up canceled the parent
		// context; presumably, the OS also sent the signal
		// to the child process, so wait for it to die
		select {
		case <-time.After(15 * time.Second):
			_ = cmd.Process.Kill()
		case <-cmdErrChan:
		}
		return ctx.Err()
	}
}

// RunXXD ...
//
//	@param file
//	@param args
//	@return []byte
//	@return error
func RunXXD(file string, args ...string) ([]byte, error) {
	args = append(args, file)
	cmd := exec.CommandContext(context.Background(), "xxd", args...)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

// RunDD @param file
//
//	@param args
//	@return []byte
//	@return error
func RunDD(args ...string) ([]byte, error) {
	cmd := exec.CommandContext(context.Background(), "dd", args...)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

// Mount
//
//	@param args
//	@return error
func Mount(args ...string) error {
	cmd := exec.Command("mount", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// UnMount
//
//	@param node
//	@return error
func UnMount(node string) error {
	cmd := exec.Command("umount", node)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// execShCommand
//
//	@param workDir
//	@param shFile
//	@param envs
//	@return error
func ExecShCommand(workDir string, shFile string, envs []string) error {
	cmd := exec.Command("bash", shFile)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = workDir
	cmd.Env = envs
	return cmd.Run()
}

// execGitCloneCommand
//
//	@param workDir
//	@param gitURL
//	@param branch
//	@param codeFolder
//	@return error
func ExecGitCloneCommand(workDir string, gitURL string, branch string, codeFolder string) error {
	cmd := exec.Command("git", "clone", "--depth=1", "--branch="+branch, gitURL, codeFolder)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = workDir
	return cmd.Run()
}

// ExecDockerSave ...
//
//	@param imageName
//	@param dest
//	@return error
func ExecDockerSave(imageName string, dest string) error {
	cmd := exec.Command("docker", "save", imageName, "-o", dest)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
