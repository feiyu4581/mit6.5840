package worker

type Worker interface {
	Run(filename string, params ...string)
	IsDone() bool
	IsFailed() bool
	GetPercent() float64
}
