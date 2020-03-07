package main

import (
	"auto-api-ui/store"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/threeq/docs"
	"github.com/threeq/docs/swagger"
	"github.com/threeq/docs/swagger/endpoint"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	Version = "0.0.1"
)

var apiDocDesc *store.DocDesc

func initApi(router *httprouter.Router) {

	endpoints := []*swagger.Endpoint{
		endpoint.Post("/api/doc", "添加文档描述",
			endpoint.Handler(addDocHandler),
			endpoint.Tags("doc"),
			endpoint.Description("添加一个新文档"),
			endpoint.Body(store.DocDesc{}, "文档地址描述", true),
			endpoint.Response(http.StatusCreated, "OK",
				"Created", "接口调用成功"),
		),
		endpoint.Put("/api/doc/:name", "修改文档描述",
			endpoint.Handler(updateDocHandler),
			endpoint.Tags("doc"),
			endpoint.Description("修改一个文档描述"),
			endpoint.Path("name", "string", "文档名字", true),
			endpoint.Body(store.DocDesc{}, "文档地址描述", true),
			endpoint.Response(http.StatusOK, "OK",
				"OK", "接口调用成功"),
		),
		endpoint.Delete("/api/doc/:name", "删除文档描述",
			endpoint.Handler(deleteDocHandler),
			endpoint.Tags("doc"),
			endpoint.Description("删除一个文档描述"),
			endpoint.Path("name", "string", "文档名字", true),
			endpoint.Response(http.StatusOK, "OK",
				"OK", "接口调用成功"),
		),
	}

	api := docs.New(
		docs.Version(Version),
		docs.Host(Args.Host),
		docs.BasePath(Args.Path),
		docs.Endpoints(endpoints...),
		docs.Title(Args.Title),
		docs.Description(Args.Description),
		docs.ContactEmail(Args.Contact),
		docs.TermsOfService(Args.TermOfService),
		docs.Tag("doc", "swagger 文档"),
	)

	api.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(http.ResponseWriter, *http.Request, httprouter.Params))
		path = docs.ColonPath(path)

		router.Handle(endpoint.Method, path, h)
	})

	router.HandlerFunc("GET", api.DocEndpointPath(), api.Handler(true))
	apiDocDesc = &store.DocDesc{Name: Args.Title, Type: "swagger", Url: api.DocEndpointPath()}
}

func deleteDocHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	err := store.DelAPI(params.ByName("name"))
	if err != nil {
		writerError(writer, "删除文档-删除错误: %v", err)
	}

	err = writerOK(writer)
	if err != nil {
		writerError(writer, "删除文档-写入返回数据错误: %v", err)
	}
}

func updateDocHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writerError(writer, "修改文档-读取请求错误: %v", err)
		return
	}

	api := &store.DocDesc{}
	err = json.Unmarshal(body, api)
	if err != nil {
		writerError(writer, "修改文档-解析请求错误: %v", err)
		return
	}
	if api.Name == "" {
		api.Name = name
	}

	err = store.UpdateAPI(name, api)
	if err != nil {
		writerError(writer, "修改文档-修改错误: %v", err)
	}
	err = writerOK(writer)
	if err != nil {
		writerError(writer, "修改文档-写入返回数据错误: %v", err)
	}
}

func addDocHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writerError(writer, "增加文档-读取请求错误: %v ", err)
		return
	}

	api := &store.DocDesc{}
	err = json.Unmarshal(body, api)
	if err != nil {
		writerError(writer, "增加文档-解析请求错误: %v", err)
		return
	}

	err = store.AddAPI(api)
	if err != nil {
		writerError(writer, "修改文档-修改错误: %v", err)
	}
	err = writerOK(writer)
	if err != nil {
		writerError(writer, "增加文档-写入返回数据错误: %v", err)
	}
}

func swaggerJsonHandle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

}

func writerError(writer io.Writer, format string, err error) {
	msg := fmt.Sprintf(format, err)
	log.Print(msg)
	_, _ = writer.Write([]byte(msg))
}

func writerOK(writer http.ResponseWriter) error {
	_, err := writer.Write([]byte("OK"))
	return err
}
