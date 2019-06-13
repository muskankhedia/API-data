package main
import (
  "net/http"
  "strings"
  "encoding/json"
  "github.com/extrame/xls"
  "fmt"
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

  r.ParseForm()

  request := r.FormValue("commodity")

  fmt.Println(request)
  var Result xldata
  var ResultArray []xldata

  if xlFile, err := xls.Open("./Table.xls", "utf-8"); err == nil {
    fmt.Print("reached")
      if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
          fmt.Print("Total Lines ", sheet1.MaxRow, sheet1.Name)
          // fmt.Print("Total Columns : ", sheet1.MaxColumn)
          for i := 1; i <= (int(sheet1.MaxRow)); i++ {
              row1 := sheet1.Row(i)
              Result.Markets = row1.Col(0)
              Result.MarPrice = row1.Col(1)
              Result.EffPrice = row1.Col(2)
              Result.Delta = row1.Col(3)
              Result.Expenses = row1.Col(4)
              fmt.Print("\n", Result)
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
  } else if (strings.ToLower(request) == "united states of america" || strings.ToLower(request) == "australia" || strings.ToLower(request) == "canada" || strings.ToLower(request) == "germany" ||
      strings.ToLower(request) == "qatar" || strings.ToLower(request) == "united kingdom" || strings.ToLower(request) == "netherlands"|| strings.ToLower(request) == "nigeria" ||
      strings.ToLower(request) == "united arab emirates" || strings.ToLower(request) == "saudi arabia") {

        var response string
        for _, data:= range ResultArray {
          if (strings.ToLower(data.Markets) == strings.ToLower(request)) {
            response = "The selling price in "+ data.Markets + " is Rs. " + data.MarPrice + " /kg. The costs incured to complete export will be Rs. " + data.Expenses + " /kg. Effectively you will earn " + data.EffPrice + " /kg, which is " + data.Delta + "% higher than tha price in India (domestic Market)."
            break;
          }
        }

        responseJSON := jsonResponse {
                Message: "The Details are: ",
                Result: response,
            }

        jData, _ := json.Marshal(responseJSON)
        w.Write(jData)
    }  
} 



func main() {
  http.HandleFunc("/", fetchData)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
  // init()
}