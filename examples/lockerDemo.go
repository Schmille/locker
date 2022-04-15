package main

import (
    "locker"
    "time"
)

func main() {
    locker.Info("Welcome to the locker demo!")
    locker.Debug("I hope you enjoy your stay ")
    go func() {
        locker.Info("Lets start a side project while we do something else")
        time.Sleep(50 * time.Microsecond)
        locker.Info("This should be good enough")
    }()

    locker.Warn("About to start a long and tedious task")
    locker.Push("LongAndTediousTask")
    locker.Debug("That is a long task")
    locker.Debug("This really is long")
    go func() {
        locker.Info("I'd like to interject here and say that it's actually go not golang, you see...")
    }()
    locker.Error("This takes a bit too long, let's stop doing this")
    locker.Pop()
    time.Sleep(5 * time.Second)
}

