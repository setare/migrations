package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jamillosantos/migrations"
	color "gopkg.in/gookit/color.v1"
)

var (
	StyleH1                     = color.New(color.Bold)
	StyleSuccess                = color.New(color.FgGreen)
	StyleWarning                = color.New(color.FgYellow)
	StyleError                  = color.New(color.FgRed)
	StyleErrorHighlight         = color.New(color.BgRed, color.FgWhite)
	StyleMigrationID            = color.New(color.FgCyan)
	StyleMigrationTitle         = color.New()
	StyleMigrationCurrentMarker = color.New(color.Green.Darken(), color.Bold)
	StyleMigrationDo            = color.New(color.BgGreen)
	StyleMigrationUndo          = color.New(color.BgCyan)
)

type Printer struct {
	prefix    string
	spinner   *spinner.Spinner
	startedAt time.Time
}

var Output Printer

func IconOK() string {
	return StyleSuccess.Sprint("✔")
}

func IconError() string {
	return StyleError.Sprint("✗")
}

func IconPending() string {
	return StyleWarning.Sprint("⧖")
}

func (p *Printer) H1(args ...interface{}) {
	fmt.Print(p.prefix)
	StyleH1.Println(args...)
}

func (p *Printer) H1f(format string, args ...interface{}) {
	fmt.Print(p.prefix)
	StyleH1.Printf(format, args...)
	fmt.Println()
}

func (p *Printer) Smigration(migration migrations.Migration) string {
	return fmt.Sprintf("%s %s", StyleMigrationID.Sprint(migration.ID().Format(migrations.DefaultMigrationIDFormat)), StyleMigrationTitle.Sprint(migration.Description()))
}

func (p *Printer) Saction(actionType migrations.ActionType) string {
	if actionType == migrations.ActionTypeUndo {
		return StyleMigrationUndo.Sprint(" rewind ")
	}
	return StyleMigrationDo.Sprint(" migrate ")
}

func (p *Printer) PlanAction(i int, icon string, action *migrations.Action) {
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

func (p *Printer) MigrationItemList(i int, icon string, migration migrations.Migration, current bool) {
	fmt.Print(p.prefix)
	suffix := ""
	if current {
		suffix = StyleMigrationCurrentMarker.Sprint(" <<< ")
	}
	fmt.Printf("  %s %s%s", icon, p.Smigration(migration), suffix)
	fmt.Println()
}

func (p *Printer) Success(args ...interface{}) {
	fmt.Print(p.prefix)
	color.Success.Println(args...)
}

func (p *Printer) Warn(args ...interface{}) {
	fmt.Print(p.prefix)
	color.Warn.Println(args...)
}

func (p *Printer) Error(args ...interface{}) {
	fmt.Print(p.prefix)
	color.Error.Println(args...)
}

func (p *Printer) Errorf(format string, args ...interface{}) {
	fmt.Print(p.prefix)
	color.Error.Printf(format, args...)
	color.Error.Println()
}

func (p *Printer) Print(args ...interface{}) {
	fmt.Print(p.prefix)
	fmt.Println(args...)
}

func (p *Printer) Printf(format string, args ...interface{}) {
	fmt.Print(p.prefix)
	fmt.Printf(format, args...)
	fmt.Println()
}

func (p *Printer) BeforeExecute(actionType migrations.ActionType, migration migrations.Migration) {
	fmt.Print(p.prefix)
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

func (p *Printer) AfterExecute(actionType migrations.ActionType, migration migrations.Migration, err error) {
	duration := time.Since(p.startedAt)
	fmt.Print(p.prefix)
	p.spinner.Stop()
	fmt.Print("  ")
	if err != nil {
		qryError, ok := err.(migrations.QueryError)
		fmt.Printf("%s %s %s %s: %s (%s)", StyleError.Sprint(IconError()), StyleErrorHighlight.Sprint("FAILED"), p.Smigration(migration), p.Saction(actionType), StyleError.Sprint(err), duration)
		if ok {
			fmt.Println()
			fmt.Println()
			fmt.Print(p.prefix)
			fmt.Print("      Query: ", qryError.Query())
		} else {
		}
		fmt.Println()
		return
	}
	fmt.Printf("%s %s %s (%s)", StyleSuccess.Sprint(IconOK()), p.Smigration(migration), p.Saction(actionType), duration)
	fmt.Println()
}
