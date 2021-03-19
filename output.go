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

type output = func(int, interface{})

type contextOutput = func(*gin.Context) output

func DefaultOutput(ctx *gin.Context, code int, obj interface{}) {
  GetDefaultOutputFunction(ctx)(code,obj)
}

func Output(ctx *gin.Context, format OutputFormat, code int, obj interface{}){
  GetOutputFunction(format)(ctx)(code,obj)
}

func GetOutputFunction(format OutputFormat) contextOutput {
  switch format {
    case JSON:
      return outputJson
    case IndentedJSON:
      return outputIndentedjson
    case YAML:
      return outputYaml
    case SecureJSON:
      return outputSecurejson
    case JSONP:
      return outputJsonp
    case AsciiJSON:
      return outputAsciijson
    case PureJSON:
      return outputPurejson
    default:
      return outputJson
  }
}

func outputJson(ctx *gin.Context) output { return ctx.JSON }
func outputIndentedjson(ctx *gin.Context) output { return ctx.IndentedJSON }
func outputYaml(ctx *gin.Context) output { return ctx.YAML }
func outputSecurejson(ctx *gin.Context) output { return ctx.SecureJSON }
func outputJsonp(ctx *gin.Context) output { return ctx.JSONP }
func outputAsciijson(ctx *gin.Context) output { return ctx.AsciiJSON }
func outputPurejson(ctx *gin.Context) output { return ctx.PureJSON }
