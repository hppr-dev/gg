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

type Comparison struct {
  Operator string
  Value interface{}
}

func (c Comparison) WhereString(param string) string {
  return fmt.Sprintf("%s %s ?", param, c.Operator)
}

type ComparisonMap map[string]Comparison

func (cm ComparisonMap) toDefaultMap() DefaultMap {
  dm := make(DefaultMap)
  for k,c := range cm {
    dm[k] = c.Value
  }
  return dm
}

func QuerySearch() gin.HandlerFunc {
  return func(ctx *gin.Context) {
    params := bindQuery(ctx)
    sch := GetModelSchema(ctx)
    convertComparisonDates(params, sch)
    results, err := search(params, ctx)
    if err != nil {
      DefaultOutput(ctx, 400, gin.H{"error": err.Error()})
      return
    }
    DefaultOutput(ctx, 200, results)
  }
}

func BodySearch() gin.HandlerFunc {
  return func(ctx *gin.Context) {
    var dataMap DefaultMap
    ctx.BindJSON(&dataMap)
    params := createComparisons(dataMap)
    results, err := search(params, ctx)
    if err != nil {
      DefaultOutput(ctx, 400, gin.H{"error": err.Error()})
      return
    }
    DefaultOutput(ctx, 200, results)
  }
}

func SearchByColumn(column, urlParam string) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    DefaultOutput(ctx, 200, "Not Implemented")
  }
}

func GetByID(urlParam string) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    db := GetDatabase(ctx)
    model := GetModel(ctx)
    sch := GetModelSchema(ctx)
    pKeyColumn := sch.PrioritizedPrimaryField.DBName
    var results []DefaultMap
    if err := db.Model(model).Where(pKeyColumn + " = ?", ctx.Param(urlParam)).First(&results).Error; err != nil {
      DefaultOutput(ctx, 404, gin.H{"error": err.Error()})
      return
    }
    DefaultOutput(ctx, 200, results[0])
  }
}

func search(params ComparisonMap, ctx *gin.Context) ([]DefaultMap, error) {
  var results []DefaultMap
  db := GetDatabase(ctx)
  sch := GetModelSchema(ctx)
  model := GetModel(ctx)
  if len(params) > 0 {
    if err := mdl.MatchAnyMapToModel(params.toDefaultMap(), sch); err != nil {
      return nil, err
    }
    db = buildWhere(db, params)
  }
  err := db.Model(model).Find(&results).Error
  if results == nil {
    results = make([]DefaultMap, 0, 0)
  }
  return results, err
}

func bindQuery(ctx *gin.Context) ComparisonMap {
  m := make(ComparisonMap)
  for k, vs := range ctx.Request.URL.Query() {
    v := strings.Join(vs, "")
    addComparison(m, k, v)
  }
  return m
}

func convertComparisonDates(comparisons ComparisonMap, sch schema.Schema) error {
  var err error
  fields := sch.Fields
  for _, field := range fields {
    if comp, present := comparisons[field.DBName]; present && field.FieldType == reflect.TypeOf(time.Time{}) {
      if comp.Value, err = time.Parse(time.RFC3339Nano, comp.Value.(string)); err != nil {
        return err
      }
    }
  }
  return nil
}

func createComparisons(m DefaultMap) ComparisonMap {
  comps := make(ComparisonMap)
  for k, v := range m {
    addComparison(comps, k, v)
  }
  return comps
}

func buildWhere(db *gorm.DB, params ComparisonMap) (*gorm.DB) {
  for param, comp := range params {
    db = db.Where(comp.WhereString(param), comp.Value)
  }
  return db
}

func addComparison(m ComparisonMap, k string, v interface{}) {
  switch {
    case strings.HasSuffix(k, "_gte"):
      m[strings.TrimSuffix(k, "_gte")] = Comparison{Value: v, Operator: ">=" }
    case strings.HasSuffix(k, "_lte"):
      m[strings.TrimSuffix(k, "_lte")] = Comparison{Value: v, Operator: "<=" }
    case strings.HasSuffix(k, "_gt"):
      m[strings.TrimSuffix(k, "_gt")] = Comparison{Value: v, Operator: ">" }
    case strings.HasSuffix(k, "_lt"):
      m[strings.TrimSuffix(k, "_lt")] = Comparison{Value: v, Operator: "<"}
    case strings.HasSuffix(k, "_ne"):
      m[strings.TrimSuffix(k, "_ne")] = Comparison{Value: v, Operator: "<>"}
    default:
      m[k] = Comparison{Value: v, Operator: "="}
  }
}