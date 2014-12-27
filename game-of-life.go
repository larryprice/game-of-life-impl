package main

import (
  "net/http"
  "net/url"
  "encoding/json"
  "log"
  "strconv"
)

type cell struct {
    Column int `json:"column"`
    Row int `json:"row"`
}

type stepArguments struct {
    Cells []cell `json:"cells"`
    Steps int `json:"steps"`
    Direction int `json:"direction"`
}

func NewStepArguments(query url.Values) stepArguments {
    var cells []cell
    json.Unmarshal([]byte(query["cells"][0]), &cells)

    steps, _ := strconv.Atoi(query["steps"][0])
    direction, _ := strconv.Atoi(query["direction"][0])

    return stepArguments{cells, steps, direction}
}

func stepHandler(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Access-Control-Allow-Origin", "*")
    res.Header().Add("Access-Control-Allow-Headers", "Content-Type")

    req.ParseForm()
    log.Println(NewStepArguments(req.Form))
}

func main() {
    http.HandleFunc("/", stepHandler)
    http.ListenAndServe(":5001", nil)
}