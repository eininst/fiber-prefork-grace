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
    grace "github.com/eininst/fiber-grace"
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
2022/09/08 06:21:01 [GRACE] Received signal: terminated
2022/09/08 06:21:01 [GRACE] Child pid 35210 Shotdown
2022/09/08 06:21:01 [GRACE] Child pid 35212 Shotdown
2022/09/08 06:21:01 [GRACE] Child pid 35213 Shotdown
2022/09/08 06:21:01 [GRACE] Child pid 35214 Shotdown
2022/09/08 06:21:01 [GRACE] Child pid 35208 Shotdown
2022/09/08 06:21:01 [GRACE] Child pid 35207 Shotdown
2022/09/08 06:21:01 [GRACE] Child pid 35209 Shotdown
2022/09/08 06:21:01 [GRACE] Child pid 35211 Shotdown
2022/09/08 06:21:01 [GRACE] Main  pid 35204 Shotdown
```

> You can customize it all you want:

```go
grace.Listen(app, ":8080", grace.Config{
    Timeout: time.Second * 15,
    Sig:     syscall.SIGILL,
})
```