package go_ratelimit

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	headerRemainingDefault string = "x-rate-limit-remaining"
	headerResetDefault     string = "x-rate-limit-reset"
)

type RateLimit struct {
	headerRemaining *string
	headerReset     *string
	groups          map[string]Status
}

type Status struct {
	remaining int64
	reset     int64
}

func NewRateLimit() *RateLimit {
	rl := RateLimit{}
	rl.groups = make(map[string]Status)

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

func (rl *RateLimit) Set(group string, response *http.Response) error {
	remaining, err := strconv.ParseInt(response.Header.Get(rl.getHeaderRemaining()), 10, 64)
	if err != nil {
		return err
	}
	reset, err := strconv.ParseInt(response.Header.Get(rl.getHeaderReset()), 10, 64)
	if err != nil {
		return err
	}

	status, ok := rl.groups[group]
	if !ok {
		status = Status{}
	}
	status.remaining = remaining
	status.reset = reset

	rl.groups[group] = status

	return nil
}

func (rl *RateLimit) Check(group string) {

	status, ok := rl.groups[group]
	if !ok {
		rl.groups[group] = Status{}
		return
	}

	fmt.Println("remaining", status.remaining)
	fmt.Println("reset", status.reset)

	return

}
