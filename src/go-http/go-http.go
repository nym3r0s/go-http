package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net"
	"os"
)

var defaultPath string
var baseURL string

func Handler(w http.ResponseWriter, req *http.Request) {

	filename := defaultPath + req.URL.Path[1:]
	if last := len(filename) - 1; last >= 0 && filename[last] == '/' && len(filename) != 1 {
		filename = filename[:last]
	}

	if filename == "" {
		// fmt.Println("request is empty ")
		filename = "./"
	}

	// If file Exists
	if file, err := os.Stat(filename); os.IsNotExist(err) {

		_, err = io.WriteString(w, "404 Not Found")

		return

	} else {

		if file.IsDir() {

			slashCheck := ""

			files, _ := ioutil.ReadDir(filename)
			if filename != "./" {
				if filename[len(filename)-1] != '/' {
					slashCheck = "/"
				}
			}

			responseString := "<html><body> <h3> Directory Listing for " + req.URL.Path[1:] + "/ </h3> <br/> <hr> "
			for _, f := range files {
				if f.Name()[0] != '.' {
					if f.IsDir() {
						newLink := "<a href=\"" + baseURL + req.URL.Path[0:] + slashCheck + f.Name() + "\">" + f.Name() + "/" + "</a><br/>"
						responseString = responseString + newLink
					} else {
						newLink := "<a href=\"" + baseURL + req.URL.Path[0:] + slashCheck + f.Name() + "\">" + f.Name() + "</a><br/>"
						responseString = responseString + newLink
					}
				}
			}
			responseString = responseString + "</body></html>"
			_, err = io.WriteString(w, responseString)
			if err != nil {
				// panic(err)
				http.Redirect(w, req, "", http.StatusInternalServerError)
			}
		} else {

			b, err := ioutil.ReadFile(filename)
			if err != nil {
				http.Redirect(w, req, "", http.StatusInternalServerError)
				return
			} else {
				str := string(b)
				_, err = io.WriteString(w, str)
				if err != nil {
					// panic(err)
					http.Redirect(w, req, "", http.StatusInternalServerError)
				}
			}
		}

	}
}

func main() {

	defaultHostPtr := flag.String("h", "", "Host address")
	defaultPortPtr := flag.String("p", "", "Port Number")
	defaultPathPtr := flag.String("d", "", "Root Directory")
	flag.Parse()

	if *defaultPathPtr != "" {
		defaultPath = "./" + *defaultPathPtr + "/"
	} else {
		defaultPath = ""
	}

	host := "localhost"
	if *defaultHostPtr != "" {
		addrs, err := net.LookupHost(*defaultHostPtr)

		if err != nil {
			fmt.Println("Impossible to resolv domain name:", err)

			ipaddr := net.ParseIP(*defaultHostPtr)

			if ipaddr != nil {
				host = *defaultHostPtr
			} else {
				log.Fatal("No address IP could be decoded")
			}
		} else {
			host = addrs[0]
		}
	}

	portNum := "8080"
	if *defaultPortPtr != "" {
		portNum = *defaultPortPtr
	}

	baseURL = "http://" + host + ":" + portNum

	fmt.Println("Serving on ", baseURL, " subdirectory ", defaultPath)

	http.HandleFunc("/", Handler)
	err := http.ListenAndServe(":"+portNum, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
