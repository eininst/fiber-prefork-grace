# Grace

[![Build Status](https://travis-ci.org/ivpusic/grpool.svg?branch=master)](https://github.com/infinitasx/easi-go-aws)

Fiber `Prefork` 模式下的优雅退出

## ⚙ Installation

```text
go get -u github.com/eininst/fiber-prefork-grace
```

## ⚡ Quickstart

```go
package main

import (
    grace "github.com/eininst/fiber-prefork-grace"
    "github.com/gofiber/fiber/v2"
    "time"
)

func main() {
    app := fiber.New(fiber.Config{
        Prefork:     true,
        ReadTimeout: time.Second * 5,
    })

    grace.Listen(app, ":8080")
}
```

> Default listen signal: TERM

```shell
kill -term pid
```

```text
[GRACE] 2022/09/08 06:29:37 Child pid 35951 Shotdown
[GRACE] 2022/09/08 06:29:37 Child pid 35953 Shotdown
[GRACE] 2022/09/08 06:29:37 Child pid 35948 Shotdown
[GRACE] 2022/09/08 06:29:37 Child pid 35946 Shotdown
[GRACE] 2022/09/08 06:29:37 Child pid 35949 Shotdown
[GRACE] 2022/09/08 06:29:37 Child pid 35952 Shotdown
[GRACE] 2022/09/08 06:29:37 Child pid 35950 Shotdown
[GRACE] 2022/09/08 06:29:37 Child pid 35947 Shotdown
[GRACE] 2022/09/08 06:29:37 Main  pid 35945 Shotdown
```

> You can customize it all you want:

```go
grace.Listen(app, ":8080", grace.Config{
    Timeout: time.Second * 15,
    Signal:  syscall.SIGTERM,
})
```


> See [examples](/examples)

## License

*MIT*