# API

## Fiber
#### 1. New
**Signature**
```
func New(config... Config) *App
```
**Example**
```
// Default config
app := fiber.New()

// Custom config
app := fiber.New(fiber.Config{
    Prefork:       true,
    CaseSensitive: true,
    StrictRouting: true,
    ServerHeader:  "Fiber",
    AppName: "Test App v1.0.1",
})
```
*[Config fields](https://docs.gofiber.io/api/fiber)*
#### 2.NewError
NewError tạo một HTTPError mới với các message tùy chọn
**Signature**
```
func NewError(code int, message ...string) *Error
```
**Example**
```
app.Get("/", func(c *fiber.Ctx) error {
    return fiber.NewError(782, "Custom error message")
})
```

#### 4. IsChild
IsChild xác định xem tiến trình hiện tại có phải kết quả của Prefork hay không
*Prefork: cho phép sử dụng SO_REUSEPORT socket. Socket này cho phép nghe nhiều socket trên cùng một địa chỉ IP và kết hợp cổng.*
**Signature**
```
func IsChild() bool
```
**Example**
```
// Prefork sẽ sinh ra các tiến trình con
app := fiber.New(fiber.Config{
    Prefork: true,
})

if !fiber.IsChild() {
    fmt.Println("I'm the parent process")
} else {
    fmt.Println("I'm a child process")
}

// ...
```

## App
#### 1. Static
Sử dụng phương thức Static cho các file tĩnh như images, CSS và JS
**Signature**
```
func (app *App) Static(prefix, root string, config ...Static) Router
```
Ví dụ với các file tĩnh trong directory ```./public```
```
app.Static("/", "./public")
// => http://localhost:3000/hello.html
// => http://localhost:3000/js/jquery.js
// => http://localhost:3000/css/style.css

app.Static("/static", "./public")

// => http://localhost:3000/static/hello.html
// => http://localhost:3000/static/js/jquery.js
// => http://localhost:3000/static/css/style.css
```
Có thể sửa cấu trúc ```fiber.Static``` như sau:
```
type Static struct{
    //Nếu true, máy chủ sẽ cố gắng giảm thiểu mức sử dụng CPU bằng cách lưu vào bộ đệm các tệp nén
    //Giá trị mặc định: False
    Compress bool `json:"compress"`
    
    //Nếu true, requests phạm vi byte sẽ được enable
    //Default value: False
    ByteRange bool `json:"byte_range"`
    
    //Nếu true, cho phép duyệt thư mục
    //Default value: False
    Browse bool `json:"browse"`
    
    //Nếu true, cho phép download trực tiếp
    //Default value: False
    Download bool `json:"download"`
    
    //Tên của file chỉ mục để serve một directory
    //Default value: "index.html"
    Index string `json:"index"`
    
    //Giới hạn thời gian với trình xử lý file không hoạt động
    //Sử dụng time.Duration giá trị âm để disnable
    //Default value: 10*time.Second
    CacheDuration time.Duration `json:"cache_duration"`
    
    //Giá trị của Cache-Control HTTP-header
    //Được cài đặt trong file response. MaxAge được định nghĩa ở mức giây
    //Default value: 0
    MaxAge int `json:"max_age"`
    
    //ModifyResponse định nghĩa hàm cho phép bạn thay đổi các response
    //Default value: nil
    ModifyResponse Handler
    
    //Next định nghãi hàm để bỏ qua các middleware khi return true
    //Default value: nil
    Next func(c *Ctx) bool
}
```

**Example**
```
// Custom config
app.Static("/", "./public", fiber.Static{
  Compress:      true,
  ByteRange:     true,
  Browse:        true,
  Index:         "john.html",
  CacheDuration: 10 * time.Second,
  MaxAge:        3600,
})
```

#### 2. Route Handlers
Khởi tạo một tuyến đường được liên kết với một phương thức HTTP cụ thể.
**Signature**
```
// Các phương thức HTTP
func (app *App) Get(path string, handlers ...Handler) Router
func (app *App) Head(path string, handlers ...Handler) Router
func (app *App) Post(path string, handlers ...Handler) Router
func (app *App) Put(path string, handlers ...Handler) Router
func (app *App) Delete(path string, handlers ...Handler) Router
func (app *App) Connect(path string, handlers ...Handler) Router
func (app *App) Options(path string, handlers ...Handler) Router
func (app *App) Trace(path string, handlers ...Handler) Router
func (app *App) Patch(path string, handlers ...Handler) Router

// ```Add``` cho phép chỉ định một phương thức làm giá trị
func (app *App) Add(method, path string, handlers ...Handler) Router

// All khởi tạo một route trên tất cả phương thức HTTP
// Gần gống như app.Use nhưng không bị bó buộc bởi các tiền tố (gì gì đó cái đoạn ni không hiểu lắm :D)
func (app *App) All(path string, handlers ...Handler) Router
```
**Example**
```
// Simple GET handler
app.Get("/api/list", func(c *fiber.Ctx) error {
  return c.SendString("I'm a GET request!")
})

// Simple POST handler
app.Post("/api/register", func(c *fiber.Ctx) error {
  return c.SendString("I'm a POST request!")
})
```

**Use**: kiểu ``./john`` sẽ math với cả ``/john/doe, /johnnnnnn``...
**Signature**
```
func (app *App) Use(args ...interface{}) Router
```
**Example**
```
// Match với bất kì request nào
app.Use(func(c *fiber.Ctx) error {
    return c.Next()
})

// Match với các requests bắt đầu bằng /api
app.Use("/api", func(c *fiber.Ctx) error {
    return c.Next()
})

// Match với các requests bắt đầu bằng /api hoặc /home (multiple-prefix support)
app.Use([]string{"/api", "/home"}, func(c *fiber.Ctx) error {
    return c.Next()
})

// Ping tới những handlers khác
app.Use("/api", func(c *fiber.Ctx) error {
  c.Set("X-Custom-Header", random.String(32))
    return c.Next()
}, func(c *fiber.Ctx) error {
    return c.Next()
})
```

#### 3. Mount
**Signature**
```
func (a *App) Mount(prefix string, app *App) Router
```
=> app là sub-app của a. Nếu truy cập vào đường dẫn có tiền đố /prefix thì sẽ được định hướng tới app để xử lý.
**Example**
```
func main() {
    app := fiber.New()
    micro := fiber.New()
    app.Mount("/john", micro) // GET /john/doe -> 200 OK

    micro.Get("/doe", func(c *fiber.Ctx) error {
        return c.SendStatus(fiber.StatusOK)
    })

    log.Fatal(app.Listen(":3000"))
}
```

#### 4. MountPath
MountPath là một cách để gắn ứng dụng con vào một đường dẫn cụ thể bên trong ứng dụng chính, mà không yêu cầu tiền tố như Mount.
**Signature**
```
func (app *App) MountPath() string
```
**Example**
```
app.MountPath("/public", publicApp)
```
=> Gắn ứng dụng con ```publicApp``` vào ứng dụng chính ```app``` tại đường dẫn ```/public```. Khi yêu cầu đến ```/public/file``` được gửi, nó sẽ được xử lý bởi ứng dụng con ```publicApp```
#### 5. Group
Nhóm các path và xử lý yêu cầu cho chúng trong một ứng dụng Fiber chính hoặc ứng dụng con.
**Signature**
```
func (app *App) Group(prefix string, handlers ...Handler) Router
```
**Example**
```
//Tạo nhóm các path có tiền tố /api và sử dụng middleware handler để xử lý tất cả các yêu cầu trong nhóm /api
  api := app.Group("/api", handler)  // /api

  v1 := api.Group("/v1", handler)   // /api/v1
  v1.Get("/list", handler)          // /api/v1/list
  v1.Get("/user", handler)          // /api/v1/user

```

#### 6. Route
Bạn có thể định nghĩa routes với tiền tố chung bên trong hàm chung
**Signature**
```
func (app *App) Route(prefix string, fn func(router Router), name ...string) Router
```
**Example**
```
app.Route("/test", func(api fiber.Router) {
      api.Get("/foo", handler).Name("foo") // /test/foo (name: test.foo)
    api.Get("/bar", handler).Name("bar") // /test/bar (name: test.bar)
  }, 
  ```
  
  #### 7. Server
  Server trả về một fasthttp server cơ bản
  **Signature**
  ```
  func (app *App) Server() *fasthttp.Server
  ```
  **Example**
  ```
  //Máy chủ Fasthttp cho phép mỗi địa chỉ IP chỉ có 1 kết nối đến máy chủ tại 1 thời điểm
  app.Server().MaxConnsPerIP = 1
  ```
  
  #### 8. Server Shutdown
  Server Shutdown tắt server mà không làm ảnh hưởng tới bất kì kết nối đang hoạt động nào. Trước tiên, nó đóng tất cả các trình nghe đang mở, sau đó đợi vô thời hạn cho tới khi tất cả kết nối trở về trạng thái không hoạt động. 
  ShutdownWithTimeout sẽ tắt mọi kết nối đang hoạt động sau khi hết tgian chờ.
  ShutdownWithContext tắt server nếu vượt quá thời hạn context
  **Signature**
  ```
  func (app *App) Shutdown() error
func (app *App) ShutdownWithTimeout(timeout time.Duration) error
func (app *App) ShutdownWithContext(ctx context.Context) error
```

#### 9. HandlersCount
Trả về lượng handlers đã được đăng kí
**Signature**
```
func (app *App) HandlersCount() uint32
```

#### 10. Stack
Truy xuất thông tin về trạng thái hiện tại của ứng dụng hoặc máy chủ Fiber
**Signature**
```
func (app *App) Stack() [][]*Route
```
**Example**
```
app.Get("/john/:age", handler)
    app.Post("/register", handler)

    data, _ := json.MarshalIndent(app.Stack(), "", "  ")
    fmt.Println(string(data))
```
**Result**
```
[
  [
    {
      "method": "GET",
      "path": "/john/:age",
      "params": [
        "age"
      ]
    }
  ],
  [
    {
      "method": "HEAD",
      "path": "/john/:age",
      "params": [
        "age"
      ]
    }
  ],
  [
    {
      "method": "POST",
      "path": "/register",
      "params": null
    }
  ]
]
```