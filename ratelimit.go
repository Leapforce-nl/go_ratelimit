package go_ratelimit

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	headerRemainingDefault string = "x-rate-limit-remaining"
	headerResetDefault     string = "x-rate-limit-reset"
)

type RateLimit struct {
	headerRemaining *string
	headerReset     *string
	endpoints       map[string]Status
}

type Status struct {
	limit     int
	remaining int
	reset     int
}

func NewRateLimit() *RateLimit {
	rl := RateLimit{}
	rl.endpoints = make(map[string]Status)

	return &rl
}

func (rm RateLimit) getHeaderRemaining() string {
	if rm.headerRemaining == nil {
		return headerRemainingDefault
	}
	if *rm.headerRemaining == "" {
		return headerRemainingDefault
	}

	return *rm.headerRemaining
}

func (rm RateLimit) getHeaderReset() string {
	if rm.headerReset == nil {
		return headerResetDefault
	}
	if *rm.headerReset == "" {
		return headerResetDefault
	}

	return *rm.headerReset
}

func (rl *RateLimit) InitEndpoint(endpoint string, limit int, remaining int, reset int) {
	if endpoint == "" {
		return
	}
	rl.endpoints[endpoint] = Status{limit, remaining, reset}
}

func (rl *RateLimit) Set(endpoint string, response *http.Response) error {
	remaining, err := strconv.Atoi(response.Header.Get(rl.getHeaderRemaining()))
	if err != nil {
		return err
	}
	reset, err := strconv.Atoi(response.Header.Get(rl.getHeaderReset()))
	if err != nil {
		return err
	}

	status, ok := rl.endpoints[endpoint]
	if !ok {
		status = Status{}
	}
	status.remaining = remaining
	status.reset = reset

	//fmt.Println(endpoint, remaining)

	rl.endpoints[endpoint] = status

	return nil
}

func (rl *RateLimit) Check(endpoint string) {

	status, ok := rl.endpoints[endpoint]
	if !ok {
		rl.endpoints[endpoint] = Status{}
		return
	}

	if status.remaining < 1 {
		reset := time.Unix(int64(status.reset), 0)
		ms := reset.Sub(time.Now()).Milliseconds()

		if ms > 0 {
			fmt.Println("waiting ms:", ms)
			time.Sleep(time.Duration(ms+1000) * time.Millisecond)
		}
	}

	return
}
