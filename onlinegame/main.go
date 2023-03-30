package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type Clan struct {
	NumberOfPlayers int     `json:"numberOfPlayers`
	Points          int     `json:"points"`
	PointsPerPlayer float64 `json:"-"`
}

func (clan Clan) String() string {
	return fmt.Sprintf("Clan{%d %d %.4f}", clan.NumberOfPlayers, clan.Points, clan.PointsPerPlayer)
}

type Clans []*Clan

func (clans Clans) Len() int {
	return len(clans)
}

func (clans Clans) Less(i, j int) bool {
	return clans[i].PointsPerPlayer < clans[j].PointsPerPlayer
}

func (clans Clans) Swap(i, j int) {
	clans[i], clans[j] = clans[j], clans[i]
}

type Request struct {
	GroupSize int   `json:"groupCount"`
	Clans     Clans `json:"clans"`
}

func parseRequest(data []byte) (Clans, int, error) {
	var request Request
	err := json.Unmarshal(data, &request)
	if err != nil {
		return nil, 0, errors.Wrap(err, "unable to parse request")
	}

	for _, clan := range request.Clans {
		clan.PointsPerPlayer = float64(clan.Points) / float64(clan.NumberOfPlayers)
	}

	sort.Sort(sort.Reverse(request.Clans))

	return request.Clans, request.GroupSize, nil
}

func calculateGroups(clans Clans, groupSize int) []Clans {
	groups := make([]Clans, 1)

	currentSize := 0
	currentGroup := 0
	for _, clan := range clans {
		if currentSize+clan.NumberOfPlayers > groupSize { // TODO(misiek08): could use bin packing instead of such naive approach - PRs welcome :)
			groups = append(groups, make([]*Clan, 0))
			currentGroup += 1
			currentSize = 0
		}
		currentSize += clan.NumberOfPlayers
		groups[currentGroup] = append(groups[currentGroup], clan)
	}

	return groups
}

func main() {
	fasthttp.ListenAndServe(":8080", func(ctx *fasthttp.RequestCtx) {
		clans, groupSize, err := parseRequest(ctx.PostBody())
		if err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Response.SetBody([]byte(fmt.Sprintf("unable to parse request: %v", err)))
		}
		groups := calculateGroups(clans, groupSize)
		groupsJSON, err := json.Marshal(groups)
		if err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Response.SetBody([]byte(fmt.Sprintf("unable to serialize groups: %v", err)))
		} else {
			ctx.Response.SetBody(groupsJSON)
		}
	})
}
