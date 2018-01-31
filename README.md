# semaphore package 

## History
* Jan 31st 2018 Release 0.0.2 for the first version of windows and linux


## example
```go
package main

import (
    "fmt"
    "github.com/jeppeter/go-semaphore"
    "os"
    "os/signal"
    "strconv"
)

func Usage(ec int, fmtstr string, a ...interface{}) {
    fp := os.Stderr
    if ec == 0 {
        fp = os.Stdout
    }

    if fmtstr != "" {
        fmt.Fprintf(fp, fmtstr, a...)
        fmt.Fprintf(fp, "\n")
    }

    fmt.Fprintf(fp, "%s [commands]\n", os.Args[0])
    fmt.Fprintf(fp, "\tcreate name cnt              to create semaphore\n")
    fmt.Fprintf(fp, "\topen name                    to open semaphore\n")
    fmt.Fprintf(fp, "\trelease name                 to release semaphore\n")
    fmt.Fprintf(fp, "\twait name [times]            to wait semaphore\n")
    os.Exit(ec)
}

func main() {
    var p *semaphore.Semaphore
    var err error
    var cnt int = 1
    var mills int = -1
    sigch := make(chan os.Signal, 1)
    signal.Notify(sigch, os.Interrupt)
    if len(os.Args) < 3 {
        Usage(3, "")
    }
    if len(os.Args) > 3 {
        cnt, _ = strconv.Atoi(os.Args[3])
    }
    if len(os.Args) > 4 {
        mills, _ = strconv.Atoi(os.Args[4])
    }

    if os.Args[1] == "create" {
        p, err = semaphore.NewSemaphore(os.Args[2], cnt)
    } else if os.Args[1] == "open" {
        p, err = semaphore.NewSemaphore(os.Args[2], cnt)
    } else if os.Args[1] == "wait" {
        p, err = semaphore.NewSemaphore(os.Args[2], cnt)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err.Error())
            os.Exit(5)
        }
        err = p.Wait(mills)
    } else if os.Args[1] == "release" {
        p, err = semaphore.NewSemaphore(os.Args[2], cnt)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err.Error())
            os.Exit(5)
        }
        err = p.Release()
    } else {
        Usage(3, "not support cmd [%s]", os.Args[1])
    }

    if err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err.Error())
        os.Exit(5)
    }
    defer p.Close()

    fmt.Fprintf(os.Stdout, "%s %s succ\n", os.Args[1], os.Args[2])

    for {
        select {
        case <-sigch:
            goto out
        }
    }
out:
    return
}
```

## to run the command 
```shell
go build -o simple simple.go
```

## run in three console
```shell
./simple wait hello_go 1
wait hello_go succ
```

```shell
./simple wait hello_go 1
```
> the second will not output "wait hello_go succ"
> if you put third console
```shell
./simple release hello_go 1
```
> the second console will output "wait hello_go succ"
> and third console will output "release hello_go succ"
