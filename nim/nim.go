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

func Find(query string) (name, tpb, s1, major string) {

	baseurl := "https://api.nim.aryuuu.ninja/get/nim/"

	// send get request to nim finder
	res, err := http.Get(baseurl+query)

	if err != nil {
		log.Println(err)
	}


	defer res.Body.Close()

	if (res.StatusCode == 204) {
		return "nothing found :(", "tpb", "s1", "major"
	}


	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	
	nim, err := getNims([]byte(body))

	max := 15
	if nim.Count < max {
		max = nim.Count
	}

	n := ""
	t := ""
	s := ""
	m := ""

	for i := 0; i < max; i++ {
		n = fmt.Sprintf("%s\n%s", n, nim.Data[i].Nama)
		t = fmt.Sprintf("%s\n%s", t, nim.Data[i].Tpb)
		s = fmt.Sprintf("%s\n%s", s, nim.Data[i].S1)
		m = fmt.Sprintf("%s\n%s", m, nim.Data[i].Jurusan)
	}

	return n, t, s, m
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
