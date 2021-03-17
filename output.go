package gg

import (
  "github.com/gin-gonic/gin"
)

type OutputFormat uint8

//Output formats
const (
  JSON = iota
  IndentedJSON
  YAML
  SecureJSON
  JSONP
  AsciiJSON
  PureJSON
)

type output func(int, interface{})

type contextOutput func(*gin.Context) output

func DefaultOutput(ctx *gin.Context, code int, obj interface{}) {
  GetDefaultOutputFunction(ctx)(code,obj)
}

func Output(ctx *gin.Context, format OutputFormat, code int, obj interface{}){
  GetOutputFunction(format)(ctx)(code,obj)
}

func GetOutputFunction(format OutputFormat) contextOutput {
  switch format {
    case JSON:
      return json
    case IndentedJSON:
      return indentedjson
    case YAML:
      return yaml
    case SecureJSON:
      return securejson
    case JSONP:
      return jsonp
    case AsciiJSON:
      return asciijson
    case PureJSON:
      return purejson
    default:
      return json
  }
}

func json(ctx *gin.Context) output { return ctx.JSON }
func indentedjson(ctx *gin.Context) output { return ctx.IndentedJSON }
func yaml(ctx *gin.Context) output { return ctx.YAML }
func securejson(ctx *gin.Context) output { return ctx.SecureJSON }
func jsonp(ctx *gin.Context) output { return ctx.JSONP }
func asciijson(ctx *gin.Context) output { return ctx.AsciiJSON }
func purejson(ctx *gin.Context) output { return ctx.PureJSON }
