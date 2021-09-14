package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
)

func started(c *cli.Context, action string) func() {
	v := c.GlobalBool("v")
	if !v {
		return func() {}
	}
	log.Println("gmig version", Version)
	log.Println("BEGIN", action)
	start := time.Now()
	return func() { log.Println("END", action, "completed in", time.Now().Sub(start)) }
}

func promptForYes(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	yn, _ := reader.ReadString('\n')
	return strings.HasPrefix(yn, "Y") || strings.HasPrefix(yn, "y")
}
