package models

import (
	"bytes"
	"net/http"
	"WorkingPromo/Utils"
	"WorkingPromo/XMLParsers"
	"io/ioutil"
)

func UpdateBigFoldersInDB() {

	//Типа новый хттп клиент дофига создаем,так же создаем реквест
	//делаем в хедаре метод пост, выбириаем урл, херачим наш хмл
	//из пакеда XMLParsers в байты, чекаем что ошибок нету и передаем реквест клиенту типа замути.
	//Получаем ответ и закрываем поток респанса боди чтобы ниче не упало попутно чекая ошибки.
	//закрывается все после этих прыжков с бубном

	//клиент и реквест
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://plati.io/xml/sections.asp",
		bytes.NewBuffer(XMLParsers.GetSectionsXML()))
	Utils.CheckError(err)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	Utils.CheckError(err)

	//Читаем байты с тела (хмл), переводим хмл в структуры нашего проекта благодаря чудо утилите в пакте XML
	//Parsers
	responseData, err := ioutil.ReadAll(resp.Body)
	Utils.CheckError(err)
	structs := XMLParsers.GetSectionsViaStructXML(responseData)

	//Создаем транзакцию (tx) и припейрейд стейтмент который захерачим в бд. В постгрес для индексации параметров
	//в prepared statement используется $
	tx, err := db.Begin()

	//TODO надо переделать стейтмент чтобы он не тупо добавлял новые записи, а при конфликте ключа обновлял те что уже есть
	//, у меня где то валялось, надо будет найти не забыть
	stmt, err := tx.Prepare("insert into public.big_folders (id, name_folder,count_goods) values ($1,$2,$3)")
	Utils.CheckError(err)

	//Тут мы опять говорим типа после танцев с бубном верни обратно в пул конекшен к базе
	defer stmt.Close()

	//Кидаем в бдшку наши данные
	for count := range structs.Folders{
		_, err = stmt.Exec(structs.Folders[count].FolderID, structs.Folders[count].FolderName,
			structs.Folders[count].CNTGoods)
		Utils.CheckError(err)
	}
	//Коммитим нашу транзакцию, вообще можно если есть ошибка делать rollback но я хз пока, для тестов и тай сойдет
	err = tx.Commit()
	Utils.CheckError(err)

}