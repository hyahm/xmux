package xmux

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("66666666")
	name := Var(r)["aaa"]
	// time.Sleep(time.Millisecond * 30)
	m1 := map[string]string{
		"message": name,
	}
	b, _ := json.Marshal(m1)
	GetInstance(r).GetFuncName()
	fmt.Println(GetInstance(r).GetUrl())
	w.Write(b)
	// GetInstance(r).Response.(*Response).Msg = time.Now().String()

}

func home2(w http.ResponseWriter, r *http.Request) {
	// name := Var(r)["name"]
	// time.Sleep(time.Millisecond * 30)
	m2 := map[string]string{
		"message": "asdfasdf",
	}
	GetInstance(r).Set("xmux_response", m2)
	w.Write([]byte("aaaaa"))
}

func grouphome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("mmmmmmmm group")
	w.Write([]byte("grouphome" + Var(r)["name"] + "-" + Var(r)["age"]))
}

func adminhandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("mmmmmmmm admin")
}

func DictMe() *RouteGroup {
	dict := NewRouteGroup().SetMeta(Meta{
		Name:     "我是菜单",
		MenuType: "f",
	})
	dict.AddGroup(adminGroup())
	return dict
}

func adminGroup() *RouteGroup {
	admin := NewRouteGroup().Prefix("test").SetMeta(Meta{
		Name:     "后台",
		MenuType: "c",
	})
	admin.Get("/admin/b", home).SetMeta(Meta{
		Name:     "admin",
		MenuType: "b",
	})

	admin.Get("/admin", adminhandle).DelPageKeys("editor")
	admin.Get("/aaa/adf{re:([a-z]{1,4})sf([0-9]{0,10})sd: name, age}", grouphome)
	return admin
}

func userGroup() *RouteGroup {
	user := NewRouteGroup().SetMeta(Meta{
		Name:     "用户管理菜单",
		MenuType: "c",
	})
	// user.Get("/group", home).Use(CombineHandlers())
	user.Get("/user/{asdfsdf}/{int:gg}", home).SetMeta(Meta{
		Name:     "用户查看",
		MenuType: "b",
	})

	user.AddGroup(DictMe()).DelPostModule(postModule)
	return user
}

func Post(w http.ResponseWriter, r *http.Request) (exit bool) {
	fmt.Println("post")
	return false
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func PermissionTemplate(w http.ResponseWriter, r *http.Request) (post bool) {

	// Get the permission of the corresponding URI, which is set by addpagekeys and delpagekeys
	pages := GetInstance(r).GetPageKeys()
	// If the length is 0, it means that anyone can access it
	if len(pages) == 0 {
		return false
	}
	// Get the corresponding role of the user and judge that it is all in
	roles := []string{"editor"} //Obtain the user's permission from the database or redis
	for _, role := range roles {
		if _, ok := pages[role]; ok {
			//What matches here is the existence of this permission. Continue to follow for so long
			return false
		}
	}
	// no permission
	w.Write([]byte("no permission"))
	return true
}

var db *gorm.DB

func InitMySQL() {
	// ====================== 这里改成你的数据库信息 ======================
	username := "testuser"   // 用户名
	password := "123456"     // 密码
	host := "172.21.174.119" // 地址
	port := "3306"           // 端口
	dbName := "test_db"      // 数据库名
	// ==================================================================

	// 拼接 DSN 连接串
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 打印 SQL 语句（开发用，上线可关闭）
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取 DB 实例失败：%v", err)
	}
	sqlDB.SetMaxOpenConns(100)                 // 最大连接数
	sqlDB.SetMaxIdleConns(20)                  // 最大空闲连接
	sqlDB.SetConnMaxLifetime(10 * time.Minute) // 连接最长存活时间

	log.Println("✅ 数据库连接成功")
}

func TestMain(t *testing.T) {
	// pool := NewPool()
	router := NewRouter().AddModule(PermissionTemplate)
	// router.HandleAll = LimitFixedWindowCounterTemplate
	// router.HandleRecover = func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("服务器错误"))
	// }
	// cth := gocache.NewCache[string, []byte](100, gocache.LFU)
	// InitResponseCache(cth)
	// router.AddModule(setkey, DefaultCacheTemplateCacheWithoutResponse).AddModule(Post)
	router.AddPageKeys("admin", "editor")
	// router.SetHeader("Access-Control-Allow-Origin", "*").
	// 	SetHeader("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	// router.AddGroup(Pprof())
	// router.Enter = enter
	// router.ModuleContinue = true
	// router.Prefix("/api")
	// router.EnableConnect = true
	router.Get("/test/{aaaa}", home).SetMeta(Meta{
		Name:     "测试接口",
		MenuType: "b",
	})
	// router.Get("/bar", home2).AddPageKeys("admin")
	// router.Get("/post", pp).Use(pool.Middleware(heavyHandler))
	// pf := router.PageKeyFuncMap()
	// fmt.Println(pf)
	// router.SetAddr(":8080")
	router.AddGroup(userGroup())
	// b, _ := json.MarshalIndent(router.menuTree, "", "  ")
	fmt.Println(len(router.Routes()))
	b, _ := json.MarshalIndent(BuildRouteTree(router.Menus()), "", "  ")
	fmt.Println(string(b))
	log.Fatal(router.SetAddr(":19999").Run())
}

type Binding struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func Recovery(key string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			key := r.URL.Path
			e := getEntry(key)

			mu.Lock()
			//---------- 生产者路径 ----------
			if !e.done {
				e.count++   // 正在计算
				mu.Unlock() // 放锁去 IO

				// 模拟耗时计算
				// resp := []byte("这是 " + key + " 的处理结果")
				next.ServeHTTP(w, r)
				mu.Lock()
				e.response = GetInstance(r).Get(key).([]byte)
				e.done = true
				e.cond.Broadcast()
				mu.Unlock()
				return
			}
			//---------- 消费者路径 ----------
			e.count++
			for !e.done {
				e.cond.Wait() // 内部会临时放锁
			}
			data := append([]byte(nil), e.response...) // 锁内深拷贝
			e.count--
			if e.count == 0 { // 最后一个离开
				delete(entries, key) // 安全删除
			}
			mu.Unlock()

			w.Write(data)

		})
	}
}

// Middleware 返回 http.Handler 中间件
// opts: 可传入 KeyFunc，默认用 r.URL.Path
func (p *Pool) Middleware(next http.HandlerFunc, opts ...KeyFunc) http.HandlerFunc {
	keyFn := func(r *http.Request) string { return r.URL.Path }
	if len(opts) > 0 && opts[0] != nil {
		keyFn = opts[0]
	}

	return func(w http.ResponseWriter, r *http.Request) {
		key := keyFn(r)
		e := p.getEntry(key)

		p.mu.Lock()
		// 生产者路径
		if !e.done {
			e.count++
			p.mu.Unlock()

			// 执行业务 handler 拿结果
			rec := &responseRecorder{ResponseWriter: w, status: 200}
			next(rec, r)

			p.mu.Lock()
			e.response = rec.body
			e.done = true
			e.cond.Broadcast()
			p.mu.Unlock()
			return
		}

		// 消费者路径
		e.count++
		for !e.done {
			e.cond.Wait()
		}
		data := append([]byte(nil), e.response...)
		e.count--
		if e.count == 0 {
			delete(p.data, key)
		}
		p.mu.Unlock()

		// 把缓存内容写回客户端
		w.Write(data)
	}
}

// 内部拿 entry（锁内）
func (p *Pool) getEntry(key string) *entry {
	p.mu.Lock()
	defer p.mu.Unlock()
	e := p.data[key]
	if e == nil {
		e = &entry{cond: sync.NewCond(&p.mu)}
		p.data[key] = e
	}
	return e
}

// KeyFunc 允许调用方自定义分组 key
type KeyFunc func(r *http.Request) string

// responseRecorder 把响应内容截下来复用
type responseRecorder struct {
	http.ResponseWriter
	body   []byte
	status int
}

func (r *responseRecorder) Write(p []byte) (int, error) {
	r.body = append(r.body, p...)
	return len(p), nil
}

func postModule(w http.ResponseWriter, r *http.Request) bool {

	fmt.Println("这是一个后置模块")
	return false
}

// Pool 按 key 缓存请求结果
type Pool struct {
	mu   sync.Mutex
	data map[string]*entry
}

type entry struct {
	done     bool
	response []byte
	count    int32
	cond     *sync.Cond
}

func NewPool() *Pool {
	return &Pool{data: make(map[string]*entry)}
}

var (
	mu      sync.Mutex // 只用一把 Lock，简化生命周期
	entries = make(map[string]*entry)
)

// 拿到或新建 entry（锁内）
func getEntry(key string) *entry {
	mu.Lock()
	defer mu.Unlock()
	e := entries[key]
	if e == nil {
		e = &entry{cond: sync.NewCond(&mu)}
		entries[key] = e
	}
	return e
}

func pp(w http.ResponseWriter, r *http.Request) {

	// a := make([]string, 0)
	key := r.URL.Path
	e := getEntry(key)

	mu.Lock()
	//---------- 生产者路径 ----------
	if !e.done {
		e.count++   // 正在计算
		mu.Unlock() // 放锁去 IO

		// 模拟耗时计算
		resp := []byte("这是 " + key + " 的处理结果")

		mu.Lock()
		e.response = resp
		e.done = true
		e.cond.Broadcast()
		mu.Unlock()
		return
	}
	//---------- 消费者路径 ----------
	e.count++
	for !e.done {
		e.cond.Wait() // 内部会临时放锁
	}
	data := append([]byte(nil), e.response...) // 锁内深拷贝
	e.count--
	if e.count == 0 { // 最后一个离开
		delete(entries, key) // 安全删除
	}
	mu.Unlock()

	w.Write(data)
}
