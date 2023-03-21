package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "run <cmd> <params>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		run()

	case "fork":

		if os.Args[0] != "/proc/self/exe" {
			fmt.Fprintln(os.Stderr, "wrong command")
			os.Exit(1)
		}

		fork()

	default:
		fmt.Fprintln(os.Stderr, "wrong command")
		os.Exit(1)
	}
}

func run() {

	cmd := exec.Command("/proc/self/exe", append([]string{"fork"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,

		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},

		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Start(); err != nil {
		checkError(err)
	}

	fmt.Printf("run %v as %v\n", os.Args[0], cmd.Process.Pid)

	err := cmd.Wait()
	checkError(err)

}

func fork() {
	fmt.Printf("\n>> namespace setup code goes here <<\n\n")

	cwd, _ := os.Getwd()

	if err := mountProc(cwd + "/tmp-rootfs"); err != nil {
		fmt.Printf("Error mounting /proc - %s\n", err)
		os.Exit(1)
	}

	if err := pivotRoot(cwd + "/tmp-rootfs"); err != nil {
		fmt.Printf("Error running pivot_root - %s\n", err)
		os.Exit(1)
	}

	if err := syscall.Sethostname([]byte("container")); err != nil {
		fmt.Printf("Error setting hostname - %s\n", err)
		os.Exit(1)
	}


	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	cmd.Env = append(cmd.Env, "TERM="+os.Getenv("TERM"), "PS1=\\u@[\\h]--[\\w] # ")

	if err := cmd.Start(); err != nil {
		checkError(err)
	}

	fmt.Printf("run %v as %v\n", os.Args[2:], cmd.Process.Pid)

	err := cmd.Wait()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
