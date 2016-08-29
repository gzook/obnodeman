package nodeman

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type Client struct {
	mutex  *sync.Mutex
	runCmd *exec.Cmd
}

func New() *Client {
	return &Client{mutex: &sync.Mutex{}}
}

func (client *Client) Running() bool {
	return client.runCmd != nil
}

func (client *Client) Stop() (err error) {

	client.mutex.Lock()
	defer client.mutex.Unlock()

	if client.runCmd != nil {
		if err = client.runCmd.Process.Signal(syscall.SIGINT); err != nil {
			err = fmt.Errorf("nodeman: failed to send sigint. %s", err.Error())
		}
		waitOnStop(client.runCmd)
	} else {
		err = errors.New("nodeman: cannot stop as not running")
	}

	if err != nil {
		log.Println(err.Error())
	}
	client.runCmd = nil
	return err
}

func (client *Client) Start() (err error) {

	client.mutex.Lock()
	defer client.mutex.Unlock()

	if client.runCmd != nil {
		return errors.New("obnodeman: cannot start as already running")
	}
	if err = checkWorkingDir(); err != nil {
		log.Println(err.Error())
		return err
	}

	args := []string{
		"openbazaard.py",
		"start",
		"-a",
		"0.0.0.0",
	}

	if err = client.startCmd(exec.Command("python", args...)); err != nil {
		err = fmt.Errorf("nodeman: failed to start. %s", err.Error())
	}

	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func waitOnStop(cmd *exec.Cmd) {
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(10 * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("nodeman: attempted, but failed to kill ob process as it did not shut down within the alloted time. Error: %s", err.Error())
		} else {
			log.Println("nodeman: killed ob process as it did not shut down within the alloted time")
		}
	case err := <-done:
		if err != nil {
			log.Printf("nodeman: ob process exited with error = %v", err)
		} else {
			log.Println("nodeman: ob process exited gracefully")
		}
	}
}
func checkWorkingDir() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("nodeman: failed to read working dir: %s", err.Error())
	}
	filename := "openbazaard.py"
	fullname := filepath.Join(workingDir, filename)
	if _, err := os.Stat(fullname); os.IsNotExist(err) {
		return fmt.Errorf("nodeman: %s not found in working dir %s", filename, workingDir)
	}
	return nil
}

func (client *Client) startCmd(cmd *exec.Cmd) error {

	if err := cmd.Start(); err != nil {
		return err
	}
	client.runCmd = cmd
	return nil
}
