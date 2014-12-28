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

func addRow(rowNumber int, baseColumn int, genMap map[int]map[int]bool) {
    if row, ok := genMap[rowNumber]; ok {
        row[baseColumn] = row[baseColumn]
        row[baseColumn+1] = row[baseColumn+1]
        row[baseColumn-1] = row[baseColumn-1]
    } else {
        genMap[rowNumber] = map[int]bool{}
        genMap[rowNumber][baseColumn] = false
        genMap[rowNumber][baseColumn+1] = false
        genMap[rowNumber][baseColumn-1] = false
    }
}

func addSurroundingCells(c cell, genMap map[int]map[int]bool) {
    addRow(c.Row-1, c.Column, genMap)
    addRow(c.Row, c.Column, genMap)
    addRow(c.Row+1, c.Column, genMap)
}

func buildGenerationMap(generation []cell) map[int]map[int]bool {
    generationMap := map[int]map[int]bool{}
    for _, c := range generation {
        if row, ok := generationMap[c.Row]; ok {
            row[c.Column] = true
        } else {
            generationMap[c.Row] = map[int]bool{}
            generationMap[c.Row][c.Column] = true
        }

        addSurroundingCells(c, generationMap)
    }
    return generationMap
}

func getNeighborCount(baseRow int, baseColumn int, neighbors map[int]map[int]bool) int {
    neighborCount := 0
    for i := baseRow-1; i<=baseRow+1; i++ {
        for j := baseColumn-1; j<=baseColumn+1; j++ {
            if neighbors[i][j] && !(i == baseRow && j == baseColumn) {
                neighborCount++
            }
        }
    }

    return neighborCount
}

func getNextGeneration(liveCells []cell) []cell {
    nextGen := []cell{}
    generationMap := buildGenerationMap(liveCells)

    for k, v := range generationMap {
        for kk, vv := range v {
            count := getNeighborCount(k, kk, generationMap)
            if vv && (count == 2 || count == 3) {
                nextGen = append(nextGen, cell{kk, k})
            } else if !vv && count == 3 {
                nextGen = append(nextGen, cell{kk, k})
            }
        }
    }

    return nextGen
}

func stepHandler(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Access-Control-Allow-Origin", "*")
    res.Header().Add("Access-Control-Allow-Headers", "Content-Type")

    req.ParseForm()
    sa := NewStepArguments(req.Form)

    json.NewEncoder(res).Encode(getNextGeneration(sa.Cells))
}

func main() {
    http.HandleFunc("/", stepHandler)
    http.ListenAndServe(":5001", nil)
}