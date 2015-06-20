package main

import (
	"io"
	"net/http"
	"flag"
	"log"
	"io/ioutil"
	"fmt"
	"os"
)

var defaultPath string
var baseURL string


func Handler(w http.ResponseWriter, req *http.Request) {
	
	// fmt.Println("defaultPath in Handler: "+defaultPath)
	// fmt.Println("Base URL",baseURL)
	filename := defaultPath + req.URL.Path[1:]
	// fmt.Println("request: "+filename)
	// Stripping the last /
	if last := len(filename) - 1; last >= 0 && filename[last] == '/' && len(filename)!=1 {
        filename = filename[:last]
    }
	
	if(filename==""){
		// fmt.Println("request is empty ")
		filename = "./"
	}

	// fmt.Println(reflect.TypeOf(filename))
	
	// If file Exists
	if file, err := os.Stat(filename); os.IsNotExist(err) {
    	
    	// http.Redirect(w, req, "", http.StatusNotFound)
    	_,err = io.WriteString(w, "404 Not Found")

    	return
	
	} else{
	// fmt.Println(file,err)
	// fmt.Println(file.IsDir())
	// fmt.Printf("%v",file)
		

		if file.IsDir() {
			
			slashCheck := ""

			files, _ := ioutil.ReadDir(filename)
			if filename != "./"{
				if filename[len(filename)-1] !='/'{
					slashCheck = "/"
					// fmt.Println("Adding / ",filename)
				}
			}

			responseString := "<html><body> <h3> Directory Listing for "+req.URL.Path[1:] +"/ </h3> <br/> <hr> "
			// fmt.Println("\nDirectory Contents \n")
	    	for _, f := range files {
	            // fmt.Println(f.Name())
	            if(f.Name()[0]!= '.') {
	            	if(f.IsDir()){
			            newLink := "<a href=\"" + baseURL + req.URL.Path[0:] + slashCheck + f.Name() + "\">" + f.Name() + "/" + "</a><br/>"
			            responseString = responseString + newLink
	            	} else {
			            newLink := "<a href=\"" + baseURL + req.URL.Path[0:] + slashCheck + f.Name() + "\">" + f.Name() + "</a><br/>"
			            responseString = responseString + newLink
	            	}
	            }
	    	}
	    	responseString = responseString + "</body></html>"
	    	_,err = io.WriteString(w, responseString)
	    	if err != nil {
		        // panic(err)
				http.Redirect(w, req, "", http.StatusInternalServerError)
		    }
		}else{
			
			
			b, err := ioutil.ReadFile(filename)
		    if err != nil {
				http.Redirect(w, req, "", http.StatusInternalServerError)
		    	return
		    }else{
			    str := string(b)
				_,err = io.WriteString(w, str)
			    if err != nil {
			        // panic(err)
					http.Redirect(w, req, "", http.StatusInternalServerError)
			    }
		    }
		}
	

	}
	// http.Redirect(w, req, "", http.StatusNotFound)
	// io.WriteString(w, "hello, world!\n")
}

func main() {

	
	defaultPortPtr := flag.String("p","","Port Number")
	defaultPathPtr := flag.String("dir","","Root Directory")
	flag.Parse()

	if *defaultPathPtr !=""{
		defaultPath = "./" + *defaultPathPtr + "/"
	} else {
		defaultPath = ""
	}

	portNum := "8080"
	if *defaultPortPtr !=""{
		portNum = *defaultPortPtr
	}

	baseURL = "http://localhost:" + portNum;

	fmt.Println("Serving on ",baseURL)
	// fmt.Println("Setting Default Path: ",defaultPath)
    // fmt.Println(os.Args[1])
	http.HandleFunc("/", Handler)
	err := http.ListenAndServe(":"+portNum, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
