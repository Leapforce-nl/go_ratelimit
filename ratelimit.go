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
	headerRemaining string
	headerReset     string
	endpoints       map[string]Status
}

type Status struct {
	limit     int
	remaining int
	reset     int
}

type ServiceConfig struct {
	HeaderRemaining *string
	HeaderReset     *string
}

func NewService(serviceConfig *ServiceConfig) *Service {
	headerRemaining := headerRemainingDefault

	if serviceConfig != nil {
		if serviceConfig.HeaderRemaining != nil {
			headerRemaining = *serviceConfig.HeaderRemaining
		}
	}

	headerReset := headerResetDefault

	if serviceConfig != nil {
		if serviceConfig.HeaderReset != nil {
			headerReset = *serviceConfig.HeaderReset
		}
	}
	return &Service{
		headerRemaining: headerRemaining,
		headerReset:     headerReset,
		endpoints:       make(map[string]Status),
	}
}

func (service *Service) InitEndpoint(endpoint string, limit int, remaining int, reset int) {
	if endpoint == "" {
		return
	}
	service.endpoints[endpoint] = Status{limit, remaining, reset}
}

func (service *Service) Set(endpoint string, response *http.Response) error {
	remaining, err := strconv.Atoi(response.Header.Get(service.headerRemaining))
	if err != nil {
		return err
	}
	reset, err := strconv.Atoi(response.Header.Get(service.headerReset))
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
