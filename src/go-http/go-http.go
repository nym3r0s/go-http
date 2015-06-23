package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var defaultPath string

var baseURL string

func Logger(req *http.Request, statusCode int) {
	const layout = "[ 2/Jan/2006 15:04:05 ]"
	fmt.Println(baseURL + " --- " + time.Now().Format(layout) + " " + req.Method + "  " + strconv.Itoa(statusCode) + "  " + req.URL.Path)
}

func Handler(w http.ResponseWriter, req *http.Request) {

	filename := defaultPath + req.URL.Path[1:]
	if last := len(filename) - 1; last >= 0 && filename[last] == '/' && len(filename) != 1 {
		filename = filename[:last]
	}

	// Empty request (Root)
	if filename == "" {
		filename = "./"
	}

	file, err := os.Stat(filename)

	// 404 if file doesn't exist
	if os.IsNotExist(err) {
		_, err = io.WriteString(w, "404 Not Found")
		Logger(req, http.StatusNotFound)
		return
	}

	// Serve directory
	if file.IsDir() {

		slashCheck := ""

		files, err := ioutil.ReadDir(filename)
		// Catch the Error in reading from directory
		if err != nil {
			http.Redirect(w, req, "", http.StatusInternalServerError)
			Logger(req, http.StatusInternalServerError)
		}
		// Checking for Root Directory
		if filename != "./" {
			if filename[len(filename)-1] != '/' {
				slashCheck = "/"
			}
		}

		responseString := "<html><body> <h3> Directory Listing for " + req.URL.Path[1:] + "/ </h3> <br/> <hr> <ul>"
		for _, f := range files {
			if f.Name()[0] != '.' {
				if f.IsDir() {
					responseString += "<li><a href=\"" + baseURL + req.URL.Path[0:] + slashCheck + f.Name() + "\">" + f.Name() + "/" + "</a></li>"
				} else {
					responseString += "<li><a href=\"" + baseURL + req.URL.Path[0:] + slashCheck + f.Name() + "\">" + f.Name() + "</a></li>"
				}
			}
		}

		//Ending the list
		responseString += "</ul><br/><hr/>"

		p := req.URL.Path

		// Display link to parent directory
		if len(p) > 1 {
			base := path.Base(p)

			slice := len(p) - len(base) - 1

			url := "/"

			if slice > 1 {
				url = req.URL.Path[:slice]
				url = strings.TrimRight(url, "/") // Remove extra / at the end
			}

			responseString += "<br/><a href=\"" + baseURL + url + "\">Parent directory</a>"
		}

		responseString = responseString + "</body></html>"
		_, err = io.WriteString(w, responseString)
		if err != nil {
			// panic(err)
			http.Redirect(w, req, "", http.StatusInternalServerError)
			Logger(req, http.StatusInternalServerError)
		} else {
			Logger(req, http.StatusOK)
		}

		return
	}

	// File exists and is no directory; Serve the file

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		http.Redirect(w, req, "", http.StatusInternalServerError)
		Logger(req, http.StatusInternalServerError)
		return
	}

	str := string(b)
	_, err = io.WriteString(w, str)
	if err != nil {
		// panic(err)
		http.Redirect(w, req, "", http.StatusInternalServerError)
	} else {
		Logger(req, http.StatusOK)
	}

}

func main() {

	defaultPortPtr := flag.String("p", "", "Port Number")
	defaultPathPtr := flag.String("d", "", "Root Directory")
	flag.Parse()

	portNum := "8080"

	// Handling the command line flags

	// Directory
	if *defaultPathPtr != "" {
		defaultPath = "./" + *defaultPathPtr + "/"
	} else {
		defaultPath = ""
	}
	// Port Number
	if *defaultPortPtr != "" {
		portNum = *defaultPortPtr
	} else {
		portNum = "8080"
	}

	baseURL = "http://localhost:" + portNum

	fmt.Println("Serving on ", baseURL, " subdirectory ", defaultPath)

	http.HandleFunc("/", Handler)
	err := http.ListenAndServe(":"+portNum, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
