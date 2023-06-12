package map_worker

type WorkerStatus int

const (
	InitStatus WorkerStatus = iota
	RunningStatus
	IdleStatus
	OfflineStatus
	FailedStatus
)
