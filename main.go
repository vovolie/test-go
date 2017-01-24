package main

import (
	"qiushibaike.com/test-go/data"
	"log"
	"github.com/graphql-go/handler"
	"net/http"
	"fmt"
	"os"
)

//func main() {
//	data.MaxMultiQueryCount = 300
//	if err := data.Init(); err != nil {
//		fmt.Println("init error", err)
//		os.Exit(1)
//	}
//
//	test_m := new(data.Material)
//	test_m.ID = 1
//	test_m.Cover = "http://7qn8gk.com1.z0.glb.clouddn.com/Fg--10OOqcl05a7TABv2Qrx5BAjv"
//	test_m.Name = "haha"
//	test_m.URL = "http://7qn8gk.com1.z0.glb.clouddn.com/Fg--10OOqcl05a7TABv2Qrx5BAjv"
//	test_m.Sha = "10OOqcl05a7TABv2Qrx5BAjv"
//	test_m.Version = "alsdkjflkasdjf"
//	test_m.MateInfo = "jflajoefijalsjdfasjdfjasd;fj"
//	test_m.HiddenAt = 1
//	test_m.CreatedAt = "2017-01-20 00:00:00"
//	test_m.Category = 1
//	err := data.InsertMaterial(test_m)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	test_two, err := data.GetMaterialById(strconv.Itoa(1))
//	fmt.Println(test_two, err)
//
//}

func main() {

	data.MaxMultiQueryCount = 300
	if err := data.Init(); err != nil {
		fmt.Println("init error", err)
		os.Exit(1)
	}

	h := handler.New(&handler.Config{
		Schema: &data.Schema,
		Pretty: true,
	})

	http.Handle("/graphql", h)

	port := ":9999"
	log.Printf(`GraphQL server starting up on http://localhost%v`, port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("ListenAndServe failed, %v", err)
	}
}
