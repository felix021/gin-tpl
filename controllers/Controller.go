package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var camel = regexp.MustCompile("(^[^A-Z0-9]*|[A-Z0-9]*)([A-Z0-9][^A-Z]+|$)")
var supportedHttpMethods = []string{"Get", "Post", "Put", "Delete", "Head", "Patch", "Options"}

type Config struct {
	HttpMethod string
	Path       string
	Handler    gin.HandlerFunc
}

func CamelToDash(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}

		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, "-"))
}

func FieldExists(c interface{}, fieldName string) bool {
	return reflect.ValueOf(c).Elem().FieldByName(fieldName) != reflect.Value{}
}

func getRegisterConfig(c interface{}) []Config {
	config := []Config{}

	rvalue := reflect.ValueOf(c)
	rtype := reflect.TypeOf(c)
	controllerName := strings.TrimSuffix(rtype.Name(), "Controller")
	prefix := "/" + CamelToDash(controllerName)

	for i := 0; i < rtype.NumMethod(); i++ {
		t_method := rtype.Method(i)
		name := t_method.Name

		v_method := rvalue.MethodByName(name)
		result := v_method.Call([]reflect.Value{})
		handler := result[0].Interface().(gin.HandlerFunc)

		params := ""
		if len(result) > 1 {
			params = "/" + strings.Trim(result[1].Interface().(string), "/")
		}

		for _, httpMethod := range supportedHttpMethods {
			if !strings.HasPrefix(name, httpMethod) {
				continue
			}
			subPath := CamelToDash(strings.TrimPrefix(name, httpMethod))
			path := fmt.Sprintf("%s/%s%s", prefix, subPath, params)

			if subPath == "index" {
				config = append(config, Config{HttpMethod: strings.ToUpper(httpMethod), Path: path, Handler: handler})
				path = prefix + params
			}

			config = append(config, Config{HttpMethod: strings.ToUpper(httpMethod), Path: path, Handler: handler})
		}
	}
	return config
}

func RegisterController(router interface{}, c interface{}) {
	r, ok := router.(*gin.Engine)
	if ok {
		for _, config := range getRegisterConfig(c) {
			r.Handle(config.HttpMethod, config.Path, config.Handler)
		}
		return
	}

	rg, ok := router.(*gin.RouterGroup)
	if ok {
		for _, config := range getRegisterConfig(c) {
			rg.Handle(config.HttpMethod, config.Path, config.Handler)
		}
		return
	}

	fmt.Println("RegisterController: invalid router")
	os.Exit(-1)
}
