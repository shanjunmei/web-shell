package cmd

import (
	"flag"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/jiangklijna/web-shell/client"
	"github.com/jiangklijna/web-shell/server"
)

// Parameter Command line parameters
type Parameter struct {
	Server      bool
	Client      bool
	HTTPS       bool
	Host        string
	Port        string
	Username    string
	Password    string
	Command     string
	ContentPath string
	CrtFile     string
	KeyFile     string
	RootCrtFile string
}

// Init Parameter
func (parms *Parameter) Init() {
	var (
		help, version bool
	)
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show version and exit")
	flag.BoolVar(&(parms.Server), "s", false, "server mode")
	flag.BoolVar(&(parms.Client), "c", false, "client mode")
	flag.BoolVar(&(parms.HTTPS), "https", false, "enable https")
	flag.StringVar(&(parms.Host), "H", "127.0.0.1", "connect to host")
	flag.StringVar(&(parms.Port), "P", "2020", "listening port")
	flag.StringVar(&(parms.Username), "u", "admin", "username")
	flag.StringVar(&(parms.Password), "p", "webshell", "password")
	flag.StringVar(&(parms.Command), "cmd", "", "command cmd or bash")
	flag.StringVar(&(parms.ContentPath), "cp", "", "content path")
	flag.StringVar(&(parms.CrtFile), "C", "", "crt file")
	flag.StringVar(&(parms.KeyFile), "K", "", "key file")
	flag.StringVar(&(parms.RootCrtFile), "RC", "", "root crt file")

	flag.Parse()
	if help {
		printUsage()
		flag.PrintDefaults()
		os.Exit(1)
	} else if version {
		printVersion()
		os.Exit(1)
	} else {
		parms.organize()
	}
}

// Run start server or client
func (parms *Parameter) Run() {
	if parms.Server {
		s := new(server.WebShellServer)
		s.Init(parms.ContentPath, strings.Split(parms.Command, " ")...)
		s.Run(parms.HTTPS, parms.Port, parms.CrtFile, parms.KeyFile, parms.RootCrtFile)
	} else if parms.Client {
		c := new(client.WebShellClient)
		c.Init(parms.HTTPS, parms.CrtFile, parms.KeyFile, parms.RootCrtFile)
		c.Run(parms.HTTPS, parms.Username, parms.Password, parms.Host, parms.Port, parms.ContentPath)
	}
}

// Organize command line parameters
func (parms *Parameter) organize() {
	if (parms.Server && parms.Client) || (!parms.Server && !parms.Client) {
		println("select server mode OR client mode.")
		os.Exit(1)
	}
	if parms.Server && parms.HTTPS && (parms.CrtFile == "" || parms.KeyFile == "") {
		println("the crt file and key file are required in server mode.")
		os.Exit(1)
	}
	_, err := strconv.Atoi(parms.Port)
	if err != nil {
		println("Port must be an int, not" + parms.Port)
		os.Exit(1)
	}
	parms.Command = strings.Trim(parms.Command, " ")
	if parms.Command == "" {
		parms.Command = defaultCommand()
	}
	parms.ContentPath = strings.Trim(parms.ContentPath, " ")
	if len(parms.ContentPath) > 0 {
		if parms.ContentPath[0] != '/' {
			println("ContentPath must start with /, not", parms.ContentPath)
			os.Exit(1)
		}
		if parms.ContentPath[len(parms.ContentPath)-1] == '/' {
			println("ContentPath connot end with /, not", parms.ContentPath)
			os.Exit(1)
		}
	}
}

func printUsage() {
	println(`Usage:
  web-shell [-s server mode] [-c client mode]  [-P port] [-u username] [-p password] [-cmd command]

Example:
  web-shell -s -P 2020 -u admin -p admin -cmd bash
  web-shell -c -H 192.168.1.1 -P 2020 -u admin -p admin

Options:`)
}

func printVersion() {
	println("web-shell server version:", server.Version)
	println("web-shell client version:", client.Version)
}

func defaultCommand() string {
	if runtime.GOOS == "windows" {
		return "cmd"
	}
	return "bash"
}
