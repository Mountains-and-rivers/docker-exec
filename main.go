package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"github.com/mileusna/crontab"
	"github.com/fsouza/go-dockerclient"
)

func main() {
	//ctx := context.Background()
	ctab := crontab.New()
	err := ctab.AddJob("*/1 * * * *", expireSession)
	if err != nil {
		log.Println(err)
		return
	}
	select {}
}

func expireSession(){
	commad := `
#!/bin/bash
docker ps|grep "k8s_nginx_nginx"|awk '{print $1}'
`
	id,err:=getExecOut(commad)
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	execInContainer(strings.Trim(id,"\n"),"/usr/sbin/nginx -help")
}

func execInContainer(id string,cmd string){
	client, _ := docker.NewClient(`unix:///run/docker.sock`)
	command := []string{"/bin/bash", "-c", cmd}

	exec, err := client.CreateExec(docker.CreateExecOptions{
		AttachStderr: true,
		AttachStdin:  false,
		AttachStdout: true,
		Tty:          false,
		Cmd:          command,
		Container:    id,
	})

	if err != nil {
		log.Println(err)
	}
	var outputBytes []byte
	outputWriter := bytes.NewBuffer(outputBytes)
	var errorBytes []byte
	errorWriter := bytes.NewBuffer(errorBytes)

	err = client.StartExec(exec.ID, docker.StartExecOptions{
		OutputStream: outputWriter,
		ErrorStream:  errorWriter,
	})
	log.Println(outputWriter.String())
	log.Println(errorWriter.String())
	log.Println(err)
}
func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

func getExecOut(commad string) (string,error) {
	cmd := exec.Command("/bin/bash", "-c",commad)
	var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
		return "",err
	}
	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// wg ensures that we finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
		wg.Done()
	}()

	stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return "",err
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
		return "",err
	}
	out, errStr := string(stdout), string(stderr)
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", out, errStr)
	return out,err
}