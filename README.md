# 🖨 Windows printing

Call Windows operating system printer in Golang.

## ✨ Features

See <http://godoc.org/github.com/chenxi2015/printers> for details.

- [AddCustomPaperSize](https://pkg.go.dev/github.com/godoes/printers#AddCustomPaperSize): add a custom paper specification to the print server;
- [Printer.Forms](https://pkg.go.dev/github.com/godoes/printers#Printer.Forms): get all paper size forms on the print server;
- [Printer.Jobs](https://pkg.go.dev/github.com/godoes/printers#Printer.Jobs): get all print job information on a printer;
- [ReadNames](https://pkg.go.dev/github.com/godoes/printers#ReadNames): get printer names on the system;
- [SetDefault](https://pkg.go.dev/github.com/godoes/printers#SetDefault): set default printer for the system;
- [GetDefault](https://pkg.go.dev/github.com/godoes/printers#GetDefault): get default printer name on the system;
- ...

## 🔰 Installation

```shell
go get -d github.com/chenxi2015/printers
```

📝 Usage

```cgo
package main

import (
    "log"

    "github.com/chenxi2015/printers"
)

func main() {
    name, err := printers.GetDefault()
    if err != nil {
        log.Fatalln("GetDefault error:", err)
    }

    printer, err := printers.Open(name)
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

Forked from [alexbrainman/printer](https://github.com/alexbrainman/printer).

## 📃 LICENSE

[BSD-3-Clause](./LICENSE)
