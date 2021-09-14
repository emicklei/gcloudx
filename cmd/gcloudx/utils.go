package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func promptForYes(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	yn, _ := reader.ReadString('\n')
	return strings.HasPrefix(yn, "Y") || strings.HasPrefix(yn, "y")
}

func logBegin(c *cli.Context) func() {
	if !c.Bool("v") {
		return func() {}
	}
	buf := new(bytes.Buffer)
	fmt.Fprint(buf, "gcloudx")

	appendFlag := func(each cli.Flag) {
		fv := reflect.ValueOf(each)
		hide := reflect.Indirect(fv).FieldByName("Hidden").Bool()
		name := each.Names()[0]
		value := c.Generic(name)
		if hide {
			value = "**hidden**"
		}
		fmt.Fprintf(buf, " %s=%v", name, value)
	}
	action := c.Command.FullName()
	fmt.Fprintf(buf, " %s (", action)
	for _, each := range c.App.Flags {
		appendFlag(each)
	}
	for _, each := range c.Command.Flags {
		appendFlag(each)
	}
	fmt.Fprint(buf, " )")
	log.Println(buf.String())
	start := time.Now()
	return func() {
		log.Println(action, "completed in", time.Since(start))
		if err := recover(); err != nil {
			// no way to communicate error to cli so exit here.
			log.Fatalln(action, "FAILED with error:", err)
		}
	}
}
