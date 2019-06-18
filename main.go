package main
import (
  "net/http"
  "strings"
  "encoding/json"
  "github.com/extrame/xls"
  "fmt"
  "os"
  "log"
)

type xldata struct {
	Markets string `json:"market"`
	MarPrice string `json:"price"`
	EffPrice string `json:"eff_price"`
  Delta string `json:"delta"`
  Expenses string `json:"expense"`
}

type jsonResponseQuery struct {
	Message string `json:"message"`
	Result []string `json:"result"`
}

type jsonResponse struct {
	Message string `json:"message"`
	Result string `json:"result"`
}

func fetchData(w http.ResponseWriter, r *http.Request) {

  w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
  r.ParseForm()

  request := r.FormValue("commodity")
  country := r.FormValue("country")
  market := r.FormValue("market")

  fmt.Println(request)
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

  if (strings.ToLower(request) == "turmeric" ) {
    var response []string 
    for _, data:= range ResultArray {
      response  = append(response, data.Markets)
    } 
    
    responseJSON := jsonResponseQuery {
				Message: "Here are the top 10 options to export Turmeric, in decreasing of profits",
				Result: response,
			}
    jData, _ := json.Marshal(responseJSON)
    w.Write(jData)
  } else if (strings.ToLower(market) == "best market" ) {
    
    responseJSON := jsonResponse {
				Message: "",
				Result: "The selling price in USA is Rs 390.60/Kg. The costs incurred to complete export will Rs. 23.86/Kg (Shipping and transportation Rs 12.14/Kg, Insurance Rs 7.21/Kg, Nudge Charges Rs. 15.62/Kg and Income of Rs 11.72/Kg as Govt subsidies). Effectively you will earn Rs 366.74/Kg, which is 432.92% higher than the price of India's domestic market.",
			}
    jData, _ := json.Marshal(responseJSON)
    w.Write(jData)
  } else if (strings.ToLower(country) == "usa" || strings.ToLower(country) == "australia" || strings.ToLower(country) == "canada" || strings.ToLower(country) == "germany" ||
      strings.ToLower(country) == "qatar" || strings.ToLower(country) == "uk" || strings.ToLower(country) == "netherlands"|| strings.ToLower(country) == "nigeria" ||
      strings.ToLower(country) == "uae" || strings.ToLower(country) == "sa") {

        var response string
        for _, data:= range ResultArray {
          if (strings.ToLower(data.Markets) == strings.ToLower(country)) {
            MarPrice := data.MarPrice
            Expenses := data.Expenses
            EffPrice := data.EffPrice
            Delta := data.Delta
            LMarPrice := len(MarPrice)
            LExpenses := len(Expenses)
            LEffPrice := len(EffPrice)
            LDelta := len(Delta)
            response = "The selling price in "+ data.Markets + " is Rs. " + MarPrice[:LMarPrice-2] + "/kg. The costs incured to complete export will be Rs. " + data.Expenses[:LExpenses-2] + "/kg. Effectively you will earn " + EffPrice[:LEffPrice-2] + "/kg, which is " + Delta[:LDelta-2] + "% higher than tha price in India (domestic Market)."
            break;
          }
        }

        responseJSON := jsonResponse {
                Message: "The Details are: ",
                Result: response,
            }

        jData, _ := json.Marshal(responseJSON)
        w.Write(jData)
    } else {

      response := "Please enter a valid query"
      responseJSON := jsonResponse {
                Message: "Invalid Input",
                Result: response,
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