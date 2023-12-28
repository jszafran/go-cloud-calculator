package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type SurveySchemaValidationMode string

const (
	FailFast         SurveySchemaValidationMode = "fail_fast"
	CaptureAllErrors                            = "capture_all_errors"
)

var ErrInvalidSchemaValidationMode = errors.New("invalid schema validation mode")
var ErrFloatValueWithFractionalPartDetected = errors.New("float error with fractional value is not allowed")
var ErrSchemaMinMaxValuesViolation = errors.New("value does not fit into schema min or max value")

// Dataset represents a survey data along with metadata (schema).
type Dataset struct {
	Schema   Schema
	OrgNodes []OrgNode
	Data     map[string]int
}

type DataError struct {
	lineNum    int
	columnName string
	message    string
}

func (de DataError) Error() string {
	return fmt.Sprintf("line %d | column %s | %s", de.lineNum, de.columnName, de.message)
}

type DatasetLoadAttempt struct {
	succeeded  bool
	dataset    Dataset
	error      error
	dataErrors []DataError
}

// stringToInt converts a numeric value from string into int (i.e.
func stringDataPointToInt(s string) (int, error) {
	var result int
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return result, err
	}

	result = int(f)
	if float64(result) != f {
		return result, ErrFloatValueWithFractionalPartDetected
	}

	return result, nil
}

type ParsedHeader struct {
	IndexToCode    map[int]string
	CodeToIndex    map[string]int
	OrgColFound    bool
	OrgColPos      int
	MissingColumns []string
	ExtraColumns   []string
}

func inSliceOfStrings(v string, slice []string) bool {
	for _, s := range slice {
		if s == v {
			return true
		}
	}
	return false
}

// ParseHeader parses given header and compares it against provided schema.
func ParseHeader(header []string, schema Schema) ParsedHeader {
	var parsedHeader ParsedHeader

	indexToCode := make(map[int]string)
	codeToIndex := make(map[string]int)
	missingColumns := make([]string, 0)
	extraColumns := make([]string, 0)

	for i, val := range header {
		indexToCode[i] = val
		codeToIndex[val] = i
	}

	parsedHeader.CodeToIndex = codeToIndex
	parsedHeader.IndexToCode = indexToCode

	orgNodePos, exists := codeToIndex[schema.OrgNodeCol]
	if !exists {
		parsedHeader.OrgColFound = false
		parsedHeader.OrgColPos = -1
	} else {
		parsedHeader.OrgColFound = true
		parsedHeader.OrgColPos = orgNodePos
	}

	for _, cn := range schema.ColumnsNames() {
		_, exists = codeToIndex[cn]
		if !exists {
			missingColumns = append(missingColumns, cn)
		}
	}

	for k, _ := range codeToIndex {
		if !inSliceOfStrings(k, schema.AllFieldsNames()) {
			extraColumns = append(extraColumns, k)
		}
	}

	parsedHeader.MissingColumns = missingColumns
	parsedHeader.ExtraColumns = extraColumns

	return parsedHeader
}

func NewDatasetFromCSV(csvPath string, schema Schema, validationMode SurveySchemaValidationMode) DatasetLoadAttempt {
	var (
		header             []string
		dataset            Dataset
		datasetLoadAttempt DatasetLoadAttempt
	)
	datasetLoadAttempt.succeeded = false

	if validationMode != FailFast && validationMode != CaptureAllErrors {
		datasetLoadAttempt.error = ErrInvalidSchemaValidationMode
		return datasetLoadAttempt
	}

	csvFile, err := os.Open(csvPath)
	if err != nil {
		datasetLoadAttempt.error = err
		return datasetLoadAttempt
	}
	r := csv.NewReader(csvFile)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			datasetLoadAttempt.error = err
			return datasetLoadAttempt
		}

		if header == nil {
			header = record
			parsedHeader := ParseHeader(header, schema)
			fmt.Printf("Parsed header data: %+v\n", parsedHeader)
			continue
		} else {
			for _, dp := range record {
				dataPoint, err := stringDataPointToInt(dp)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Converting %s to %d\n", dp, dataPoint)
			}

		}
	}

	datasetLoadAttempt.succeeded = true
	datasetLoadAttempt.dataset = dataset

	return datasetLoadAttempt
}
