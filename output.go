package gg

import (
	"github.com/gin-gonic/gin"
)

type OutputFormat uint8

//Available Output formats
const (
	// JSON will use ctx.JSON
	JSON = iota
	// IndentedJSON will use ctx.IndentedJSON
	IndentedJSON
	// YAML will use ctx.YAML
	YAML
	// SecureJSON will use ctx.SecureJSON
	SecureJSON
	// JSONP will use ctx.JSONP
	JSONP
	// AsciiJSON will use ctx.AsciiJSON
	AsciiJSON
	// PureJSON will use ctx.PureJSON
	PureJSON
)

type output = func(int, interface{})

type contextOutput = func(*gin.Context) output

// DefaultOutput uses the gg.Config's DefaultOutputFormat to output API information
// i.e. if DefaultOutputFormat is JSON, calling DefaultOutput(ctx, 200, gin.H{"msg": "hello"}) would be equivalent to ctx.JSON(200, gin.H{"msg": "hello"})
func DefaultOutput(ctx *gin.Context, code int, obj interface{}) {
	getDefaultOutputFunction(ctx)(code, obj)
}

func getOutputFunction(format OutputFormat) contextOutput {
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

func outputJson(ctx *gin.Context) output         { return ctx.JSON }
func outputIndentedjson(ctx *gin.Context) output { return ctx.IndentedJSON }
func outputYaml(ctx *gin.Context) output         { return ctx.YAML }
func outputSecurejson(ctx *gin.Context) output   { return ctx.SecureJSON }
func outputJsonp(ctx *gin.Context) output        { return ctx.JSONP }
func outputAsciijson(ctx *gin.Context) output    { return ctx.AsciiJSON }
func outputPurejson(ctx *gin.Context) output     { return ctx.PureJSON }
