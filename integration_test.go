package gg

import (
  "os"
  "io"
  "fmt"
  "bytes"
  "testing"
  "reflect"
  "net/http"
  "net/http/httptest"
  "encoding/json"
  "gorm.io/gorm"
  "gorm.io/gorm/schema"
  "github.com/gin-gonic/gin"
)

type TestModel struct {
  gorm.Model
  Name string
  Rating float64
}

var testConfig = Config{
  Gorm: &gorm.Config{},
  Database: InMemConfig{},
  Models: []interface{}{&TestModel{}},
}

func TestMain(m *testing.M) {
  gin.SetMode(gin.TestMode)
  initDB()
  os.Exit(m.Run())
}

func TestSetModel(t *testing.T) {
  router := getRouter()
  var model interface{}
  var sch schema.Schema
  router.GET("/", SetModel(&TestModel{}), func (ctx *gin.Context) {
    model = GetModel(ctx)
    sch = GetModelSchema(ctx)
  })
  runRequest(router, "GET", "/", nil)

  assertEqual(t, reflect.TypeOf(model), reflect.TypeOf(&TestModel{}))
  assertEqual(t, sch.Name, "TestModel")
}

func TestBodyCreate(t *testing.T) {
  router := getRouter()
  router.POST("/", SetModel(&TestModel{}), BodyCreate())

  body := buildJsonBytes("name", "tester", "rating", 0.0)
  resp := runRequest(router, "POST", "/", body)
  jsonResp := extractMapBody(resp)
  expected := buildMap("name", "tester", "rating", 0.0)

  assertEqual(t, 201, resp.Code)
  assertNotEqual(t, 0, jsonResp["id"])
  assertMapEqual(t, expected, jsonResp)
  t.Logf("%v", jsonResp)
}

func TestBodyCreateMissingArgs(t *testing.T) {
  router := getRouter()
  router.POST("/", SetModel(&TestModel{}), BodyCreate())

  body := buildJsonBytes("name", "tester")
  resp := runRequest(router, "POST", "/", body)
  jsonResp := extractMapBody(resp)
  expected := buildMap("error", "rating is required")

  assertEqual(t, 400, resp.Code)
  assertMapEqual(t, expected, jsonResp)
  t.Logf("%v", jsonResp)
}

func TestDeleteByID(t *testing.T) {
  router := getRouter()
  router.DELETE("/:id", SetModel(&TestModel{}), DeleteByID("id"))
  id := create("ponce", 0.0).ID

  resp := runRequest(router, "DELETE", fmt.Sprintf("/%d", id), nil)

  assertEqual(t, resp.Code, 200)
}

func TestDeleteByIDBadID(t *testing.T) {
  router := getRouter()
  router.DELETE("/:id", SetModel(&TestModel{}), DeleteByID("id"))

  resp := runRequest(router, "DELETE", "/4000", nil)

  assertEqual(t, 404, resp.Code)
}

func TestMutateByID(t *testing.T) {
  router := getRouter()
  router.GET("/:id", SetModel(&TestModel{}), MutateByID("id", func(m interface{}) interface{} {
    n := m.(*TestModel)
    n.Name += " world"
    return n
  }))
  id := create("hello", 0.0).ID

  resp := runRequest(router, "GET", fmt.Sprintf("/%d", id), nil)
  jsonResp := extractMapBody(resp)
  expected := buildMap("name", "hello world")

  assertEqual(t, 200, resp.Code)
  assertMapEqual(t, expected, jsonResp)
  t.Logf("%v", jsonResp)
}

func TestMutateByIDMissingID(t *testing.T) {
  router := getRouter()
  router.GET("/:id", SetModel(&TestModel{}), MutateByID("id", func(m interface{}) interface{} {
    n := m.(*TestModel)
    n.Name += " world"
    return n
  }))

  resp := runRequest(router, "GET", "/8000", nil)
  jsonResp := extractMapBody(resp)

  assertEqual(t, 404, resp.Code)
  t.Logf("%v", jsonResp)
}

func TestQuerySearchSimple(t *testing.T) {
  router := getRouter()
  router.GET("/testmodel", SetModel(&TestModel{}), QuerySearch())

  resp := runRequest(router, "GET", "/testmodel?name=harry", nil)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 2, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestQuerySearchMultiple(t *testing.T) {
  router := getRouter()
  router.GET("/testmodel", SetModel(&TestModel{}), QuerySearch())

  resp := runRequest(router, "GET", "/testmodel?rating_gt=2.0&rating_lte=5.0", nil)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 5, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestQuerySearchDates(t *testing.T) {
  router := getRouter()
  router.GET("/testmodel", SetModel(&TestModel{}), QuerySearch())

  resp := runRequest(router, "GET", "/testmodel?created_at=2010-10-10", nil)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 0, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchSimple(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel", SetModel(&TestModel{}), BodySearch())

  body := buildJsonBytes("name", "harry")
  resp := runRequest(router, "POST", "/testmodel", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 2, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchMultiple(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel", SetModel(&TestModel{}), BodySearch())

  body := buildJsonBytes("rating_gt", 2.0, "rating_lte", 5.0)
  resp := runRequest(router, "POST", "/testmodel", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 5, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchNoResults(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel", SetModel(&TestModel{}), BodySearch())

  body := buildJsonBytes("name", "doesnt exist")
  resp := runRequest(router, "POST", "/testmodel", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 0, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchNonmatching(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel", SetModel(&TestModel{}), BodySearch())

  body := buildJsonBytes("nomatch", 0)
  resp := runRequest(router, "POST", "/testmodel", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 400, resp.Code)
  t.Logf("%v", jsonResp)
}

func TestBodySearchNonmatchingWithMatchgin(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel", SetModel(&TestModel{}), BodySearch())

  body := buildJsonBytes("name", "harry", "nomatch", 0)
  resp := runRequest(router, "POST", "/testmodel", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 2, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchGreaterThan(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel", SetModel(&TestModel{}), BodySearch())

  body := buildJsonBytes("rating_gt", 5.0)
  resp := runRequest(router, "POST", "/testmodel", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 3, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchGreaterThanEqual(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel", SetModel(&TestModel{}), BodySearch())

  body := buildJsonBytes("rating_gte", 8.0)
  resp := runRequest(router, "POST", "/testmodel", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 3, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchLessThan(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel", SetModel(&TestModel{}), BodySearch())

  body := buildJsonBytes("rating_lt", 0.0)
  resp := runRequest(router, "POST", "/testmodel", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 3, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchLessThanEqual(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel", SetModel(&TestModel{}), BodySearch())

  body := buildJsonBytes("rating_lte", -0.2)
  resp := runRequest(router, "POST", "/testmodel", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 3, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchNotEqual(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel", SetModel(&TestModel{}), BodySearch())

  body := buildJsonBytes("rating_ne", 0.0)
  resp := runRequest(router, "POST", "/testmodel", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 14, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestQuerySearchByColumnSimple(t *testing.T) {
  router := getRouter()
  router.GET("/:column", SetModel(&TestModel{}), QuerySearchByColumn("column", "name"))

  resp := runRequest(router, "GET", "/ross", nil)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 2, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchByColumnSimple(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel/:column", SetModel(&TestModel{}), BodySearchByColumn("column", "name"))
  ost1id := create("ostio", 0.0).ID

  body := buildJsonBytes("id", ost1id)
  resp := runRequest(router, "POST", "/testmodel/ostio", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 1, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchByColumnComparison(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel/:column", SetModel(&TestModel{}), BodySearchByColumn("column", "rating_ne"))

  body := buildJsonBytes("name", "greco")
  resp := runRequest(router, "POST", "/testmodel/0", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 2, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestGetByID(t *testing.T) {
  router := getRouter()
  router.GET("/:id", SetModel(&TestModel{}), GetByID("id"))
  id := create("bunko", 0.0).ID

  resp := runRequest(router, "GET", fmt.Sprintf("/%d", id), nil)
  jsonResp := extractMapBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, "bunko", jsonResp["name"])
  assertEqual(t, id, jsonResp["id"])
  t.Logf("%v", jsonResp)
}

func TestGetByIDMissingID(t *testing.T) {
  router := getRouter()
  router.GET("/:id", SetModel(&TestModel{}), GetByID("id"))

  resp := runRequest(router, "GET", fmt.Sprintf("/%d", 5000), nil)
  jsonResp := extractMapBody(resp)

  assertEqual(t, 404, resp.Code)
  t.Logf("%v", jsonResp)
}

func TestUpdateByID(t *testing.T) {
  router := getRouter()
  router.PUT("/:id", SetModel(&TestModel{}), UpdateByID("id"))
  id := create("cool", 0.0).ID

  body := buildJsonBytes("name", "runnings")
  resp := runRequest(router, "PUT", fmt.Sprintf("/%d", id), body)
  jsonResp := extractMapBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, "runnings", jsonResp["name"])
  assertEqual(t, id, jsonResp["id"])
  t.Logf("%v", jsonResp)
}

func TestUpdateByIDBadArgs(t *testing.T) {
  router := getRouter()
  router.PUT("/:id", SetModel(&TestModel{}), UpdateByID("id"))
  id := create("brick", 0.0).ID

  body := buildJsonBytes("lamp", "ramp")
  resp := runRequest(router, "PUT", fmt.Sprintf("/%d", id), body)
  jsonResp := extractMapBody(resp)

  _, hasErrorMsg := jsonResp["error"]
  assertEqual(t, 400, resp.Code)
  assertEqual(t, true, hasErrorMsg)
  t.Logf("%v", jsonResp)
}

func TestUpdateByIDBadID(t *testing.T) {
  router := getRouter()
  router.PUT("/:id", SetModel(&TestModel{}), UpdateByID("id"))

  body := buildJsonBytes("rating", 1.9)
  resp := runRequest(router, "PUT", "/9000", body)
  jsonResp := extractMapBody(resp)

  _, hasErrorMsg := jsonResp["error"]
  assertEqual(t, 404, resp.Code)
  assertEqual(t, true, hasErrorMsg)
  t.Logf("%v", jsonResp)
}

func create(name string, rating float64) TestModel {
  db, _ := testConfig.OpenDB()
  tm := TestModel{ Name: name , Rating: rating }
  db.Create(&tm)
  return tm
}

func extractSliceBody(resp *httptest.ResponseRecorder) []interface{} {
  var jsonBody []interface{}
  respBody, _ := io.ReadAll(resp.Body)
  json.Unmarshal(respBody, &jsonBody)
  return jsonBody
}

func extractMapBody(resp *httptest.ResponseRecorder) map[string]interface{} {
  jsonBody := make(map[string]interface{})
  respBody, _ := io.ReadAll(resp.Body)
  json.Unmarshal(respBody, &jsonBody)
  return jsonBody
}

func buildJsonBytes(args... interface{}) []byte {
  jsonMap := buildMap(args...)
  json, _ := json.Marshal(jsonMap)
  return json
}

func buildMap(args... interface{}) map[string]interface{} {
  m := make(map[string]interface{})
  for i := 0; i < len(args); i+=2 {
    m[args[i].(string)] = args[i+1]
  }
  return m
}

func runRequest(router *gin.Engine, method, url string, bodyBytes []byte) *httptest.ResponseRecorder{
  var body io.Reader
  if bodyBytes != nil {
      body = bytes.NewBuffer(bodyBytes)
  }
  rec := httptest.NewRecorder()
  req, _ := http.NewRequest(method, url, body)
  router.ServeHTTP(rec, req)
  return rec
}

func initDB() {
  db, _ := testConfig.OpenDB()
  db.AutoMigrate(testConfig.Models...)
  create("harry",  5.0 )
  create("harry",  2.2 )
  create("delia",  4.0 )
  create("greco",  9.0 )
  create("greco",  0.0 )
  create("greco",  8.1 )
  create("manny",  8.5 )
  create("fareek", -0.2 )
  create("joanne", 1.0 )
  create("tregan", 5.0 )
  create("ross",   2.0 )
  create("ross",   -3.0 )
  create("chris",  1.0 )
  create("ostio",  -3.0 )
  create("reg",    4.0 )
}

func getRouter() *gin.Engine {
  r := gin.New()
  r.Use(Middleware(testConfig))
  return r
}

func assertMapEqual(t *testing.T, exp, act map[string]interface{}) {
  for k, v := range exp {
    if ! isEquivalent(v, act[k]) {
      t.Logf("Expected %v != Actual %v", exp, act)
      t.Fail()
    }
  }
}

func assertEqual(t *testing.T, exp, act interface{}) {
  if ! isEquivalent(exp, act) {
    t.Logf("Expected: %v != Actual %v", exp, act)
    t.Fail()
  }
}

func assertNotEqual(t *testing.T, exp, act interface{}) {
  if isEquivalent(exp, act) {
    t.Logf("Expected: %v != Actual %v", exp, act)
    t.Fail()
  }
}

func isEquivalent(exp, act interface{}) bool{
  switch {
    case act == exp:
      return true
    case act == nil && exp != nil:
      return false
  }
  expType := reflect.TypeOf(exp)
  return reflect.ValueOf(act).Convert(expType).Interface() == exp
}
