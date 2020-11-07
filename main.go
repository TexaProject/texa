package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/md5"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/big"

	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/TexaProject/texajson"
	"github.com/TexaProject/texalib"

	store "github.com/TexaProject/store"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// AIName exports form value from /welcome globally
var AIName string

// IntName exports form value from /texa globally
var IntName string

//Config is used to extract data and credentials from config.json(this file is gitignore'ed)
type Config struct {
	EthereumRPCEndpoint    string `json:"ethereum_rpc_endpoint"`
	WalletPrivateKey       string `json:"wallet_privatekey"`
	StorageContractAddress string `json:"storage_contract_address"`
}

// GetConfigData reads the data from config.json and loads it into a variable
func GetConfigData() Config {
	var configData Config
	bytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Panicln("GetConfigData(): Issue in reading config file. Please have a check.")
	}
	json.Unmarshal(bytes, &configData)
	return configData
}

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
		QSA := r.Form.Get("scoreArray")
		SlabName := r.Form.Get("SlabName")
		slabSequence := r.Form.Get("slabSequence")
		justificationSequence := r.Form.Get("justificationSeq")
		transactionLines := r.Form.Get("elizaLines")

		if AIName == "" {
			AIName = "unknown"
		}
		fmt.Println("###", AIName)
		fmt.Println("###", IntName)
		fmt.Println("###", QSA)
		fmt.Println("###", SlabName)
		fmt.Println("###", slabSequence)
		fmt.Println("###", justificationSequence)
		fmt.Println("###", transactionLines)

		// LOGIC
		re := regexp.MustCompile("[0-1]+")
		array := re.FindAllString(QSA, -1)

		SlabNameArray := regexp.MustCompile("[,]").Split(SlabName, -1)
		slabSeqArray := regexp.MustCompile("[,]").Split(slabSequence, -1)

		justificationSeqArray := regexp.MustCompile("[\n]").Split(justificationSequence, -1)
		transactionSeqArray := regexp.MustCompile("[\n]").Split(transactionLines, -1)

		fmt.Println("###Resulting Array:")
		for x := range array {
			fmt.Println(array[x])
		}

		fmt.Println("###SlabNameArray: ")
		fmt.Println(SlabNameArray)

		fmt.Println("###slabSeqArray: ")
		fmt.Println(slabSeqArray)

		fmt.Println("###justificationSeqArray: ")
		fmt.Println(justificationSeqArray)

		fmt.Println("###transactionSeqArray: ")
		fmt.Println(transactionSeqArray)

		ArtiQSA := texalib.Convert(array)
		fmt.Println("###ArtiQSA:")
		fmt.Println(ArtiQSA)

		HumanQSA := texalib.SetHumanQSA(ArtiQSA)
		fmt.Println("###HumanQSA:")
		fmt.Println(HumanQSA)

		TSA := texalib.GetTransactionSeries(ArtiQSA, HumanQSA)
		fmt.Println("###TSA:")
		fmt.Println(TSA)

		ArtiMts := texalib.GetMeanTestScore(ArtiQSA)
		HumanMts := texalib.GetMeanTestScore(HumanQSA)

		fmt.Println("###ArtiMts: ", ArtiMts)
		fmt.Println("###HumanMts: ", HumanMts)

		PageArray := texajson.GetPages()
		fmt.Println("###PageArray")
		fmt.Println(PageArray)
		for _, p := range PageArray {
			fmt.Println(p)
		}

		newPage := texajson.ConvtoPage(AIName, IntName, ArtiMts, HumanMts)

		PageArray = texajson.AddtoPageArray(newPage, PageArray)
		fmt.Println("###AddedPageArray")
		fmt.Println(PageArray)

		JsonPageArray := texajson.ToJson(PageArray)
		fmt.Println("###jsonPageArray:")
		fmt.Println(JsonPageArray)

		////
		fmt.Println("### SLAB LOGIC")

		slabPageArray := texajson.GetSlabPages()
		fmt.Println("###slabPageArray")
		fmt.Println(slabPageArray)

		slabPages := texajson.ConvtoSlabPage(ArtiQSA, SlabNameArray, slabSeqArray, justificationSeqArray, transactionSeqArray)
		fmt.Println("###slabPages")
		fmt.Println(slabPages)
		for z := 0; z < len(slabPages); z++ {
			slabPageArray = texajson.AddtoSlabPageArray(slabPages[z], slabPageArray)
		}
		fmt.Println("###finalslabPageArray")
		fmt.Println(slabPageArray)

		JsonSlabPageArray := texajson.SlabToJson(slabPageArray)
		fmt.Println("###JsonSlabPageArray: ")
		fmt.Println(JsonSlabPageArray)

		////
		fmt.Println("### CAT LOGIC")

		CatPageArray := texajson.GetCatPages()
		fmt.Println("###CatPageArray")
		fmt.Println(CatPageArray)

		CatPages := texajson.ConvtoCatPage(AIName, slabPageArray, SlabNameArray)
		fmt.Println("###CatPages")
		fmt.Println(CatPages)
		CatPageArray = texajson.AddtoCatPageArray(CatPages, CatPageArray)

		fmt.Println("###finalCatPageArray")
		fmt.Println(CatPageArray)

		JsonCatPageArray := texajson.CatToJson(CatPageArray)
		fmt.Println("###JsonCatPageArray: ")
		fmt.Println(JsonCatPageArray)

		ResultObject := texajson.NewResultObject(AIName)

		newSessionData := texajson.NewInterrogationObject(IntName, ArtiMts, HumanMts, CatPages.CatVal)
		fmt.Println("PRINTING NEW SESSION DATA BEFORE ADDING: ", newSessionData)

		ResultObject.Interrogations = append(ResultObject.Interrogations, newSessionData)
		fmt.Println("PRINTING UPDATED RESULT OBJECT: ", ResultObject)

		cid := texajson.WriteDataToIPFS("https://ipfs.infura.io:5001", ResultObject)
		if len(cid) > 0 {
			fmt.Println("Successfully wrote the session data to IPFS at ", cid)
		}

		pin, err := texajson.GetInstance().PinCid(cid)
		if err != nil {
			log.Println("Pining failed for Cid:", cid, " error :", err)
		} else {
			fmt.Println("CID is successfully pinned : ", pin)
		}
		configData := GetConfigData()
		tx := SubmitTxnToBlockchain(configData, AIName, cid)

		txnURL := "https://kovan.etherscan.io/tx/" + tx
		globalURL := "https://explore.ipld.io/#/explore/" + cid

		fmt.Fprint(w, "<html><head><link rel=\"stylesheet\" href=\"http://localhost:3030/css/bootstrap.min.css\"><title>File Ack | TEXA Project</title></head><body>ACKNOWLEDGEMENT: Received the scores. <br /><br />Info:<br />")
		fmt.Fprint(w, "<br /><br />VISIT: <b> <a href=\"", globalURL, "\">", globalURL, "</a></b> for interrogation results. This link is public and the result data is written to the public IPFS network!")
		fmt.Fprint(w, "<br /><br />VISIT: <b> <a href=\"", txnURL, "\">", txnURL, "</a></b> for transaction details. This link is public and written to the Ethereum Kovan tesnet!")
		fmt.Fprintf(w, "<br /><br /><input type=\"button\" class=\"btn info\" style=\"border: 2px solid black;\" onclick=\"location.href='http://localhost:3030/result';\" value=\"Visit /result\" /></body></html>")
	}
}

// SubmitTxnToBlockchain is used to sign a transaction with the new CID of the session result data
// along with the name of the corresponding AI (to maintain the latest CID, yet maintain provenance)
func SubmitTxnToBlockchain(configData Config, AIName, cid string) string {
	client, err := ethclient.Dial(configData.EthereumRPCEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := crypto.HexToECDSA(configData.WalletPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
	address := common.HexToAddress(configData.StorageContractAddress)
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := instance.LogTexaResultURL(auth, AIName, cid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction committed for the new session: %s", tx.Hash().Hex())
	return tx.Hash().Hex()
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

// upload logic
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
		handler.Filename = "elizadata.js"
		AIName = r.FormValue("AIName")
		fmt.Println(AIName)
		defer file.Close()

		fmt.Fprint(w, "<html><head><link rel=\"stylesheet\" href=\"http://localhost:3030/css/bootstrap.min.css\"><title>File Ack | TEXA Project</title></head><body>ACKNOWLEDGEMENT: Uploaded the file. <br /><br />Header Info:<br />")
		fmt.Fprintf(w, "%v", handler.Header)
		fmt.Fprintf(w, "<br /><br />Saved As: www/js/"+handler.Filename)
		fmt.Fprint(w, "<br /><br />VISIT: /texa for interrogation.")
		fmt.Fprintf(w, "<br /><br /><input type=\"button\" class=\"btn info\" style=\"border: 2px solid black;\" onclick=\"location.href='http://localhost:3030/texa';\" value=\"Visit /texa\" /></body></html>")
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
	fmt.Println("method:", r.Method) //get	request	method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("www/result.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
	}
}

func getCatJSON(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get	request	method
	catPages := texajson.GetCatPages()
	bs, err := json.Marshal(catPages)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func getMtsJSON(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get	request	method
	mtsPage := texajson.GetPages()
	bs, err := json.Marshal(mtsPage)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func getSlabJSON(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get	request	method
	slabPages := texajson.GetSlabPages()
	bs, err := json.Marshal(slabPages)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func main() {

	host := flag.String("host", "localhost", "identifier for ipfs and ipfs cluster")
	port := flag.String("port", "9094", "Source port number")

	flag.Parse()
	texajson.InitCluster(*host, *port)

	fmt.Println("--TEXA SERVER--")
	fmt.Println("STATUS: INITIATED")
	fmt.Println("ADDR: http://127.0.0.1:3030")

	fsc := http.FileServer(http.Dir("www/css"))
	http.Handle("/css/", http.StripPrefix("/css/", fsc))
	fsj := http.FileServer(http.Dir("www/js"))
	http.Handle("/js/", http.StripPrefix("/js/", fsj))
	fsd := http.FileServer(http.Dir("www/data"))
	http.Handle("/data/", http.StripPrefix("/data/", fsd))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/welcome", welcomeHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/texa", texaHandler)
	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/cat", getCatJSON)
	http.HandleFunc("/mts", getMtsJSON)
	http.HandleFunc("/slab", getSlabJSON)

	http.ListenAndServe(":3030", nil)
}
