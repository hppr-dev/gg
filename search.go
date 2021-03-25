package gg

import (
  "fmt"
  "reflect"
  "time"
  "strings"
  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
  "gorm.io/gorm/schema"
  "hppr.dev/gg/mdl"
)

// Comparison is the basic building block of a search query
type Comparison struct {
  Operator string
  Value interface{}
}

// WhereString creates the substitution string necesarry for a db.Where() call
func (c Comparison) WhereString(param string) string {
  return fmt.Sprintf("%s %s ?", param, c.Operator)
}

type ComparisonMap map[string][]Comparison

func (cm ComparisonMap) toKeyMap() DefaultMap {
  dm := make(DefaultMap)
  for k := range cm {
    dm[k] = true
  }
  return dm
}

// QuerySearch returns a handler suitable for searching using the query parameters in the url
func QuerySearch() gin.HandlerFunc {
  return QuerySearchByColumn("", "")
}

// BodySearch returns a handler suitable for searching using the request post data
func BodySearch() gin.HandlerFunc {
  return BodySearchByColumn("", "")
}

// BodySearchByColumn returns a handler for searching where one column is provided in the url as a url parameter and more search parameters may be provided through the post data
func BodySearchByColumn(urlParam, column string) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    var dataMap DefaultMap
    ctx.BindJSON(&dataMap)
    params := createComparisons(dataMap)
    if column != "" && urlParam != "" {
      params.addComparison(column, ctx.Param(urlParam))
    }
    searchAndOutput(params, ctx)
  }
}

// QuerySearchByColumn returns a handler for searching where one column is provided in the url and other search parameters may be provided through the query string
func QuerySearchByColumn(urlParam, column string) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    sch := GetModelSchema(ctx)
    params := bindQuery(ctx)
    if column != "" && urlParam != "" {
      params.addComparison(column, ctx.Param(urlParam))
    }
    convertComparisonDates(params,sch)
    searchAndOutput(params, ctx)
  }
}

// GetByID returns a handler that searches for a particular record by it's primary key provided through a url parameter
func GetByID(urlParam string) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    db := GetDatabase(ctx)
    model := GetModel(ctx)
    sch := GetModelSchema(ctx)
    pKeyColumn := sch.PrioritizedPrimaryField.DBName
    data := make(DefaultMap)
    result := db.Model(model).Where(pKeyColumn + " = ?", ctx.Param(urlParam)).Find(&data)
    if result.RowsAffected == 0 {
      DefaultOutput(ctx, 404, gin.H{"error": "not found"})
      return
    }
    DefaultOutput(ctx, 200, data)
  }
}

func searchAndOutput(params ComparisonMap, ctx *gin.Context) {
  results, err := search(params, ctx)
  if err != nil {
    DefaultOutput(ctx, 400, gin.H{"error": err.Error()})
    return
  }
  DefaultOutput(ctx, 200, results)
}

func search(params ComparisonMap, ctx *gin.Context) ([]DefaultMap, error) {
  db := GetDatabase(ctx)
  sch := GetModelSchema(ctx)
  model := GetModel(ctx)
  var results []DefaultMap
  if len(params) > 0 {
    if err := mdl.MatchAnyMapToModel(params.toKeyMap(), sch); err != nil {
      return nil, err
    }
    db = buildWhere(db, params, sch)
  }
  db.Model(model).Find(&results)
  if results == nil {
    results = make([]DefaultMap, 0, 0)
  }
  return results, nil
}

func bindQuery(ctx *gin.Context) ComparisonMap {
  m := make(ComparisonMap)
  for k, vs := range ctx.Request.URL.Query() {
    v := strings.Join(vs, "")
    m.addComparison(k, v)
  }
  return m
}

func convertComparisonDates(comparisons ComparisonMap, sch schema.Schema) error {
  var err error
  fields := sch.Fields
  for _, field := range fields {
    if comps, present := comparisons[field.DBName]; present && field.FieldType == reflect.TypeOf(time.Time{}) {
      for _, comp := range comps {
        if comp.Value, err = time.Parse(time.RFC3339Nano, comp.Value.(string)); err != nil {
          return err
        }
      }
    }
  }
  return nil
}

func createComparisons(m DefaultMap) ComparisonMap {
  comps := make(ComparisonMap)
  for k, v := range m {
    comps.addComparison(k, v)
  }
  return comps
}

func buildWhere(db *gorm.DB, params ComparisonMap, sch schema.Schema) (*gorm.DB) {
  for param, comps := range params {
    if _, exist := sch.FieldsByDBName[param] ; exist {
      for _, comp := range comps {
        db = db.Where(comp.WhereString(param), comp.Value)
      }
    }
  }
  return db
}

func (m ComparisonMap) addComparison(k string, v interface{}) {
  var key, op string
  switch {
    case strings.HasSuffix(k, "_gte"):
      key = strings.TrimSuffix(k, "_gte")
      op = ">="
    case strings.HasSuffix(k, "_lte"):
      key = strings.TrimSuffix(k, "_lte")
      op = "<="
    case strings.HasSuffix(k, "_gt"):
      key = strings.TrimSuffix(k, "_gt")
      op = ">"
    case strings.HasSuffix(k, "_lt"):
      key = strings.TrimSuffix(k, "_lt")
      op = "<"
    case strings.HasSuffix(k, "_ne"):
      key = strings.TrimSuffix(k, "_ne")
      op = "<>"
    case strings.HasSuffix(k, "_contains"):
      key = strings.TrimSuffix(k, "_contains")
      op = "LIKE"
      v = fmt.Sprintf("%%%v%%", v)
    default:
      key = k
      op = "="
  }
  m.addValueToList(key, op, v)
}

func (m ComparisonMap) addValueToList(k, op string, v interface{}) {
  comp := Comparison{Value: v, Operator: op}
  sl, exist := m[k]
  if !exist {
    m[k] = []Comparison{comp}
  } else {
    m[k] = append(sl, comp)
  }
}
