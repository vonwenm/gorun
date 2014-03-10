package main

import (
    "flag"
    "fmt"
)

var (
    FilePattern = `(.+\.go|.+\.c)$`
)

func main() {
    flag.String("p", FilePattern, "Pattern of watched files")
    flag.String("command", "", "Command to run and restart after build")

    fmt.Println(1)
    fmt.Println(1)
}
