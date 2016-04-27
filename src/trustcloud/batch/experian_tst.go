package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("Experian test...")

	out, err := exec.Command("sh", "-c", "./run.sh").Output()

	fmt.Println("out... ", string(out))

	fmt.Println("err... ", err)

	//	out, err := exec.Command("sh", "-c", cmd).Output()

}
