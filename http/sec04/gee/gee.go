package gee

import (
	"log"
	"net/http"
)

// RouterGroup engine上包装了一层RouteGroup
type RouterGroup struct {
	prefix string
	// 支持中间件
	middlewares []HandlerFunc
	// 支持嵌套
	parent *RouterGroup
	// 所有的Group共享一个engine实例
	engine *Engine
}

// Engine 拥有RouteGroup的所有能力
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

// New 此时的Engine是rg父类 因为其中包含了rg
// engine可以调用子类的方法 可以调用完new 再group
func New() *Engine {
	// 初始化 engine
	engine := &Engine{router: newRouter()}
	// 初始化rg
	engine.RouterGroup = &RouterGroup{engine: engine}
	// 初始化group
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {

	// 创建新的group时 提取出其中的engine
	engine := group.engine

	// 创建一个新的rg 其中的前缀 父级关系 还有
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	// 每个前缀都增加进engine的group
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 不同的group增加路由时 自动增加进对应的engin 多个group共用一个engine
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}
