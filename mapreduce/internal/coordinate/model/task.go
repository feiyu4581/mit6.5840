package model

type TaskStatus int

const (
	InitStatus TaskStatus = iota
	RunningMapStatus
	RunningReduceStatus
	FinishedStatus
	FailedStatus
)

type TaskRequest struct {
	FileNames []string `json:"file_names" binding:"required"`
	Name      string   `json:"name" binding:"required"`
	Worker    string   `json:"worker" binding:"required"`
	MNums     int      `json:"m_nums" binding:"required"`
	RNums     int      `json:"r_nums" binding:"required"`
}

type Task struct {
	FileNames []string   `json:"file_names"`
	Name      string     `json:"name"`
	Worker    string     `json:"worker"`
	MNums     int        `json:"m_nums"`
	RNums     int        `json:"r_nums"`
	Status    TaskStatus `json:"status"`
	Workers   []*Worker  `json:"workers"`
}

func NewTask(fileNames []string, name string, worker string, mNums int, rNums int) *Task {
	return &Task{
		Name:      name,
		FileNames: fileNames,
		MNums:     mNums,
		RNums:     rNums,
		Worker:    worker,
		Status:    InitStatus,
	}
}
