package query

import (
	"fmt"
	"strings"
)

type Query interface {
	Build() string
}

type JQLQueryBuilder struct {
	conditions []string
}

func NewJQLQuery() JQLQueryBuilder {
	return JQLQueryBuilder{}
}

func (b JQLQueryBuilder) Equals(field, value string, isFunction bool) JQLQueryBuilder {
	if isFunction {
		b.conditions = append(b.conditions, fmt.Sprintf("%s = %s", field, value))
		return b
	}

	b.conditions = append(b.conditions, fmt.Sprintf("%s = \"%s\"", field, value))
	return b
}

func (b JQLQueryBuilder) In(field string, values []string) JQLQueryBuilder {
	quotedValues := make([]string, len(values))

	for i, v := range values {
		quotedValues[i] = fmt.Sprintf("\"%s\"", v)
	}

	condition := fmt.Sprintf("%s IN (%s)", field, strings.Join(quotedValues, ", "))
	b.conditions = append(b.conditions, condition)

	return b
}

func (b JQLQueryBuilder) NotIn(field string, values []string) JQLQueryBuilder {
	quotedValues := make([]string, len(values))

	for i, v := range values {
		quotedValues[i] = fmt.Sprintf("\"%s\"", v)
	}

	condition := fmt.Sprintf("%s NOT IN (%s)", field, strings.Join(quotedValues, ", "))
	b.conditions = append(b.conditions, condition)

	return b
}

func (b *JQLQueryBuilder) Build() string {
	if len(b.conditions) == 0 {
		return ""
	}

	return strings.Join(b.conditions, " AND ")
}
