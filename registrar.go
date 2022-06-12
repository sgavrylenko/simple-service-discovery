package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"time"
)

type Registrar struct {
	Interval  time.Duration
	DockerCLI *client.Client
	SRegistry *ServiceRegistry
}

const (
	HelloServiceImageName = "hello"
	ContainerRunningState = "running"
)

func (r *Registrar) Observe() {
	for range time.Tick(r.Interval) {
		cList, _ := r.DockerCLI.ContainerList(context.Background(), types.ContainerListOptions{
			All: true,
		})

		if len(cList) == 0 {
			r.SRegistry.RemoveAll()
			continue
		}

		for _, c := range cList {
			if c.Image != HelloServiceImageName {
				continue
			}

			_, exist := r.SRegistry.GetByContainerID(c.ID)

			if c.State == ContainerRunningState {
				//fmt.Printf("found container %s with port %d\n", c.ID, c.Ports[0].)
				if !exist {
					var containerPort uint16
					var containerAddr string

					for networkContainer, _ := range c.NetworkSettings.Networks {
						if c.Ports[0].PublicPort != 0 {
							containerPort = c.Ports[0].PublicPort
							containerAddr = "localhost"
						} else {
							containerPort = c.Ports[0].PrivatePort
							containerAddr = c.NetworkSettings.Networks[networkContainer].IPAddress
						}

						addr := fmt.Sprintf("http://%s:%d", containerAddr, containerPort)
						r.SRegistry.Add(c.ID, addr)
						fmt.Printf("Add container %s from network %s with address %s:%d\n", c.ID, networkContainer, containerAddr, containerPort)
					}
				}
			} else {
				if exist {
					r.SRegistry.RemoveByContainerID(c.ID)
				}
			}
		}
	}
}
