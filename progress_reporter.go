package migrations

import (
	"fmt"

	"github.com/gosuri/uiprogress"
)

type progressReporter struct {
	bar         *uiprogress.Bar
	steps       []string
	stepsMaxLen int
	currentStep int
	elapsedTime bool
}

func NewProgressReporter(progress *uiprogress.Progress) ProgressReporter {
	bar := progress.AddBar(0)
	reporter := &progressReporter{
		bar: bar,
	}
	bar.PrependFunc(reporter.prependFunc).AppendFunc(reporter.appendFunc)
	return reporter
}

func (reporter *progressReporter) prependFunc(b *uiprogress.Bar) string {
	return fmt.Sprintf(fmt.Sprintf("Step %%d/%%d: %%-%ds", reporter.stepsMaxLen), reporter.currentStep, len(reporter.steps), reporter.steps[reporter.currentStep-1])
}

func (reporter *progressReporter) appendFunc(b *uiprogress.Bar) string {
	return fmt.Sprintf("%0.1f", float64(reporter.bar.Current())/float64(reporter.bar.Total))
}

func (reporter *progressReporter) SetStep(step int) {
	reporter.currentStep = step
	reporter.SetProgress(0) // Reset progress
}

func (reporter *progressReporter) SetSteps(steps []string) {
	reporter.steps = steps
	reporter.stepsMaxLen = 0
	if reporter.currentStep > len(steps) {
		reporter.currentStep = len(steps) - 1
	}

	if reporter.currentStep < 0 {
		reporter.currentStep = 0
	}

	// Find the max len for the steps
	for _, step := range steps {
		if len(step) > reporter.stepsMaxLen {
			reporter.stepsMaxLen = len(step)
		}
	}
}

func (reporter *progressReporter) SetProgress(current int) {
	reporter.bar.Set(current)
}

func (reporter *progressReporter) SetTotal(total int) {
	reporter.bar.Total = total
}
