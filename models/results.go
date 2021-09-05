package models

import (
	"errors"
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

var ErrWrongItemType = errors.New("last can't be called on text items")

func (set ResultSet) Min() (float64, error) {
	if set.Type() == Text {
		logger.Error(resultSetLoggingArea, "Something tried to calculate min for item with wrong type!")
		return 0, ErrWrongItemType
	}
	var min float64

	for i, k := range set.Results {
		if i == 0 {
			min = k.ValueNumeric
		}

		if k.ValueNumeric < min {
			min = k.ValueNumeric
		}
	}
	return min, nil
}

func (set ResultSet) Max() (float64, error) {
	if set.Type() == Text {
		logger.Error(resultSetLoggingArea, "Something tried to calculate max for item with wrong type!")
		return 0, ErrWrongItemType
	}
	var max float64

	for i, k := range set.Results {
		if i == 0 {
			max = k.ValueNumeric
		}

		if k.ValueNumeric > max {
			max = k.ValueNumeric
		}
	}
	return max, nil
}

func (set ResultSet) Avg() (float64, error) {
	if set.Type() == Text {
		logger.Error(resultSetLoggingArea, "Something tried to calculate avg for item with wrong type!")
		return 0, ErrWrongItemType
	}

	var sum float64

	for _, k := range set.Results {
		sum += k.ValueNumeric
	}
	return sum / float64(len(set.Results)), nil
}

func (set ResultSet) NumericLast(Count int) (ResultSet, error) {
	var dummyResultSet ResultSet
	dummyResultSet.Results = make([]Result, 0)

	if set.Type() == Text {
		logger.Error(resultSetLoggingArea, "Something tried to get last numeric values for item with wrong type!")
		return dummyResultSet, ErrWrongItemType
	}

	if Count < 1 {
		logger.Error(resultSetLoggingArea, "Something tried to get the last", Count, "results of a resultset, which isn't a valid query")
		return dummyResultSet, errors.New("count of last items can't be less than 0")
	}

	if len(set.Results) == 0 {
		return dummyResultSet, nil
	}

	indexPointer := len(set.Results)
	indexTarget := indexPointer - Count

	//If we request more items than present in the result set, just return everything we have
	if indexTarget < 0 {
		indexTarget = 0
	}

	slicedResults := ResultSet{
		Results: set.Results[indexTarget:indexPointer],
	}

	return slicedResults, nil
}
