package models

import (
	"errors"
	"math"
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Result stores a single result from a item
type Result struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ItemID       primitive.ObjectID
	HostID       primitive.ObjectID
	Type         ReturnType
	CapturedAt   time.Time
	ValueString  string
	ValueNumeric float64
	Error        string
}

//ResultSet stores a collection of results
type ResultSet struct {
	Results []Result
}

//Type returns the ReturnType used by the items in the collection
//All items in a resultset have the same resulttype -> Return the returntype of the first item
//Fallback to Numeric if there are no results as to not brick arithemetic calculations
func (set ResultSet) Type() ReturnType {
	if len(set.Results) == 0 {
		return Numeric
	}

	return set.Results[0].Type
}

const resultSetLoggingArea = "EVAL"

//ErrWrongItemType is returned if numeric functions are used on text items or vice versa
var ErrWrongItemType = errors.New("function can't be called on text items")

//ErrNoResults is returned if a ResultSet has no / not enough values
var ErrNoResults = errors.New("there are no results to return")

//Min returns the minimum value within the ResultSet
//If Limit != 0 only the last N values are evaluated
//Limit is a float64 because of govaluate, which always seems to pass arguments as float64 (even if no decimal point is present)
func (set ResultSet) Min(Limit float64) (float64, error) {
	if set.Type() == Text {
		logger.Error(resultSetLoggingArea, "Something tried to calculate min for item with wrong type!")
		return 0, ErrWrongItemType
	}

	integerLimit := int(Limit)

	if len(set.Results) < integerLimit || integerLimit == 0 {
		integerLimit = len(set.Results)
	}

	var min float64

	for i, k := range set.Results[:integerLimit] {
		if i == 0 {
			min = k.ValueNumeric
		}

		if k.ValueNumeric < min {
			min = k.ValueNumeric
		}
	}
	return min, nil
}

//Max returns the maximum value within the ResultSet
//If Limit != 0 only the last N values are evaluated
//Limit is a float64 because of govaluate, which always seems to pass arguments as float64 (even if no decimal point is present)
func (set ResultSet) Max(Limit float64) (float64, error) {
	if set.Type() == Text {
		logger.Error(resultSetLoggingArea, "Something tried to calculate max for item with wrong type!")
		return 0, ErrWrongItemType
	}

	integerLimit := int(Limit)

	if len(set.Results) < integerLimit || integerLimit == 0 {
		integerLimit = len(set.Results)
	}

	var max float64

	for i, k := range set.Results[:integerLimit] {
		if i == 0 {
			max = k.ValueNumeric
		}

		if k.ValueNumeric > max {
			max = k.ValueNumeric
		}
	}
	return max, nil
}

//Avg returns the average value within the ResultSet
//If Limit != 0 only the last N values are evaluated
//Limit is a float64 because of govaluate, which always seems to pass arguments as float64 (even if no decimal point is present)
func (set ResultSet) Avg(Limit float64) (float64, error) {
	if set.Type() == Text {
		logger.Error(resultSetLoggingArea, "Something tried to calculate avg for item with wrong type!")
		return 0, ErrWrongItemType
	}

	integerLimit := int(Limit)

	if len(set.Results) < integerLimit || integerLimit == 0 {
		integerLimit = len(set.Results)
	}

	var sum float64

	for _, k := range set.Results[:integerLimit] {
		sum += k.ValueNumeric
	}
	return sum / float64(len(set.Results)), nil
}

//Diff returns the difference between the last two value in the ResultSet
func (set ResultSet) Diff() (float64, error) {
	if set.Type() == Text {
		logger.Error(resultSetLoggingArea, "Something tried to calculate avg for item with wrong type!")
		return 0, ErrWrongItemType
	}

	if len(set.Results) == 0 {
		return 0, ErrNoResults
	}

	secondItem := 1
	if len(set.Results) == 1 {
		secondItem = 0
	}

	return math.Abs(math.Abs(set.Results[0].ValueNumeric) - math.Abs(set.Results[secondItem].ValueNumeric)), nil
}

//LastNumeric returns the last numeric value in the given ResultSet
func (set ResultSet) LastNumeric() (float64, error) {
	if set.Type() == Text {
		logger.Error(resultSetLoggingArea, "Something tried to calculate avg for item with wrong type!")
		return 0, ErrWrongItemType
	}

	if len(set.Results) == 0 {
		return 0, ErrNoResults
	}

	return set.Results[0].ValueNumeric, nil
}
