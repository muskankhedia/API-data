package main
import (
  "net/http"
  "strings"
  "encoding/json"
  "os"
  "log"
)

type jsonResponseQuery struct {
	Result []string `json:"result"`
}

func fetchData(w http.ResponseWriter, r *http.Request) {

  w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
  r.ParseForm()

  request := r.FormValue("commodity")

  if (strings.ToLower(request) == "turmeric" ) {

    response := []string {
      "Here is best market for your product based on Income potential (on 14/06/2019)",
      "USA (Market Price: Rs. 390.60/Kg, Effective Price: Rs. 366.74/Kg, Higher income by 423.92%)",
      "The selling price in USA is Rs 390.60/Kg. The costs incurred to complete export will Rs. 23.86/Kg (Shipping and transportation Rs 12.14/Kg, Insurance Rs 7.21/Kg, Nudge Charges Rs. 15.62/Kg and Income of Rs 11.72/Kg as Govt subsidies).",
      "Effectively you will earn Rs 366.74/Kg, which is 432.92% higher than the price of India's domestic market. Also check the following top 10 destinations for your product, in decreasing order of profitability.",
      "1. United States of America",
      "2. Australia",
      "3. Canada",
      "4. Germany",
      "5. Qatar",
    } 
    
    responseJSON := jsonResponseQuery {
				Result: response,
			}
    jData, _ := json.Marshal(responseJSON)
    w.Write(jData)
  } else {

      response := []string {"Please enter a valid query",}
      responseJSON := jsonResponseQuery {
                Result: response,
            }

        jData, _ := json.Marshal(responseJSON)
        w.Write(jData)
    }
} 

func main() {

  http.HandleFunc("/fetchdata", fetchData)
  port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
  }

  port = ":" + port

  if err := http.ListenAndServe(port, nil); err != nil {
    panic(err)
  }

}