package main

import (
    "lincast/webui/backend"
)

func main() {
    if r := run(); r != nil {
        panic("Error on run: " + r.Error())
    }
}

func run() error {
    // Arguments should be obtained from settings.
    err := backend.Run(8080, true, false, true)
    if err != nil {
        return err
    }

    return nil
}
