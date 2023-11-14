/*
Copyright © 2023 Ahmed Salah <ahmedsalah.yousuf@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func cpDir(src, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Join(dest, srcInfo.Name()), srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err = cpDir(path.Join(src, entry.Name()), path.Join(dest, entry.Name()))
			if err != nil {
				return err
			}
		} else {
			err = cpFileToDir(path.Join(src, entry.Name()), path.Join(dest, srcInfo.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func cpFileToDir(src, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(path.Join(dest, srcInfo.Name()))
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func cpFileToFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func cp(cmd *cobra.Command, args []string) {
	if len(args) == 2 {
		// copy {file | directory} to {file | directory}
		source := args[0]
		sourceInfo, err := os.Stat(source)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cp: cannot stat '%s': No such file or directory\n", source)
			os.Exit(1)
		}
		isSourceDir := sourceInfo.IsDir()

		dest := args[1]
		destInfo, err := os.Stat(dest)
		var isDestDir bool
		if err != nil {
			isDestDir = false
		} else {
			isDestDir = destInfo.IsDir()
		}

		if isSourceDir && isDestDir {
			isRecursive, err := cmd.Flags().GetBool("recursive")
			if err != nil {
				os.Exit(1)
			}
			if !isRecursive {
				fmt.Fprintf(os.Stderr, "cp: -r not specified; omitting directory '%s'\n", source)
				os.Exit(1)
			}

			err = cpDir(source, dest)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cp: cannot copy directory '%s' to directory '%s'\n", source, dest)
				os.Exit(1)
			}
		} else if !isSourceDir && !isDestDir {
			err := cpFileToFile(source, dest)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cp: cannot copy file '%s' to file '%s'\n", source, dest)
				os.Exit(1)
			}
		} else if !isSourceDir && isDestDir {
			err := cpFileToDir(source, dest)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cp: cannot copy file '%s' to directory '%s'\n", source, dest)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "cp: cannot copy directory '%s' to file '%s'\n", source, dest)
			os.Exit(1)
		}
	} else {
		// copy multiple files to directory
		dest := args[len(args)-1]
		destInfo, err := os.Stat(dest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cp: cannot stat '%s': No such file or directory\n", dest)
			os.Exit(1)
		}
		isDestDir := destInfo.IsDir()

		if !isDestDir {
			fmt.Fprintf(os.Stderr, "cp: target '%s' is not a directory\n", dest)
			os.Exit(1)
		}

		for _, source := range args[:len(args)-1] {
			sourceInfo, err := os.Stat(source)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cp: cannot stat '%s': No such file or directory\n", source)
				os.Exit(1)
			}
			isSourceDir := sourceInfo.IsDir()
			if isSourceDir {
				isRecursive, err := cmd.Flags().GetBool("recursive")
				if err != nil {
					os.Exit(1)
				}
				if !isRecursive {
					fmt.Fprintf(os.Stderr, "cp: -r not specified; omitting directory '%s'\n", source)
					continue
				}
				err = cpDir(source, dest)
				if err != nil {
					fmt.Fprintf(os.Stderr, "cp: cannot copy directory '%s' to directory '%s'\n", source, dest)
					os.Exit(1)
				}
			} else {
				err := cpFileToDir(source, dest)
				if err != nil {
					fmt.Fprintf(os.Stderr, "cp: cannot copy file '%s' to directory '%s'\n", source, dest)
					os.Exit(1)
				}
			}
		}
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cp {SOURCE DEST | SOURCE... DIRECTORY}",
	Short: "Copy SOURCE to DEST, or multiple SOURCE(s) to DIRECTORY.",
	Run:   cp,
	Args:  cobra.MinimumNArgs(2),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ch4.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("recursive", "r", false, "copy directories recursively")
}
