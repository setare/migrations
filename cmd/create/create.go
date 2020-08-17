package create

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jamillosantos/migrations"
	"github.com/jamillosantos/migrations/cmd"
	"github.com/ory/viper"
	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use: "create <description>",
	Args: func(cobraCmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("migration description required")
		}
		return nil
	},
	Short: "Create a migration file",
	Long:  ``,
}

var DefaultTemplate func(io.Writer, *TemplateInput) = WriteMigration

func init() {
	var (
		typeFlag   string
		noUndoFlag bool
	)

	CreateCmd.Flags().StringVarP(&typeFlag, "type", "t", "code", `Define the type of the migration. (code or sql)`)
	CreateCmd.Flags().BoolVarP(&noUndoFlag, "no-undo", "n", false, `Do not generate the "undo" step for the migration`)

	CreateCmd.Run = func(cobraCmd *cobra.Command, args []string) {
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
				cmd.Output.Error("unknown error: ", err)
				os.Exit(1)
			}
			if createFolder {
				err = os.MkdirAll(dir, os.ModePerm)
				if err != nil {
					cmd.Output.Errorf("could not create directory %s: %s", dir, err)
					os.Exit(206)
				}
			} else {
				cmd.Output.Warn("user cancelled operation")
				os.Exit(cmd.EC_CANCELLED)
			}
		} else if err == nil && !migrationsDirStats.IsDir() {
			cmd.Output.Errorf("%s is not a directory", dir)
			os.Exit(205)
		}

		switch typeFlag {
		case "code":
			fileName := path.Join(dir, migrationNamePrefix+".go")

			fMigration, err := os.Create(fileName)
			if err != nil {
				cmd.Output.Errorf("error creating file: %s", fileName)
				os.Exit(204)
			}
			defer fMigration.Close()

			DefaultTemplate(fMigration, &TemplateInput{
				Package:  path.Base(dir),
				DontUndo: noUndoFlag,
			})

			cmd.Output.Printf("%s created", cmd.StyleSuccess.Sprint(fileName))
		case "sql":
			fileNameDo := path.Join(dir, migrationNamePrefix+".do.sql")

			fileInfoDo, err := os.Stat(fileNameDo)
			if !errors.Is(err, os.ErrNotExist) {
				cmd.Output.Error("failed checking file: ", err)
				os.Exit(200)
			} else if err == nil && fileInfoDo.IsDir() {
				cmd.Output.Errorf("there is a directory with the same migration name (%s): %s", fileNameDo, err)
				os.Exit(201)
			} else if err == nil {
				cmd.Output.Errorf("migration file already exists (%s): %s", fileNameDo, err)
				os.Exit(203)
			}

			fDo, err := os.Create(fileNameDo)
			if err != nil {
				cmd.Output.Errorf("error creating file: %s", fileNameDo)
				os.Exit(204)
			}
			defer fDo.Close()

			fDo.WriteString("-- SQL here")

			cmd.Output.Printf("%s created", cmd.StyleSuccess.Sprint(fileNameDo))

			if noUndoFlag {
				// Do not generate undo file
				return
			}

			fileNameUndo := path.Join(dir, migrationNamePrefix+".undo.sql")

			fileInfoUndo, err := os.Stat(fileNameUndo)
			if !errors.Is(err, os.ErrNotExist) {
				cmd.Output.Error("failed checking file: ", err)
				os.Exit(200)
			} else if err == nil && fileInfoUndo.IsDir() {
				cmd.Output.Errorf("there is a directory with the same migration name (%s): %s", fileNameUndo, err)
				os.Exit(201)
			} else if err == nil {
				cmd.Output.Errorf("migration file already exists (%s): %s", fileNameUndo, err)
				os.Exit(203)
			}

			fUndo, err := os.Create(fileNameUndo)
			if err != nil {
				cmd.Output.Errorf("error creating file: %s", fileNameUndo)
				os.Exit(204)
			}
			defer fUndo.Close()

			fUndo.WriteString("-- SQL here")
			cmd.Output.Printf("%s created", cmd.StyleSuccess.Sprint(fileNameUndo))
		}
	}
}
