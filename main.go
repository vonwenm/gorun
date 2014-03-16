package main

import (
    "github.com/howeyc/fsnotify"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
)

var (
    runningApp     *exec.Cmd
    defaultPath, _ = filepath.Abs("./")
    appName        = filepath.Base(defaultPath)
    appArgs        = strings.Join(os.Args[1:], " ")
)

func main() {
    // Get subfolder path
    paths, err := Walk(defaultPath)
    if err != nil {
        log.Fatalln(err)
    }

    // Watch and run
    Build()
    go Start()
    Watch(paths)
}

func Walk(rootDir string) (paths []string, err error) {
    err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() || strings.Contains(path, ".git") {
            return nil
        }
        paths = append(paths, path)
        return nil
    })
    if err != nil {
        return
    }
    return
}

func Watch(paths []string) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatalln(err)
    }

    done := make(chan bool)

    go func() {
        var prevActionSecond int
        for {
            select {
            case ev := <-watcher.Event:
                if filepath.Ext(ev.Name) == ".go" {
                    // Prevent the same action output many times.
                    if prevActionSecond-time.Now().Second() >= -1 {
                        continue
                    }
                    // Must be put after ignoring file extension checking, because arise bug if first .fff.swp second fff
                    prevActionSecond = time.Now().Second()
                    //log.Println("Rebuild")
                    Rebuild()
                }
            case err := <-watcher.Error:
                log.Println("error:", err)
            }
        }
    }()

    for _, path := range paths {
        err = watcher.Watch(path)
        if err != nil {
            log.Fatalln(err)
        }
    }

    //log.Println("Begin to watch app:", appName)
    <-done
    watcher.Close()
}

func Build() (err error) {
    begin := time.Now().UnixNano()
    cmd := exec.Command("go", "build")

    // Let standard output and error to show on the screen
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Wait for build
    if err = cmd.Run(); err != nil {
        log.Println("cmd.Run error : " + err.Error())
    }
    log.Println("Build passed:", (time.Now().UnixNano()-begin)/1000/1000, "ms")
    return
}

func Rebuild() {
    err := Build()
    if err != nil {
        log.Println("Rebuild fail : " + err.Error())
    } else {
        ReStart()
    }
}

func ReStart() {
    if runningApp != nil {
        //log.Println("Kill old running app:", appName)
        runningApp.Process.Kill()
    }
    Start()
}

func Start() {
    //log.Println("./" + appName + appArgs)
    runningApp = exec.Command("./"+appName, appArgs)
    runningApp.Stdout = os.Stdout
    runningApp.Stderr = os.Stderr
    //log.Println("Start running app:", appName)
    go runningApp.Run()
}
