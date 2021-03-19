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
}

var testConfig = Config{
  Gorm: &gorm.Config{},
  Database: InMemConfig{},
  Models: []interface{}{&TestModel{}},
}

func TestMain(m *testing.M) {
  gin.SetMode(gin.TestMode)
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

  body := buildJsonBytes("name", "tester")
  resp := runRequest(router, "POST", "/", body)
  jsonResp := extractMapBody(resp)
  expected := buildMap("id", 1)

  assertEqual(t, 201, resp.Code)
  assertMapEqual(t, expected, jsonResp)
  t.Logf("%v", jsonResp)
}

func TestDeleteByID(t *testing.T) {
  router := getRouter()
  router.DELETE("/:id", SetModel(&TestModel{}), DeleteByID("id"))
  id := create("ponce").ID

  resp := runRequest(router, "DELETE", fmt.Sprintf("/%d", id), nil)

  assertEqual(t, resp.Code, 200)
}

func TestMutateByID(t *testing.T) {
  router := getRouter()
  router.GET("/:id", SetModel(&TestModel{}), MutateByID("id", func(m interface{}) interface{} {
    n := m.(*TestModel)
    n.Name += " world"
    return n
  }))
  id := create("hello").ID

  resp := runRequest(router, "GET", fmt.Sprintf("/%d", id), nil)
  jsonResp := extractMapBody(resp)
  expected := buildMap("name", "hello world")

  assertEqual(t, 200, resp.Code)
  assertMapEqual(t, expected, jsonResp)
  t.Logf("%v", jsonResp)
}

func TestQuerySearchSimple(t *testing.T) {
  router := getRouter()
  router.GET("/testmodel", SetModel(&TestModel{}), QuerySearch())
  create("harry")
  create("harry")
  create("delia")

  resp := runRequest(router, "GET", "/testmodel?name=harry", nil)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 2, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchSimple(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel", SetModel(&TestModel{}), BodySearch())
  create("greco")
  create("greco")
  create("manny")

  body := buildJsonBytes("name", "harry")
  resp := runRequest(router, "POST", "/testmodel", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 2, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestQuerySearchByColumnSimple(t *testing.T) {
  router := getRouter()
  router.GET("/:column", SetModel(&TestModel{}), QuerySearchByColumn("column", "name"))
  create("ross")
  create("ross")
  create("chris")

  resp := runRequest(router, "GET", "/ross", nil)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 2, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestBodySearchByColumnSimple(t *testing.T) {
  router := getRouter()
  router.POST("/testmodel/:column", SetModel(&TestModel{}), BodySearchByColumn("column", "name"))
  ost1id := create("ostio").ID
  create("ostio")
  create("reg")

  body := buildJsonBytes("id", ost1id)
  resp := runRequest(router, "POST", "/testmodel/ostio", body)
  jsonResp := extractSliceBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, 1, len(jsonResp))
  t.Logf("%v", jsonResp)
}

func TestGetByID(t *testing.T) {
  router := getRouter()
  router.GET("/:id", SetModel(&TestModel{}), GetByID("id"))
  id := create("bunko").ID

  resp := runRequest(router, "GET", fmt.Sprintf("/%d", id), nil)
  jsonResp := extractMapBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, "bunko", jsonResp["name"])
  assertEqual(t, id, jsonResp["id"])
  t.Logf("%v", jsonResp)
}

func TestUpdateByID(t *testing.T) {
  router := getRouter()
  router.PUT("/:id", SetModel(&TestModel{}), UpdateByID("id"))
  id := create("cool").ID

  body := buildJsonBytes("name", "runnings")
  resp := runRequest(router, "PUT", fmt.Sprintf("/%d", id), body)
  jsonResp := extractMapBody(resp)

  assertEqual(t, 200, resp.Code)
  assertEqual(t, "runnings", jsonResp["name"])
  assertEqual(t, id, jsonResp["id"])
  t.Logf("%v", jsonResp)
}

func create(name string) TestModel {
  db, _ := testConfig.OpenDB()
  tm := TestModel{ Name: name }
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

func getRouter() *gin.Engine {
  r := gin.Default()
  if db, err := testConfig.OpenDB(); err == nil {
    db.AutoMigrate(testConfig.Models...)
  }
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
