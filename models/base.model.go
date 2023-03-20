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
	Page   int64    `query:"page"`
	Limit  int64    `query:"limit"`
	Sort   string   `query:"sort"`
	Filter []string `query:"filter"`
}

func (filteredQuery FilteredQuery) GetLimit() (limit int64) {
	if filteredQuery.Limit <= 0 {
		return 10
	} else {
		return filteredQuery.Limit
	}
}

func (filteredQuery FilteredQuery) GetPage() (limit int64) {
	if filteredQuery.Page <= 0 {
		return 1
	} else {
		return filteredQuery.Page
	}
}

func (filteredQuery FilteredQuery) GetSort() (sortMap bson.M) {
	sortMap = bson.M{}
	if filteredQuery.Sort != "" {
		if strings.HasPrefix(filteredQuery.Sort, "-") {
			field := strings.Split(filteredQuery.Sort, "-")[1]
			sortMap[field] = -1
		} else {
			sortMap[filteredQuery.Sort] = 1
		}
	} else {
		sortMap["created_at"] = -1
	}
	return sortMap
}

var filterOperators = []string{"$eq", "$gt", "$gte", "$lt", "$lte", "$regex", "$ne", "$text"}

func (filteredQuery FilteredQuery) GetFilter() (queryParams bson.M) {
	queryParams = bson.M{}
	if filteredQuery.Filter != nil {
		for _, filterParam := range filteredQuery.Filter {
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
	if filteredQuery.Sort != "" {
		findOptions.Sort = filteredQuery.GetSort()
	}
	return findOptions
}

func MaxTeamsPerPlan(plan string) (maxGroups int) {
	if plan == StarterPlan {
		return 5
	} else if plan == BasicPlan {
		return 10
	} else if plan == PremiumPlan {
		return 20
	} else {
		return 0
	}
}
