package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/setare/migrations"
	"gopkg.in/gookit/color.v1"
)

var (
	styleH1                     = color.New(color.Bold)
	styleSuccess                = color.New(color.FgGreen)
	styleWarning                = color.New(color.FgYellow)
	styleError                  = color.New(color.FgRed)
	styleErrorHighlight         = color.New(color.BgRed, color.FgWhite)
	styleMigrationID            = color.New(color.FgCyan)
	styleMigrationTitle         = color.New()
	styleMigrationCurrentMarker = color.New(color.Green.Darken(), color.Bold)
	styleMigrationDo            = color.New(color.BgGreen)
	styleMigrationUndo          = color.New(color.BgCyan)
)

type printer struct {
	prefix    string
	spinner   *spinner.Spinner
	startedAt time.Time
}

var output printer

func iconOK() string {
	return styleSuccess.Sprint("✔")
}

func iconError() string {
	return styleError.Sprint("✗")
}

func iconPending() string {
	return styleWarning.Sprint("⧖")
}

func (p *printer) H1(args ...interface{}) {
	fmt.Print(p.prefix)
	styleH1.Println(args...)
}

func (p *printer) H1f(format string, args ...interface{}) {
	fmt.Print(p.prefix)
	styleH1.Printf(format, args...)
	fmt.Println()
}

func (p *printer) Smigration(migration migrations.Migration) string {
	return fmt.Sprintf("%s %s", styleMigrationID.Sprint(migration.ID().Format(migrations.DefaultMigrationIDFormat)), styleMigrationTitle.Sprint(migration.Description()))
}

func (p *printer) Saction(actionType migrations.ActionType) string {
	if actionType == migrations.ActionTypeUndo {
		return styleMigrationUndo.Sprint(" rewind ")
	}
	return styleMigrationDo.Sprint(" migrate ")
}

func (p *printer) PlanAction(i int, icon string, action *migrations.Action) {
	fmt.Print(p.prefix)
	var actionPrefix string
	if icon == "" {
		actionPrefix = p.Saction(action.Action)
	} else {
		actionPrefix = icon
	}
	fmt.Printf("  %s %s", actionPrefix, p.Smigration(action.Migration))
	fmt.Println()
}

func (p *printer) MigrationItemList(i int, icon string, migration migrations.Migration, current bool) {
	fmt.Print(p.prefix)
	suffix := ""
	if current {
		suffix = styleMigrationCurrentMarker.Sprint(" <<< ")
	}
	fmt.Printf("  %s %s%s", icon, p.Smigration(migration), suffix)
	fmt.Println()
}

func (p *printer) Success(args ...interface{}) {
	fmt.Print(p.prefix)
	color.Success.Println(args...)
}

func (p *printer) Warn(args ...interface{}) {
	fmt.Print(p.prefix)
	color.Warn.Println(args...)
}

func (p *printer) Error(args ...interface{}) {
	fmt.Print(p.prefix)
	color.Error.Println(args...)
}

func (p *printer) Errorf(format string, args ...interface{}) {
	fmt.Print(p.prefix)
	color.Error.Printf(format, args...)
	color.Error.Println()
}

func (p *printer) Print(args ...interface{}) {
	fmt.Print(p.prefix)
	fmt.Println(args...)
}

func (p *printer) Printf(format string, args ...interface{}) {
	fmt.Print(p.prefix)
	fmt.Printf(format, args...)
	fmt.Println()
}

func (p *printer) BeforeExecute(actionType migrations.ActionType, migration migrations.Migration) {
	fmt.Print(output.prefix)
	p.spinner = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	p.spinner.Prefix = "  "
	p.spinner.Suffix = " "
	switch actionType {
	case migrations.ActionTypeDo:
		p.spinner.Suffix += "Migrating"
	case migrations.ActionTypeUndo:
		p.spinner.Suffix += "Reverting"
	}
	p.spinner.Suffix += fmt.Sprintf(" %s", p.Smigration(migration))
	p.startedAt = time.Now()
	p.spinner.Start()
}

func (p *printer) AfterExecute(actionType migrations.ActionType, migration migrations.Migration, err error) {
	duration := time.Since(p.startedAt)
	fmt.Print(output.prefix)
	p.spinner.Stop()
	fmt.Print("  ")
	if err != nil {
		qryError, ok := err.(migrations.QueryError)
		fmt.Printf("%s %s %s %s: %s (%s)", styleError.Sprint(iconError()), styleErrorHighlight.Sprint("FAILED"), p.Smigration(migration), p.Saction(actionType), styleError.Sprint(err), duration)
		if ok {
			fmt.Println()
			fmt.Println()
			fmt.Print(output.prefix)
			fmt.Print("      Query: ", qryError.Query())
		} else {
		}
		fmt.Println()
		return
	}
	fmt.Printf("%s %s %s (%s)", styleSuccess.Sprint(iconOK()), p.Smigration(migration), p.Saction(actionType), duration)
	fmt.Println()
}
