package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
)

type AtmRequests []AtmRequest

type AtmResponse struct {
	Region int `json:"region"`
	AtmID  int `json:"atmId"`
}

type AtmRequest struct {
	Region      int    `json:"region"`
	RequestType string `json:"requestType"`
	AtmID       int    `json:"atmId"`
	Value       int    `json:"-"`
}

func (reqs AtmRequests) Len() int {
	return len(reqs)
}

func (reqs AtmRequests) Less(i, j int) bool {
	if reqs[i].Region != reqs[j].Region {
		return reqs[i].Region < reqs[j].Region
	}
	return reqs[i].Value > reqs[j].Value
}

func (reqs AtmRequests) Swap(i, j int) {
	reqs[i], reqs[j] = reqs[j], reqs[i]
}

func parseRequests(data []byte) (AtmRequests, error) {
	var requests AtmRequests
	err := json.Unmarshal(data, &requests)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse transactions")
	}
	return requests, nil
}

func calculateValue(requestType string) int {
	switch requestType {
	case "STANDARD":
		return 10
	case "SIGNAL_LOW":
		return 20
	case "PRIORITY":
		return 30
	case "FAILURE_RESTART":
		return 40
	}
	return 0
}

func sortRequests(reqs AtmRequests) []AtmResponse {
	responseList := make([]AtmResponse, 0, len(reqs))

	seen := make(map[AtmResponse]bool)

	for i, req := range reqs {
		reqs[i].Value = calculateValue(req.RequestType)
	}

	sort.Sort(reqs)

	for _, req := range reqs {
		res := AtmResponse{Region: req.Region, AtmID: req.AtmID}
		if _, ok := seen[res]; ok {
			continue
		}
		seen[res] = true
		responseList = append(responseList, res)
	}

	return responseList
}

func main() {
	decimal.MarshalJSONWithoutQuotes = true

	fasthttp.ListenAndServe(":8080", func(ctx *fasthttp.RequestCtx) {
		requests, err := parseRequests(ctx.PostBody())
		if err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Response.SetBody([]byte(fmt.Sprintf("unable to parse request: %v", err)))
		}
		atmList := sortRequests(requests)
		atmListJSON, err := json.Marshal(atmList)
		if err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Response.SetBody([]byte(fmt.Sprintf("unable to serialize accounts: %v", err)))
		} else {
			ctx.Response.SetBody(atmListJSON)
		}
	})
}
