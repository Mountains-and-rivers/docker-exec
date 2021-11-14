package main

import (
	"fmt"
	"log"
	"github.com/mileusna/crontab"
	"github.com/ahmetalpbalkan/dexec"
	"github.com/fsouza/go-dockerclient"
)

func test(){
	cl, _ := docker.NewClient("unix:///run/docker.sock")
	d := dexec.Docker{cl}

	m, _ := dexec.ByCreatingContainer(docker.CreateContainerOptions{
		Config: &docker.Config{Image: "8c811b4aec35"}})

	cmd := d.Command(m, "ifconfig", `eth0`)
	b, err := cmd.Output()
	if err != nil { log.Fatal(err) }
	log.Printf("%s", b)
}
func main() {
	ctab := crontab.New()
	err := ctab.AddJob("*/1 * * * *", myFunc)
	if err != nil {
		log.Println(err)
		return
	}
	select {

	}
}

func myFunc() {
	fmt.Println("Helo, world")
}

func myFunc3() {
	fmt.Println("Noon!")
}

func myFunc2(s string, n int) {
	fmt.Println("We have params here, string", s, "and number", n)
}
