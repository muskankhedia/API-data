package main

import (
	"encoding/json"
	"fmt"
	"github.com/extrame/xls"
	"log"
	"net/http"
	"os"
	"strings"
)

type xldata struct {
	Markets  string `json:"market"`
	MarPrice string `json:"price"`
	EffPrice string `json:"eff_price"`
	Delta    string `json:"delta"`
	Expenses string `json:"expense"`
}

type jsonResponseQuery struct {
	Message string   `json:"message"`
	Result  []string `json:"result"`
}

type jsonResponse struct {
	Message string `json:"message"`
	Result  string `json:"result"`
}

func fetchData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()

	request := r.FormValue("commodity")
	country := r.FormValue("country")
  market := r.FormValue("market")
  LCountry := strings.ToLower(country)

  fmt.Println(request)
  fmt.Println(market)
  fmt.Println(LCountry)
	var Result xldata
	var ResultArray []xldata

	if xlFile, err := xls.Open("./Table1.xls", "utf-8"); err == nil {
		fmt.Print("reached")
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			for i := 1; i <= (int(sheet1.MaxRow)); i++ {
				row1 := sheet1.Row(i)
				Result.Markets = row1.Col(0)
				Result.MarPrice = row1.Col(1)
				Result.EffPrice = row1.Col(2)
				Result.Delta = row1.Col(3)
				Result.Expenses = row1.Col(4)
				ResultArray = append(ResultArray, Result)
			}
		}
	} else {
		panic(err)
	}

	if strings.ToLower(request) == "turmeric" {
		var response []string
		for _, data := range ResultArray {
			response = append(response, data.Markets)
		}

		responseJSON := jsonResponseQuery{
			Message: "Here are the top 10 options to export Turmeric, in decreasing of profits",
			Result:  response,
		}
		jData, _ := json.Marshal(responseJSON)
		w.Write(jData)
	} else if strings.ToLower(market) == "best" {

		responseJSON := jsonResponse{
			Message: "",
			Result:"The selling price in USA is Rs 390.60/Kg. The costs incurred to complete export will Rs. 23.86/Kg (Shipping and transportation Rs 12.14/Kg, Insurance Rs 7.21/Kg, Nudge Charges Rs. 15.62/Kg and Income of Rs 11.72/Kg as Govt subsidies). Effectively you will earn Rs 366.74/Kg, which is 432.92% higher than the price of India's domestic market.",
		}
		jData, _ := json.Marshal(responseJSON)
		w.Write(jData)
	} else if strings.Contains("usa", LCountry) || strings.Contains("australia", LCountry) || strings.Contains("canada", LCountry) || strings.Contains("germany", LCountry) ||
		strings.Contains("qatar", LCountry) {

		var response string
		for _, data := range ResultArray {
			if strings.ToLower(data.Markets) == strings.ToLower(country) {
        Price := data.MarPrice
        Expense := data.Expenses
        EffPrice := data.EffPrice
        Delta := data.Delta 
				response = "The selling price in " + data.Markets + " is Rs. " + Price + "/kg. The costs incured to complete export will be Rs. " + Expense + "/kg. Effectively you will earn " + EffPrice + "/kg, which is " + Delta + "% higher than tha price in India (domestic Market)."
				break
			}
		}

		responseJSON := jsonResponse{
			Message: "The Details are: ",
			Result:  response,
		}

		jData, _ := json.Marshal(responseJSON)
		w.Write(jData)
	} else {

		response := "Please enter a valid query"
		responseJSON := jsonResponse{
			Message: "Invalid Input",
			Result:  response,
		}

		jData, _ := json.Marshal(responseJSON)
		w.Write(jData)
	}
}

func main() {
	http.HandleFunc("/fetchdata", fetchData)
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	//   panic(err)
	// }
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	port = ":" + port

	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
