package botclient

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func getResult(response *http.Response, result any) {
	all, err := io.ReadAll(response.Body)
	if err != nil {
		log.Panicln(err.Error())
	}
	err = json.Unmarshal(all, result)
	if err != nil {
		log.Panicln(err)
	}
}

type UrlResponse struct {
	Url string
}

// ==================== Messages ==================== //
type InitMessage struct {
	Op int
	D  struct {
		HeartbeatInterval int
	}
}
