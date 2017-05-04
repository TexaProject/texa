package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
	//Import this by exec in CLI: `go get -u github.com/TexaProject/texalib`
	"github.com/TexaProject/texalib"
)

type Page struct {
	AIName   string  `json:"AIName"`
	IntName  string  `json:"IntName"`
	ArtiMts  float64 `json:"ArtiMts"`
	HumanMts float64 `json:"HumanMts"`
}

func (p Page) toString() string {
	return toJson(p)
}

//getPages() returns a converted Page Array persistent to the mts.json
func getPages() []Page {
	raw, err := ioutil.ReadFile("./www/data/mts.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []Page
	json.Unmarshal(raw, &c)
	return c
}

// convtoPage converts a set of data vars into a Page struct variable
func convtoPage(AIName string, IntName string, ArtiMts float64, HumanMts float64) Page {
	var newPage Page
	newPage.AIName = AIName
	newPage.IntName = IntName
	newPage.ArtiMts = ArtiMts
	newPage.HumanMts = HumanMts
	return newPage
}

// AddtoPageArray() Appends a new page 'p' to the specified target PageArray 'pa'
func AddtoPageArray(p Page, pa []Page) []Page {
	for x := range pa {
		if p == pa[x] {
			panic("JSON ERROR: Can't append a Duplicate Page into PageArray")
		}
	}
	return (append(pa, p))
}

// toJSON Marshals PageArray data into JSON format
func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	ioutil.WriteFile("./www/data/mts.json", bytes, 0644)
	return string(bytes)
}

// AIName exports form value from /welcome
var AIName string

// IntName exports form value from /texa
var IntName string

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

		fmt.Println("--INTERROGATION FORM DATA--")
		IntName = r.Form.Get("IntName")

		fmt.Println("###", AIName)
		fmt.Println("###", IntName)

		QSA := r.Form.Get("scoreArray")
		fmt.Println(QSA)

		re := regexp.MustCompile("[0-1]+")

		array := re.FindAllString(QSA, -1)

		fmt.Println("Resulting Array:")
		for x := range array {
			fmt.Println(array[x])
		}

		ArtiQSA := texalib.Convert(array)
		fmt.Println("ArtiQSA:")
		fmt.Println(ArtiQSA)

		HumanQSA := texalib.SetHumanQSA(ArtiQSA)
		fmt.Println("HumanQSA:")
		fmt.Println(HumanQSA)

		TSA := texalib.GetTransactionSeries(ArtiQSA, HumanQSA)
		fmt.Println("TSA:")
		fmt.Println(TSA)

		ArtiMts := texalib.GetMeanTestScore(ArtiQSA)
		HumanMts := texalib.GetMeanTestScore(HumanQSA)

		fmt.Println("ArtiMts: ", ArtiMts)
		fmt.Println("HumanMts: ", HumanMts)

		PageArray := getPages()
		fmt.Println(PageArray)
		for _, p := range PageArray {
			fmt.Println(p)
		}

		newPage := convtoPage(AIName, IntName, ArtiMts, HumanMts)

		PageArray = AddtoPageArray(newPage, PageArray)
		fmt.Println(PageArray)

		jsonPageArray := toJson(PageArray)
		fmt.Println("jsonPageArray:")
		fmt.Println(jsonPageArray)
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
		AIName = r.FormValue("AIName")
		fmt.Println(AIName)
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
