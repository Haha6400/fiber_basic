# BASIC FIBER

## 1. Installation
- [Tải](https://go.dev/dl/) và cài đặt Go bản 1.17 trở lên
- Sử dụng ```go get``` command:
```
go get github.com/gofiber/fiber/v2
```

## 2. Zero Allocation
Theo mặt định, các giá trị được trả về từ **fiber.Ctx** không phải là bất biến để Fiber có thể được tối ưu hóa, đạt hiệu suất cao hơn. 
Bạn **chỉ có thể** sử dụng các giá trị context trong trình xử lý, mọi giá trị bạn nhận lại được từ context sẽ được sử dụng lại trong các yêu cầu sau này và có thể được thay đổi.
```
func handler(c *fiber.Ctx) error {
    // Biến chỉ hợp lệ trong trình xử lý này thui
    result := c.Params("foo") 
    // Nếu bạn cần duy trì các giá trị như vậy bên ngoài trình xử lý
    //Hãy tạo các bản sao của bộ đệm cơ bản của chúng bằng cách sử dụng nội dung sao chép:
    buffer := make([]byte, len(result))
    copy(buffer, result)
    resultCopy := string(buffer)
}
```
Để dễ dàng thực hiện các chức năng trên, **gofiber/utils** có một hàm ```CopyString```:
```
app.Get("/:foo", func(c *fiber.Ctx) error {
    // Variable is now immutable
    result := utils.CopyString(c.Params("foo")) 

    // ...
})
```
hoặc sử dụng cài đặt **Immutable** để làm cho tất cả giá trị trả về từ context trở thành giá trị bất biến và có thể sử dụng ở bất cứ đâu:
```
app := fiber.New(fiber.Config{
    Immutable: true,
})
```
*For more in4: [#426](https://github.com/gofiber/fiber/issues/426); [#185](https://github.com/gofiber/fiber/issues/185)*

## 3. Hello, world!
Ví dụ đơn giản về Fiber:
```
package main

import "github.com/gofiber/fiber/v2"

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    app.Listen(":3000")
}
```
```
go run server.go
```

## 4. Định tuyến cơ bản(Basic Routing)
- URI (hoặc path)
- Phương thức HTTP request cụ thể (GET, PUT, POST...)
Cấu trúc routing: 
```
app.Method(path string, ...func(*fiber.Ctx) error)
```
Trong đó:
- ```app``` là một phiên bản của Fiber: ```app := fiber.New()```
- ```Method``` là phương thức HTTP request(GET, PUT, POST, DELETE...)
- ```path``` là một đường dẫn trên server 
- ```func(*fiber.Ctx) error``` là một hàm gọi chứa context khi route matched

***[Ví dụ đơn giản ](https://pastecord.com/zonezyvyte.swift)***

## 5. Static files
```
app.Static(prefix, root string, config... Static)
```
Lưu trữ các files tĩnh ở directory tên ```./public```
```
app := fiber.New()
app.Static("/", "./public")
app.Listen(":3000")
```