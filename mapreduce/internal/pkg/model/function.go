package model

import (
	"github.com/pkg/errors"
	"mapreduce/internal/pkg/utils"
)

type MapFunction func(string, string) []KeyValue
type ReduceFunction func(string, []string) string

type Function struct {
	Name   string
	Map    MapFunction
	Reduce ReduceFunction
}

func (f *Function) ExecuteMap(filenames []string) (chan []KeyValue, error) {
	lines, err := utils.ReadFiles(filenames)
	if err != nil {
		return nil, errors.Wrap(err, "execute map function error")
	}

	ans := make(chan []KeyValue)

	go func() {
		defer close(ans)

		results := make([]KeyValue, 0, len(lines))
		for _, line := range lines {
			results = append(results, f.Map(filenames[0], string(line))...)
		}

		ans <- results
	}()

	return ans, nil
}

func (f *Function) ExecuteReduce(key string, values []string) string {
	return f.Reduce(key, values)
}
