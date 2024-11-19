package command

import (
	"os"
	"path/filepath"

	"github.com/crackeer/go-cmd/util"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

func NewZip(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(2),
		Run:   doZip,
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", true, "verbose")
	cmd.SetHelpTemplate(`./got zip /path/to/zip file.zip	
`)
	return cmd
}

func NewUnzip(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(1),
		Run:   doUnzip,
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", true, "verbose")
	cmd.SetHelpTemplate(`./got unzip file.zip	
`)
	return cmd
}

// doZip compresses the source directory or file specified by args[0]
// into a zip file at the destination specified by args[1].
// If verbose is true, it prints the directory of the destination path.
func doZip(cmd *cobra.Command, args []string) {
	if filename, err := filepath.Rel(args[0], args[1]); err == nil && len(filename) > 0 {
		if f, err := os.Stat(args[1]); err == nil && !f.IsDir() {
			panic(args[1] + " exists in " + args[0])
		}
	}
	err := util.Zip(args[0], args[1], verbose)
	if err != nil {
		panic(err)
	}
}

func doUnzip(cmd *cobra.Command, args []string) {
	err := util.Unzip(args[0], args[1], verbose)
	if err != nil {
		panic(err)
	}
}
