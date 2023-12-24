package main

import (
	"errors"
	"strconv"
	"strings"
)

type ColumnType string

const (
	Question    ColumnType = "question"
	Demographic            = "demographic"
	Digits      string     = "0123456789"
)

var ErrMinValueGreaterThanMaxValue = errors.New("min value cannot be greater than max value")
var ErrEmptyColumnCode = errors.New("column must have non-empty code")
var ErrEmptyColumnText = errors.New("column must have non-empty text")
var ErrEmptyOrgNodeColumn = errors.New("org node column must be non-empty")
var ErrSchemaContainsColumnsWithDuplicatedCodes = errors.New("schema cannot have columns with duplicated codes")
var ErrInvalidOrgNodeString = errors.New("invalid org node string")

// Column represents a column of a survey dataset.
type Column struct {
	Code       string     `yaml:"code" json:"code"`
	Text       string     `yaml:"text" json:"text"`
	MinValue   int        `yaml:"min_value" json:"min_value"`
	MaxValue   int        `yaml:"max_value" json:"max_value"`
	Nullable   bool       `yaml:"nullable" json:"nullable"`
	ColumnType ColumnType `yaml:"column_type" json:"column_type"`
}

// Schema represents a schema of survey (metadata).
type Schema struct {
	OrgNodeCol string   `yaml:"org_node_col" json:"org_node_col"`
	Columns    []Column `yaml:"columns" json:"columns"`
}

func validateColumnsCodeUniqueness(columns []Column) error {
	codes := make(map[string]struct{})
	for _, c := range columns {
		_, exists := codes[c.Code]
		if exists {
			return ErrSchemaContainsColumnsWithDuplicatedCodes
		}
		codes[c.Code] = struct{}{}
	}
	return nil
}

func validateSchema(schema Schema) error {
	if schema.OrgNodeCol == "" {
		return ErrEmptyOrgNodeColumn
	}

	err := validateColumnsCodeUniqueness(schema.Columns)
	if err != nil {
		return err
	}

	return nil
}

// NewColumn returns a new column and validates input data.
func NewColumn(code, text string, minValue, maxValue int, nullable bool, columnType ColumnType) (Column, error) {
	var col Column

	if minValue > maxValue {
		return col, ErrMinValueGreaterThanMaxValue
	}

	if code == "" {
		return col, ErrEmptyColumnCode
	}

	if text == "" {
		return col, ErrEmptyColumnText
	}

	return Column{
		Code:       code,
		Text:       text,
		MinValue:   minValue,
		MaxValue:   maxValue,
		Nullable:   nullable,
		ColumnType: columnType,
	}, nil
}

type OrgNode struct {
	levels []int
}

func OrgNodeFromString(s string, sep string) (OrgNode, error) {
	var orgNode OrgNode

	allowedChars := Digits + sep
	chars := make([]string, 0)
	for _, c := range s {
		if strings.Contains(allowedChars, string(c)) {
			chars = append(chars, string(c))
		}
	}

	levels := make([]int, 0)
	for _, p := range strings.Split(strings.Join(chars, ""), sep) {
		if p == "" {
			continue
		}
		level, err := strconv.Atoi(p)
		if err != nil {
			return orgNode, ErrInvalidOrgNodeString
		}
		levels = append(levels, level)
	}
	return OrgNode{levels: levels}, nil
}

type OrgNodes []OrgNode
