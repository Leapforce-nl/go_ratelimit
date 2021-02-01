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

type Service struct {
	headerRemaining *string
	headerReset     *string
	endpoints       map[string]Status
}

type Status struct {
	limit     int
	remaining int
	reset     int
}

func NewService() *Service {
	service := Service{}
	service.endpoints = make(map[string]Status)

	return &service
}

func (service Service) getHeaderRemaining() string {
	if service.headerRemaining == nil {
		return headerRemainingDefault
	}
	if *service.headerRemaining == "" {
		return headerRemainingDefault
	}

	return *service.headerRemaining
}

func (service Service) getHeaderReset() string {
	if service.headerReset == nil {
		return headerResetDefault
	}
	if *service.headerReset == "" {
		return headerResetDefault
	}

	return *service.headerReset
}

func (service *Service) InitEndpoint(endpoint string, limit int, remaining int, reset int) {
	if endpoint == "" {
		return
	}
	service.endpoints[endpoint] = Status{limit, remaining, reset}
}

func (service *Service) Set(endpoint string, response *http.Response) error {
	remaining, err := strconv.Atoi(response.Header.Get(service.getHeaderRemaining()))
	if err != nil {
		return err
	}
	reset, err := strconv.Atoi(response.Header.Get(service.getHeaderReset()))
	if err != nil {
		return err
	}

	status, ok := service.endpoints[endpoint]
	if !ok {
		status = Status{}
	}
	status.remaining = remaining
	status.reset = reset

	//fmt.Println(endpoint, remaining)

	service.endpoints[endpoint] = status

	return nil
}

func (service *Service) Check(endpoint string) {

	status, ok := service.endpoints[endpoint]
	if !ok {
		service.endpoints[endpoint] = Status{}
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
