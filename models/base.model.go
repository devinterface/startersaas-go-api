package models

import (
	"strings"

	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseModel struct
type BaseModel struct{}

type PaginatedResultSet struct {
	Records []map[string]interface{} `json:"records"`
	Meta    map[string]int64         `json:"meta"`
}

type FilteredQuery struct {
	page   int64    `query:"page"`
	limit  int64    `query:"limit"`
	sort   string   `query:"sort"`
	filter []string `query:"filter"`
}

func (filteredQuery FilteredQuery) GetLimit() (limit int64) {
	if filteredQuery.limit <= 0 {
		return 10
	} else {
		return filteredQuery.limit
	}
}

func (filteredQuery FilteredQuery) GetPage() (limit int64) {
	if filteredQuery.page <= 0 {
		return 1
	} else {
		return filteredQuery.page
	}
}

func (filteredQuery FilteredQuery) GetSort() (sortMap bson.M) {
	if filteredQuery.sort != "" {
		sortMap = bson.M{}
		if strings.HasPrefix(filteredQuery.sort, "-") {
			field := strings.Split(filteredQuery.sort, "-")[1]
			sortMap[field] = -1
		} else {
			sortMap[filteredQuery.sort] = 1
		}
	}
	return sortMap
}

var filterOperators = []string{"$eq", "$in", "$gt", "$gte", "$lt", "$lte", "$regex", "$ne", "$text"}

func (filteredQuery FilteredQuery) GetFilter() (queryParams bson.M) {
	if filteredQuery.filter != nil {
		queryParams = bson.M{}
		for _, filterParam := range filteredQuery.filter {
			splittedFilterParam := strings.Split(filterParam, ":")
			if len(splittedFilterParam) != 3 {
				continue
			}
			field := splittedFilterParam[0]
			operator := splittedFilterParam[1]
			if !funk.Contains(filterOperators, operator) {
				continue
			}
			value := splittedFilterParam[2]
			queryParams[field] = bson.M{operator: value}
		}
	}
	return queryParams
}

func (filteredQuery FilteredQuery) BuildPaginatedFindOptions() (findOptions *options.FindOptions) {
	findOptions = options.Find()
	findOptions.SetLimit(filteredQuery.GetLimit())
	skip := (filteredQuery.GetPage() - 1) * filteredQuery.GetLimit()
	findOptions.SetSkip(int64(skip))
	findOptions.SetBatchSize(100)
	if filteredQuery.sort != "" {
		findOptions.Sort = filteredQuery.GetSort()
	}
	return findOptions
}
