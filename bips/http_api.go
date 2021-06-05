package bips

import (
	"encoding/json"
	"log"
)

func MarshalJSON(in interface{}) string {
	bs, e := json.Marshal(in)
	if e != nil {
		log.Println("marshal json fail:", e)
		return "{}"
	}
	return string(bs)
}

func MarshalJSONIndent(in interface{}) string {
	bs, e := json.MarshalIndent(in, "", "  ")
	if e != nil {
		log.Println("marshal json with indent fail:", e)
		return "{}"
	}
	return string(bs)
}
