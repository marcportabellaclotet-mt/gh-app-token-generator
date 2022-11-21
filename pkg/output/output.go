package output

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Token  string `json:"token"`
	Status string `json:"status"`
	Info   string `json:"info"`
	Output string `json:"output,omitempty"`
}

func ReturnResponse(response Response) {

	switch output := response.Output; output {
	case "json":
		r := Response{
			Status: response.Status,
			Token:  response.Token,
			Info:   response.Info,
		}
		jr, _ := json.Marshal(r)
		fmt.Printf(string(jr))
		fmt.Println("")
		return
	case "txt":
		if response.Status == "success" {
			fmt.Println(response.Token)
			return
		} else {
			fmt.Printf("Error: %s", response.Info)
			return
		}
	case "export":
		if response.Status == "success" {
			fmt.Printf("export GH_TOKEN=%s", response.Token)
		} else {
			fmt.Printf("Error: %s", response.Info)
			return
		}
	}
}
