package main

import (
	"log"

	"github.com/ahmetalpbalkan/dexec"
	"github.com/fsouza/go-dockerclient"
)

func main(){
	cl, _ := docker.NewClient("unix:///run/docker.sock")
	d := dexec.Docker{cl}

	m, _ := dexec.ByCreatingContainer(docker.CreateContainerOptions{
		Config: &docker.Config{Image: "8c811b4aec35"}})

	cmd := d.Command(m, "ifconfig", `eth0`)
	b, err := cmd.Output()
	if err != nil { log.Fatal(err) }
	log.Printf("%s", b)
}
