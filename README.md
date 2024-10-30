# 🖨 Windows printing

Call Windows operating system printer in Golang.

## ✨ Features

See <http://godoc.org/github.com/chenxi2015/winprinters> for details.

- [AddCustomPaperSize](https://pkg.go.dev/github.com/chenxi2015/winprinters#AddCustomPaperSize): add a custom paper specification to the print server;
- [Printer.Forms](https://pkg.go.dev/github.com/chenxi2015/winprinters#Printer.Forms): get all paper size forms on the print server;
- [Printer.Jobs](https://pkg.go.dev/github.com/chenxi2015/winprinters#Printer.Jobs): get all print job information on a printer;
- [ReadNames](https://pkg.go.dev/github.com/chenxi2015/winprinters#ReadNames): get printer names on the system;
- [SetDefault](https://pkg.go.dev/github.com/chenxi2015/winprinters#SetDefault): set default printer for the system;
- [GetDefault](https://pkg.go.dev/github.com/chenxi2015/winprinters#GetDefault): get default printer name on the system;
- ...

## 🔰 Installation

```shell
go get -d github.com/chenxi2015/winprinters
```

📝 Usage

```cgo
package main

import (
    "log"

    "github.com/chenxi2015/winprinters"
)

func main() {
    name, err := winprinters.GetDefault()
    if err != nil {
        log.Fatalln("GetDefault error:", err)
    }

    printer, err := winprinters.Open(name)
    if err != nil {
        log.Fatalln("Open error:", err)
    }
    defer func() {
        _ = printer.Close()
    }()

    jobs, err := printer.Jobs()
    if err != nil {
        log.Fatalln("Jobs error:", err)
    }
    log.Println("jobs:", jobs)
}
```

---

Forked from [godoes/printers](https://github.com/godoes/printers).

## 📃 LICENSE

[BSD-3-Clause](./LICENSE)
