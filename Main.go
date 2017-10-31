package main

import (
	"fmt"
	"net/http"
	_"github.com/lib/pq"
	"WorkingPromo/Models"
	"WorkingPromo/Utils"
	"WorkingPromo/XMLParsers"
	"bytes"
	"io/ioutil"
)

//Единственный недохендлер, юзаю для тестов
func handler(w http.ResponseWriter, r *http.Request) {

	//Типа новый хттп клиент дофига создаем,так же создаем реквест
	// делаем в хедаре метод пост, выбириаем урл, херачим наш хмл
	// из пакеда XMLParsers в байты, чекаем что ошибок нету и передаем реквест клиенту типа замути.
	// Получаем ответ и закрываем поток респанса боди чтобы ниче не упало попутно чекая ошибки.
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://plati.io/xml/sections.asp",
		bytes.NewBuffer([]byte(XMLParsers.GetSectionsXML())))
	Utils.CheckError(err)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	Utils.CheckError(err)

	//Читает байты с бади респанса и херачим их в стрингу, ну а потом печатаем, теперь мы знаем что нам
	//делать дальше! Юху!
	responseData, err := ioutil.ReadAll(resp.Body)
	Utils.CheckError(err)
	responseString := string(responseData)


	fmt.Println(responseString)
	fmt.Fprint(w, "Response \n" + responseString)
	fmt.Fprintln(w , resp)

	//
}

func main() {
	models.InitDB(fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		Utils.DB_USER, Utils.DB_PASSWORD, Utils.DB_NAME))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}



