package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/webview/webview"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Action struct {
	Name    string `json:"name"`
	Command string `json:"command,omitempty"`
}

//go:embed index.html
var index []byte

var commands []Action
var DEV bool

func main() {

	DEV = os.Getenv("DEV") == "1"

	commands = loadCommands()
	r := gin.Default()
	r.GET("/actions", handleListActions)
	r.GET("/", handleHome)

	r.GET("/task", handleTask)

	if DEV {
		err := r.Run(":8144")
		if err != nil {
			log.Fatalln("could not start app")
		}

	} else {
		go func() {
			err := r.Run(":8144")
			if err != nil {
				log.Fatalln("could not start app")
			}
		}()
		w := webview.New(true)
		defer w.Destroy()
		w.SetTitle("plip")
		w.SetSize(800, 600, webview.HintNone)
		w.Navigate("http://localhost:8144")
		w.Run()
	}
}

func handleTask(c *gin.Context) {
	action := c.Query("action")
	if action == "quit" {
		c.JSON(200, gin.H{})
		os.Exit(0)
	}
	for _, item := range commands {
		if item.Name == action {
			args := strings.Fields(item.Command)
			cmd := exec.Command(args[0], args[1:]...)
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
	if !DEV {
		os.Exit(0)
	}
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
