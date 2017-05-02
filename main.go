package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/TexaProject/texalib"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/welcome", 301)
}

func texaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get	request	method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("www/index.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// fmt.Printf("%+v\n", r.Form)

		QSA := r.Form.Get("scoreArray")
		fmt.Println(QSA)

		re := regexp.MustCompile("[0-1]+")

		array := re.FindAllString(QSA, -1)

		fmt.Println("Resulting Array:")
		for x := range array {
			fmt.Println(array[x])
		}

		fmt.Println("Converted Array:")
		convArray := texalib.Convert(array)
		fmt.Println(texalib.Convert(array))

		fmt.Println("Total of Converted Array")
		fmt.Println(texalib.Total(convArray))
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get	request	method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("www/welcome.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
	}
}

// upload logic for JS dictionary/bot data file
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("login.html")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprint(w, "ACKNOWLEDGEMENT:\nUploaded the file. Header Info:\n")
		fmt.Fprintf(w, "%v", handler.Header)
		fmt.Fprint(w, "\n\nVISIT: /texa for interrogation.")
		f, err := os.OpenFile("./www/js/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Selected file: ", handler.Filename)
		defer f.Close()
		io.Copy(f, file)
		// http.Redirect(w, r, "/texa", 301)
	}
}

func resultHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	fmt.Println("--TEXA SERVER--")
	fmt.Println("STATUS: INITIATED")
	fmt.Println("ADDR: http://127.0.0.1:3030")

	fs := http.FileServer(http.Dir("www/js"))
	http.Handle("/js/", http.StripPrefix("/js/", fs))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/welcome", welcomeHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/texa", texaHandler)
	http.HandleFunc("/result", resultHandler)

	http.ListenAndServe(":3030", nil)
}
