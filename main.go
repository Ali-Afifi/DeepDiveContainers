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

	case "child":
		child()

	default:
		fmt.Fprintln(os.Stderr, "wrong command")
		os.Exit(1)
	}
}

func run() {
	fmt.Printf("run %v as %v\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	err := cmd.Run()
	checkError(err)

}

func child() {
	fmt.Printf("run %v as %v\n", os.Args[2:], os.Getpid())

	syscall.Sethostname([]byte("container"))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	checkError(err)

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
