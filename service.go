package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	run(false)
}

func (p *program) Stop(s service.Service) error {
	shutdownSignal <- syscall.SIGINT
	return nil
}

func manageService(cmd string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error when trying to get the working directory:", err.Error())
		os.Exit(1)
	}

	serviceConfig := &service.Config{
		Name:             "lincast",
		DisplayName:      "LinCast",
		Description:      "LinCast, your podcasts player",
		WorkingDirectory: wd,
	}

	p := &program{}

	s, err := service.New(p, serviceConfig)
	if err != nil {
		fmt.Println("Error when trying to create a new instance of Service:", err.Error())
		os.Exit(1)
	}

	switch cmd {
	case "install":
		{
			err := s.Install()
			if err != nil {
				fmt.Println("Error installing the service:", err.Error())
				os.Exit(1)
			}

			fmt.Println("Service installed!")
		}

	case "uninstall":
		{
			err := s.Uninstall()
			if err != nil {
				fmt.Println("Error uninstalling the service:", err.Error())
				os.Exit(1)
			}

			fmt.Println("Service uninstalled!")
		}

	case "start":
		{
			err := s.Start()
			if err != nil {
				fmt.Println("Error starting the service:", err.Error())
				os.Exit(1)
			}

			fmt.Println("Service started!")
		}

	case "stop":
		{
			err := s.Stop()
			if err != nil {
				fmt.Println("Error stopping the service:", err.Error())
				os.Exit(1)
			}

			fmt.Println("Service stopped!")
		}

	case "restart":
		{
			err := s.Restart()
			if err != nil {
				fmt.Println("Error restarting the service:", err.Error())
				os.Exit(1)
			}

			fmt.Println("Service restarted!")
		}

	case "status":
		{
			status, err := s.Status()
			if err != nil {
				fmt.Println("Error getting the status of the service:", err.Error())
				os.Exit(1)
			}

			fmt.Print("Status of the service: ")

			switch status {
			case service.StatusRunning:
				fmt.Println("running")
			case service.StatusStopped:
				fmt.Println("stopped")
			case service.StatusUnknown:
				fmt.Println("unknown")
			}
		}

	default:
		{
			fmt.Printf("Unknown command '%s'\n", cmd)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
