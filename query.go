package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/dustin/go-jsonpointer"
)

func processDoc(collection [][]*string, doc string, ptrs []string) {
	j := map[string]interface{}{}
	err := json.Unmarshal([]byte(doc), &j)
	if err != nil {
		for i := range ptrs {
			collection[i] = append(collection[i], nil)
		}
		return
	}
	for i, p := range ptrs {
		val := jsonpointer.Get(j, p)
		switch x := val.(type) {
		case string:
			collection[i] = append(collection[i], &x)
		case int, uint, int64, float64, uint64, bool:
			v := fmt.Sprintf("%v", val)
			collection[i] = append(collection[i], &v)
		default:
			log.Printf("Ignoring %T", val)
			collection[i] = append(collection[i], nil)
		}
	}
}

type Reducer func(input []*string) interface{}

func reduce(collection [][]*string, reducers []Reducer) []interface{} {
	rv := make([]interface{}, len(collection))
	for i, a := range collection {
		rv[i] = reducers[i](a)
	}
	return rv
}

func convertToint64(in []*string) []int64 {
	rv := make([]int64, 0, len(in))
	for _, v := range in {
		if v != nil {
			x, err := strconv.ParseInt(*v, 10, 64)
			if err == nil {
				rv = append(rv, x)
			}
		}
	}
	return rv
}

var reducers = map[string]Reducer{
	"any": func(input []*string) interface{} {
		for _, v := range input {
			if v != nil {
				return *v
			}
		}
		return nil
	},
	"count": func(input []*string) interface{} {
		rv := 0
		for _, v := range input {
			if v != nil {
				rv++
			}
		}
		return rv
	},
	"sum": func(input []*string) interface{} {
		rv := int64(0)
		for _, v := range convertToint64(input) {
			rv += v
		}
		return rv
	},
	"max": func(input []*string) interface{} {
		rv := int64(math.MinInt64)
		for _, v := range convertToint64(input) {
			if v > rv {
				rv = v
			}
		}
		return rv
	},
	"min": func(input []*string) interface{} {
		rv := int64(math.MaxInt64)
		for _, v := range convertToint64(input) {
			if v < rv {
				rv = v
			}
		}
		return rv
	},
	"avg": func(input []*string) interface{} {
		nums := convertToint64(input)
		sum := int64(0)
		for _, v := range nums {
			sum += v
		}
		return float64(sum) / float64(len(nums))
	},
}