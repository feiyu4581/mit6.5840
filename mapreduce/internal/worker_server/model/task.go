package model

import (
	"fmt"
	"github.com/pkg/errors"
	"mapreduce/internal/pkg/log"
	"mapreduce/internal/pkg/model"
	"mapreduce/internal/pkg/utils"
	"os"
	"strings"
	"time"
)

type Status int

const (
	InitStatus Status = iota
	RunningStatus
	DoneStatus
	FailedStatus
)

type Task struct {
	Name      string
	TaskId    int64
	TaskIndex int64
	Filename  string
	Filenames []string
	Nums      int64
	StartTime time.Time
	EndTime   time.Time
	Function  *model.Function
	Status    Status
	Err       error

	DoneC        chan struct{}
	ResFileNames []string
}

func NewTask(taskId, taskIndex int64, name string, filenames []string, nums int64) *Task {
	return &Task{
		TaskId:    taskId,
		TaskIndex: taskIndex,
		Name:      name,
		Filenames: filenames,
		Nums:      nums,
		Status:    InitStatus,
		DoneC:     make(chan struct{}),
	}
}

func (task *Task) SplitMapResults(res []model.KeyValue) ([]string, error) {
	log.Info("[%s] compute map results, length=%d", task.Name, len(res))
	storeValues := make(map[string][]string)
	for i := range res {
		storeValues[res[i].Key] = append(storeValues[res[i].Key], res[i].Value)
	}

	lines := make([]utils.Line, 0, len(storeValues))
	for key, values := range storeValues {
		builder := strings.Builder{}
		builder.WriteString(key)
		builder.WriteString(",")
		builder.WriteString(strings.Join(values, ","))

		lines = append(lines, utils.Line(builder.String()))
	}

	fileNames, err := utils.SplitLinesToFile(lines, int(task.Nums), fmt.Sprintf("map-%s-%d-%d", task.Name, task.TaskId, time.Now().UnixMilli()))
	if err != nil {
		return nil, errors.Wrap(err, "split map results to file error")
	}
	log.Info("[%s] split lines to file success, filenames=%s", task.Name, strings.Join(fileNames, ","))
	resFileName := fmt.Sprintf("/tmp/%s.%d.%d", task.Name, task.TaskId, task.TaskIndex)

	resFileNames := make([]string, 0, len(fileNames))
	for index, fileName := range fileNames {
		resFileNames = append(resFileNames, fmt.Sprintf("%s.%d.map", resFileName, index))
		if err = os.Rename(fileName, resFileNames[len(resFileNames)-1]); err != nil {
			return nil, errors.Wrap(err, "rename file error")
		}
	}

	log.Info("[%s] rename file success, filenames=%s", task.Name, strings.Join(resFileNames, ","))
	return resFileNames, nil
}

func (task *Task) SplitReduceResults(values []model.KeyValue) (string, error) {
	lines := make([]utils.Line, 0, len(values))
	for _, value := range values {
		lines = append(lines, utils.Line(fmt.Sprintf("%s,%s", value.Key, value.Value)))
	}

	fileNames, err := utils.SplitLinesToFile(lines, 1, fmt.Sprintf("reduce-%s-%d-%d", task.Name, task.TaskId, time.Now().UnixMilli()))
	if err != nil {
		return "", errors.Wrap(err, "split reduce results to file error")
	}

	if len(fileNames) != 1 {
		return "", errors.New(fmt.Sprintf("split reduce results to file error: %s", strings.Join(fileNames, ",")))
	}

	filename := fmt.Sprintf("/tmp/%s.%d.%d.reduce", task.Name, task.TaskId, task.TaskIndex)
	if err = os.Rename(fileNames[0], filename); err != nil {
		return "", errors.Wrap(err, "rename file error")
	}

	return filename, nil
}

func (task *Task) ParseReduceLines(lines []utils.Line) (map[string][]string, error) {
	res := make(map[string][]string, len(lines))
	for _, line := range lines {
		kv := strings.Split(string(line), ",")
		if len(kv) < 2 {
			return nil, errors.New(fmt.Sprintf("line[%s] is invalid", line))
		}

		res[kv[0]] = append(res[kv[0]], kv[1:]...)
	}

	return res, nil
}

func (task *Task) Start(mode model.WorkerMode) (err error) {
	defer func() {
		if err != nil {
			task.SetError(err)
		}
	}()

	task.StartTime = time.Now()
	task.Status = RunningStatus
	switch mode {
	case model.MapMode:
		resC, mapErr := task.Function.ExecuteMap(task.Filenames)
		if mapErr != nil {
			return errors.Wrap(err, "execute map function error")
		}
		go func() {
			res := <-resC
			if filenames, err := task.SplitMapResults(res); err != nil {
				task.SetError(err)
			} else {
				task.SetDone(filenames)
			}
		}()
		return nil
	case model.ReduceMode:
		lines, err := utils.ReadFiles(task.Filenames)
		if err != nil {
			return errors.Wrap(err, "execute reduce function error")
		}

		go func() {
			reduceValues, err := task.ParseReduceLines(lines)
			if err != nil {
				task.SetError(errors.Wrap(err, "parse reduce lines error"))
				return
			}

			res := make([]model.KeyValue, 0, len(reduceValues))
			for key, values := range reduceValues {
				res = append(res, model.KeyValue{
					Key:   key,
					Value: task.Function.ExecuteReduce(key, values),
				})
			}

			if filename, err := task.SplitReduceResults(res); err != nil {
				task.SetError(err)
			} else {
				task.SetDone([]string{filename})
			}
		}()

		return nil
	}

	return errors.New(fmt.Sprintf("mode[%d] is not supported", mode))
}

func (task *Task) SetError(err error) {
	defer close(task.DoneC)

	task.Err = err
	task.EndTime = time.Now()
	task.Status = FailedStatus
}

func (task *Task) SetDone(filenames []string) {
	defer close(task.DoneC)

	task.EndTime = time.Now()
	task.Status = DoneStatus
	task.ResFileNames = filenames
}

func (task *Task) SetFunction(function *model.Function) {
	task.Function = function
}

func (task *Task) GetStatus() Status {
	return task.Status
}

func (task *Task) IsFinish() bool {
	return task.Status == DoneStatus || task.Status == FailedStatus
}
