package mdl

import (
	"errors"
	"fmt"
	"gorm.io/gorm/schema"
)

type StringSet map[string]bool

func MatchAnyMapToModel(dataMap map[string]interface{}, schema schema.Schema) error {
	columnSet := getModelColumnNames(schema, true)
	for key := range dataMap {
		if columnSet[key] {
			return nil
		}
	}
	return errors.New("No matchable data")
}

func MatchAllMapToModel(dataMap map[string]interface{}, schema schema.Schema) error {
	columnSet := getModelColumnNames(schema, false)
	for column, required := range columnSet {
		if _, present := dataMap[column]; required && !present {
			return errors.New(fmt.Sprintf("%s is required", column))
		}
	}
	return nil
}

func getModelColumnNames(schema schema.Schema, includeAuto bool) StringSet {
	names := make(StringSet)
	for _, field := range schema.Fields {
		names[field.DBName] = !isAutoCreatable(field) || includeAuto
	}
	return names
}

func isAutoCreatable(f *schema.Field) bool {
	return f.Creatable && (f.HasDefaultValue || f.AutoIncrement || f.AutoCreateTime != 0 || f.AutoUpdateTime != 0 || f.DBName == "deleted_at")
}
