package util

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const BACKGROUND_HEALTHCHECK_TIMEOUT = 10 * time.Second

func RunWithBackGroundLoop(timeout time.Duration, worker func(context.Context) error) {
	ctx := context.Background()
	for {
		fmt.Println("---------- background run ----------")
		err := worker(ctx)
		if err != nil {
			LogErr(err)
		}
		time.Sleep(timeout)
		if err == nil {
			return
		}
		fmt.Println("---------- background run finished ----------")
	}
}

func RunWithBackGround(worker func(context.Context) error) {
	ctx := context.Background()
	LogErr(worker(ctx))
}

func UnMarshalStruct(in interface{}, out interface{}) error {
	bytes, err := json.Marshal(in)
	if err != nil {
		return AddStack(err)
	}
	return AddStack(json.Unmarshal(bytes, out))
}
