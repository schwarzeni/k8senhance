package main

import (
	"html/template"
	"log"
	"net/http"
)

type Data struct {
	Name    string
	Num     int
	SubData struct {
		SubName string
		Items   []struct {
			Name string
			Num  int
		}
	}
}

func main() {
	data := Data{
		Name: "header",
		Num:  233,
		SubData: struct {
			SubName string
			Items   []struct {
				Name string
				Num  int
			}
		}{
			SubName: "subheader",
			Items: []struct {
				Name string
				Num  int
			}{
				{Name: "m1", Num: 1},
				{Name: "m2", Num: 2},
				{Name: "m3", Num: 3},
			},
		},
	}

	temp, err := template.ParseFiles("node.html", "menu.html")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/aa", func(resp http.ResponseWriter, req *http.Request) {
		log.Println("hit")
		if err := temp.Execute(resp, data); err != nil {
			panic(err)
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
