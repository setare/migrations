package uiutils

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
)

func Spin(fnc func(*spinner.Spinner) error) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stdout))
	s.Start()
	defer s.Stop()
	return fnc(s)
}
