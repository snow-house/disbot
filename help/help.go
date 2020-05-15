package help

import (
	"log"
	"io/ioutil"
)

func Show(filename string) string {

	// dir, err := ioutil.ReadDir(".")
	// if err != nil {
	// 	log.Println(err)
	// }

	// for _, f := range(dir) {
	// 	log.Println(f.Name())
	// }

	dat, err := ioutil.ReadFile("help/"+filename) 
	if err != nil {
		log.Println(err)
		return "no help available for " + filename 
	}

	return string(dat)
}