package main

import "errors"

type ColumnType string

const (
	Question    ColumnType = "question"
	Demographic            = "demographic"
)

var ErrMinValueGreaterThanMaxValue = errors.New("min value cannot be greater than max value")
var ErrEmptyColumnCode = errors.New("column must have non-empty code")
var ErrEmptyColumnText = errors.New("column must have non-empty text")

// Column represents a column of a survey dataset.
type Column struct {
	Code     string `yaml:"code" json:"code"`
	Text     string `yaml:"text" json:"text"`
	MinValue int    `yaml:"min_value" json:"min_value"`
	MaxValue int    `yaml:"max_value" json:"max_value"`
	Nullable bool   `yaml:"nullable" json:"nullable"`
	OfType   string `yaml:"of_type" json:"of_type"`
}

// NewColumn returns a new column and validates input data.
func NewColumn(code, text, ofType string, minValue, maxValue int, nullable bool) (Column, error) {
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
		Code:     code,
		Text:     text,
		MinValue: minValue,
		MaxValue: maxValue,
		Nullable: nullable,
		OfType:   ofType,
	}, nil
}
