package model

type WorkerMode int

const (
	MapMode WorkerMode = iota
	ReduceMode
)
