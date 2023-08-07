package model

type Task struct {
	FileNames []string `json:"file_names" binding:"required"`
	Name      string   `json:"name" binding:"required"`
	Worker    string   `json:"worker" binding:"required"`
	MNums     int      `json:"m_nums" binding:"required"`
	RNums     int      `json:"r_nums" binding:"required"`
}
