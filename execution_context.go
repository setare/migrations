package migrations

type ProgressReporter interface {
	SetStep(current int)
	SetSteps(steps []string)
	SetTotal(total int)
	SetProgress(progress int)
}

type ExecutionContext interface {
	EnableProgress() ProgressReporter
	Data() interface{}
}

type executionContext struct {
	data interface{}
}

func NewExecutionContext(data interface{}) ExecutionContext {
	return &executionContext{
		data: data,
	}
}

func (ctx *executionContext) EnableProgress() ProgressReporter {
	panic("not implemented") // TODO: Implement
}

func (ctx *executionContext) Data() interface{} {
	return ctx.data
}
