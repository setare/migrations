package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/setare/migrations"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use: "create <type:sql|source> <description>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("migration description required")
		}
		return nil
	},
	Short: "Create a migration file",
	Long:  `This command performs a "rewind" command followed by a "migrate".`,
}

func init() {
	rootCmd.AddCommand(createCmd)

	var (
		noUndoFlag bool
	)

	createCmd.Flags().BoolVarP(&noUndoFlag, "no-undo", "n", false, `Do not generate the "undo" step for the migration`)

	createCmd.Run = func(cmd *cobra.Command, args []string) {
		migrationNamePrefix := fmt.Sprintf("%s_%s", time.Now().UTC().Format(migrations.DefaultMigrationIDFormat), strings.Join(args, "_"))

		dir := viper.GetString("directory")

		migrationsDirStats, err := os.Stat(dir)
		if errors.Is(err, os.ErrNotExist) {
			createFolder := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("The folder %s does not exists. Shall I create it?", dir),
				Default: true,
			}
			err = survey.AskOne(prompt, &createFolder)
			if err != nil {
				output.Error("unknown error: ", err)
				os.Exit(1)
			}
			if createFolder {
				err = os.MkdirAll(dir, os.ModePerm)
				if err != nil {
					output.Errorf("could not create directory %s: %s", dir, err)
					os.Exit(206)
				}
			} else {
				output.Warn("user cancelled operation")
				os.Exit(EC_CANCELLED)
			}
		} else if err == nil && !migrationsDirStats.IsDir() {
			output.Errorf("%s is not a directory")
			os.Exit(205)
		}

		// if migrationType == "sql" {
		fileNameDo := path.Join(dir, migrationNamePrefix+".do.sql")

		fileInfoDo, err := os.Stat(fileNameDo)
		if !errors.Is(err, os.ErrNotExist) {
			output.Error("failed checking file: ", err)
			os.Exit(200)
		} else if err == nil && fileInfoDo.IsDir() {
			output.Errorf("there is a directory with the same migration name (%s): %s", fileNameDo, err)
			os.Exit(201)
		} else if err == nil {
			output.Errorf("migration file already exists (%s): ", fileNameDo, err)
			os.Exit(203)
		}

		fDo, err := os.Create(fileNameDo)
		if err != nil {
			output.Errorf("error creating file: %s", fileNameDo)
			os.Exit(204)
		}
		defer fDo.Close()

		fDo.WriteString("-- SQL here")

		output.Printf("%s created", styleSuccess.Sprint(fileNameDo))

		if noUndoFlag {
			// Do not generate undo file
			return
		}

		fileNameUndo := path.Join(dir, migrationNamePrefix+".undo.sql")

		fileInfoUndo, err := os.Stat(fileNameUndo)
		if !errors.Is(err, os.ErrNotExist) {
			output.Error("failed checking file: ", err)
			os.Exit(200)
		} else if err == nil && fileInfoUndo.IsDir() {
			output.Errorf("there is a directory with the same migration name (%s): %s", fileNameUndo, err)
			os.Exit(201)
		} else if err == nil {
			output.Errorf("migration file already exists (%s): ", fileNameUndo, err)
			os.Exit(203)
		}

		fUndo, err := os.Create(fileNameUndo)
		if err != nil {
			output.Errorf("error creating file: %s", fileNameUndo)
			os.Exit(204)
		}
		defer fUndo.Close()

		fUndo.WriteString("-- SQL here")
		output.Printf("%s created", styleSuccess.Sprint(fileNameUndo))
		// }
	}
}
