package main

import (
	"bytes"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/cihub/seelog"
)

func RunCommand(tDuration time.Duration, commands string) (string, string) {
	log.Debug(commands)
	command := strings.Fields(commands)
	cmd := exec.Command(command[0])
	if len(command) > 0 {
		cmd = exec.Command(command[0], command[1:]...)
	}
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Start()
	if err != nil {
		log.Error(err)
		log.Flush()
		os.Exit(1)
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(tDuration):
		if err := cmd.Process.Kill(); err != nil {
			log.Debug("failed to kill: ", err)
		}
		log.Debugf("%s killed as timeout reached", commands)
	case err := <-done:
		if err != nil {
			log.Debugf("err running %s: %s", commands, err.Error())
		} else {
			log.Debugf("%s done gracefully without error", commands)
		}
	}
	return strings.TrimSpace(outb.String()), strings.TrimSpace(errb.String())
}

func Average(nums []float64) float64 {
	total := float64(0)
	for _, num := range nums {
		total += num
	}
	return float64(int(total/float64(len(nums))*10)) / 10
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func RandomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}


func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
