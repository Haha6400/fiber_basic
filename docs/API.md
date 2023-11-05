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

// ```All``` khởi tạo một route trên tất cả phương thức HTTP
// Gần gống như ```app.Use``` nhưng không bị bó buộc bởi các tiền tố (gì gì đó cái đoạn ni không hiểu lắm :D)
func (app *App) All(path string, handlers ...Handler) Router
```