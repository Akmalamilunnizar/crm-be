package helpers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func HttpRequestHelpers(url string, method string, payload string) (string, error) {

	// url := "http://103.148.18.189/rest/ip/hotspot/ip-binding"
	// method := "PUT"

	// 	payload := strings.NewReader(`{
	//     "mac-address": "14:0b:22:1F:00:14",
	//     "to-address": "192.192.192.192",
	//     "address": "192.192.192.192",
	//     "type": "bypassed"
	// }`)

	payloadReader := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic cm5kQGxpbGx5Lm5ldC5pZDpybmRAMTIz")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err

	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err

	}
	fmt.Println(string(body))
	return string(body), nil
}

type PutMikrotik struct {
	MacAddress string `json:"mac-address"`
	ToAddress  string `json:"to-address"`
	Address    string `json:"address"`
	Type       string `json:"type"`
}
