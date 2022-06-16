package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type Action struct {
	Name     string   `json:"name"`
	Command  string   `json:"command"`
	Commands []Action `json:"commands"`
}

//go:embed index.html
var index []byte

var commands []Action

func main() {

	commands = loadCommands()
	r := gin.Default()
	r.GET("/actions", handleListActions)
	r.GET("/", handleHome)

	r.GET("/trigger", handleTask)

	err := r.Run(":8144")
	if err != nil {
		log.Fatalln("could not start app")
	}

}

func handleTask(c *gin.Context) {
	action := c.Query("action")
	position := c.Query("position")
	if action == "quit" {
		c.JSON(200, gin.H{})
		os.Exit(0)
	}

	var selectedCommands []Action
	if position != "undefined" {
		x, err := strconv.Atoi(position)
		if err != nil {
			c.JSON(400, gin.H{"message": "wrong item position selection"})
			return
		}
		selectedCommands = commands[x].Commands
	} else {
		selectedCommands = commands
	}

	for _, item := range selectedCommands {
		if item.Name == action {
			fmt.Println("do command", "./commands/"+item.Command)
			cmd := exec.Command("./commands/" + item.Command)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(out.String())
		}
	}
	c.JSON(200, gin.H{})
}

func handleHome(c *gin.Context) {
	c.Data(http.StatusOK, "text/html", index)
}

func handleListActions(c *gin.Context) {
	d, err := json.Marshal(commands)
	if err != nil {
		log.Fatalln("Error serializing commands")
	}
	c.Data(http.StatusOK, "application/json", d)
}

func loadCommands() []Action {
	var confPath string
	if os.Getenv("ACTIONS_PATH") == "" {
		confPath = ""
	} else {
		confPath = os.Getenv("ACTIONS_PATH")
	}
	var commands []Action
	d, err := ioutil.ReadFile(confPath + ".actions.json")
	if err != nil {
		log.Fatalln("Error reading actions file")
	}
	err = json.Unmarshal(d, &commands)
	if err != nil {
		log.Fatalln("Error parsing actions file")
	}
	return commands
}
