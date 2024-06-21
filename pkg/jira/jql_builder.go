package jira

import (
	"fmt"
	"strings"
)

type JQLBuilder struct {
	conditions []string
}

func NewJQLBuilder() *JQLBuilder {
	return &JQLBuilder{}
}

func (b *JQLBuilder) Equals(field, value string, isFunction bool) *JQLBuilder {
	if isFunction {
		b.conditions = append(b.conditions, fmt.Sprintf("%s = %s", field, value))
		return b
	}

	b.conditions = append(b.conditions, fmt.Sprintf("%s = \"%s\"", field, value))
	return b
}

func (b *JQLBuilder) In(field string, values []string) *JQLBuilder {
	quotedValues := make([]string, len(values))

	for i, v := range values {
		quotedValues[i] = fmt.Sprintf("\"%s\"", v)
	}

	condition := fmt.Sprintf("%s IN (%s)", field, strings.Join(quotedValues, ", "))
	b.conditions = append(b.conditions, condition)

	return b
}

func (b *JQLBuilder) NotIn(field string, values []string) *JQLBuilder {
	quotedValues := make([]string, len(values))

	for i, v := range values {
		quotedValues[i] = fmt.Sprintf("\"%s\"", v)
	}

	condition := fmt.Sprintf("%s NOT IN (%s)", field, strings.Join(quotedValues, ", "))
	b.conditions = append(b.conditions, condition)

	return b
}

func (b *JQLBuilder) Build() string {
	if len(b.conditions) == 0 {
		return ""
	}

	return strings.Join(b.conditions, " AND ")
}
