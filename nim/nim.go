package nim

import (
	"log"
	"fmt"
	// "strings"
	"net/http"
	// "reflect"
	"io/ioutil"
	"encoding/json"
)

type Student struct {
	Nama string `json:"nama"`
	Fakultas string `json:"fakultas"`
	Tpb string `json:"tpb"`
	S1 string `json:"s1"`
	S2 string `json:"s2"`
	S3 string `json:"s3"`
	Jurusan string `json:"jurusan"`
}

type NimAPIResponse struct {
	Message string `json:"message"`
	Count int `json:"count"`
	Data []Student `json:"data"`
}

func Find(query string) (string){

	baseurl := "http://35.240.223.196:6969/get/nim/"

	// send get request to nim finder
	res, err := http.Get(baseurl+query)

	if err != nil {
		log.Println(err)
	}


	defer res.Body.Close()

	if (res.StatusCode == 204) {
		return "nothing found :("
	}


	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	
	nim, err := getNims([]byte(body))

	reply := ""

	max := 15
	if nim.Count < max {
		max = nim.Count
	}

	for i := 0; i < max; i++ {
		reply = fmt.Sprintf("%s\n%s %s %s %s", 
							reply,
							nim.Data[i].Nama, 
							nim.Data[i].Tpb,
							nim.Data[i].S1,
							nim.Data[i].Jurusan)
	}

	return reply
}


func getNims(body []byte) (*NimAPIResponse, error) {
	nim := new(NimAPIResponse)
	// unmarshal response body
	err := json.Unmarshal(body, &nim)	
	if err != nil {
		log.Println(err)
	}

	return nim, err
}