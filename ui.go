package main

import (
	"auto-api-ui/store"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

var htmlTplEngine = template.New("htmlTplEngine")

func initUI(router *httprouter.Router) {
	uiModules := loadTemplates(path.Clean(Args.UIPath))

	if len(uiModules) == 0 {
		log.Printf("not found UI module")
		return
	}
	confUiRouter(uiModules, router)
}

func confUiRouter(uiModules []string, router *httprouter.Router) {
	fileServer := http.FileServer(http.Dir(Args.UIPath))
	for i := 0; i < len(uiModules); i++ {
		module := uiModules[i]
		router.Handle("GET", "/"+module+"/*uriPath", moduleHandle(module, fileServer))
	}
}

func moduleHandle(module string, fileServer http.Handler) func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		data := map[string]interface{}{
			"Urls": store.APIs(module),
		}
		indexTpl := module + "/index"
		uriPath := params.ByName("uriPath")
		if uriPath == "/" {
			err := renderTpl(indexTpl, data, writer)
			if err != nil {
				log.Printf("%s render error： %v", indexTpl, err)
			}
			return
		} else if strings.HasSuffix(uriPath, ".html") {
			tpl := module + strings.Replace(path.Clean(uriPath), ".html", "", 1)
			err := renderTpl(tpl, data, writer)
			if err != nil {
				log.Printf("%s render error： %v", indexTpl, err)
			}
			return
		}
		fileServer.ServeHTTP(writer, request)
	}
}

// loadTemplates 加载模板
// 返回包含 index.gohtml 的一级目录列表（包含 index.gohtml 的一级目录就是一个 UI 模板）
func loadTemplates(uiPath string) []string {
	pattern := path.Join(uiPath, "*/*.gohtml")
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("tmeplate match fail")
	}
	if len(filenames) == 0 {
		log.Printf("not match any tmeplate")
	}

	var modules []string

	for _, filename := range filenames {

		name := strings.Replace(filename, uiPath+"/", "", 1)
		name = strings.Replace(name, ".gohtml", "", 1)

		b, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("load template [%s] from %s error：%v", name, filenames, err)
		}
		s := string(b)

		_, err = htmlTplEngine.New(name).Parse(s)
		if err != nil {
			log.Printf("load template [%s] from %s error：%v", name, filenames, err)
		}
		log.Printf("load template [%s] from %s success", name, filename)

		part := strings.Split(name, "/")
		if len(part) == 2 && part[1] == "index" {
			modules = append(modules, part[0])
		}
	}
	log.Printf("load template [%s] success.", htmlTplEngine.DefinedTemplates())
	return modules
}

func renderTpl(t string, data interface{}, writer io.Writer) error {
	return htmlTplEngine.ExecuteTemplate(writer, t, data, )
}
