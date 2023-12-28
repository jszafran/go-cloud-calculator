package main

import (
	"errors"
	"strconv"
	"strings"
)

type ColumnType string

const (
	Question    ColumnType = "question"
	Demographic ColumnType = "demographic"
	Digits      string     = "0123456789"
)

var (
	ErrMinValueGreaterThanMaxValue              = errors.New("min value cannot be greater than max value")
	ErrEmptyColumnCode                          = errors.New("column must have non-empty code")
	ErrEmptyColumnText                          = errors.New("column must have non-empty text")
	ErrEmptyOrgNodeColumn                       = errors.New("org node column must be non-empty")
	ErrSchemaContainsColumnsWithDuplicatedCodes = errors.New("schema cannot have columns with duplicated codes")
	ErrInvalidOrgNodeString                     = errors.New("invalid org node string")
	ErrInvalidColumnType                        = errors.New("invalid column type")
	ErrColNotFoundInSchema                      = errors.New("column was not found in schema")
)

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

// ColumnByCode returns a schema's column with matching code.
func (s Schema) ColumnByCode(code string) (Column, error) {
	var column Column
	for _, c := range s.Columns {
		if c.Code == code {
			return c, nil
		}
	}
	return column, ErrColNotFoundInSchema
}

// Questions returns all columns of type "question".
func (s Schema) Questions() []Column {
	res := make([]Column, 0)
	for _, c := range s.Columns {
		if c.ColumnType == Question {
			res = append(res, c)
		}
	}
	return res
}

// Demographics returns all columns of type "demographics".
func (s Schema) Demographics() []Column {
	res := make([]Column, 0)
	for _, c := range s.Columns {
		if c.ColumnType == Demographic {
			res = append(res, c)
		}
	}
	return res
}

func (s Schema) ColumnsNames() []string {
	fields := make([]string, 0)

	for _, c := range s.Columns {
		fields = append(fields, c.Code)
	}

	return fields
}

func (s Schema) AllFieldsNames() []string {
	allFields := make([]string, 0)
	allFields = append(allFields, s.OrgNodeCol)
	allFields = append(allFields, s.ColumnsNames()...)
	return allFields
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

	if columnType != Question && columnType != Demographic {
		return col, ErrInvalidColumnType
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

// OrgNode represents a node (an organization unit or manager) from organizational structure.
type OrgNode struct {
	levels []int
}

// OrgNodeFromString converts a string into an OrgNode struct or returns error if input data is incorrect.
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
