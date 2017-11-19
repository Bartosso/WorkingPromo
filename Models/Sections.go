package models

import (
	"bytes"
	"net/http"
	"WorkingPromo/Utils"
	"WorkingPromo/Parsers/XMLParsers"
	"WorkingPromo/Parsers/JsonParsers"
	"encoding/json"
	"io/ioutil"
	"fmt"
)
//Обновляет в бд папки игр и секции в данных папках
func updateGameFoldersAndInnerSectionsInDB(xmlStruct *XMLParsers.SectionXML)  {

	tx, err := db.Begin()
	Utils.CheckError(err)

	stmt, err := tx.Prepare("insert into public.games_folders (id, name_folder,count_goods) values ($1,$2,$3) " +
		"on conflict (id) do update set name_folder = $2, count_goods = $3 ")
	Utils.CheckError(err)
	defer stmt.Close()

	innerstmt, err := tx.Prepare("insert into public.games_sections " +
		"(id, name_section,count_goods,parent_folder) values ($1,$2,$3,$4) " +
		"on conflict (id) do update set name_section = $2, count_goods = $3, parent_folder = $4 ")
	Utils.CheckError(err)
	defer innerstmt.Close()



	//Ищем игоры и переводим их в бд, 7682 это айди игор
	for count := range xmlStruct.Folders {
		if xmlStruct.Folders[count].FolderID == "7682" {
			//Какой то лапшекод, надо посмотреть мб возможно все красивше сделать
			gamesFolders := xmlStruct.Folders[count]
			for countInner := range gamesFolders.InnerFolders {
				_, err = stmt.Exec(
					gamesFolders.InnerFolders[countInner].FolderID,
					gamesFolders.InnerFolders[countInner].FolderName,
					gamesFolders.InnerFolders[countInner].CNTGoods)
				Utils.CheckError(err)
				folderInGamesFolders := gamesFolders.InnerFolders[countInner]
				for innerCounterFromInner := range folderInGamesFolders.InnerSections{
					innerstmt.Exec(
						folderInGamesFolders.InnerSections[innerCounterFromInner].SectionID,
						folderInGamesFolders.InnerSections[innerCounterFromInner].SectionName,
						folderInGamesFolders.InnerSections[innerCounterFromInner].CNTGoods,
						folderInGamesFolders.FolderID)
				Utils.CheckError(err)
				}
			}
			break
		}
	}
	err = tx.Commit()
	Utils.CheckError(err)
}

//Обновляет игры в бд без папок, т.е. те которые секции
func updateGameSectionsWithoutFoldersInDB(xmlStruct *XMLParsers.SectionXML) {
	tx, err := db.Begin()
	Utils.CheckError(err)

	//Т.к. паповский папки у данных секций нету то в поле парент фолдер передаем ноль чтобы быть осведомленным
	//можно конечно не заполнять её но хз, лучше перестраховатся я думаю
	stmt, err := tx.Prepare("insert into public.games_sections " +
		"(id, name_section,count_goods,parent_folder) values ($1,$2,$3,'0') " +
		"ON CONFLICT (id) DO UPDATE SET name_section = $2, count_goods = $3, parent_folder = '0' ")
	Utils.CheckError(err)
	defer stmt.Close()

	for count := range xmlStruct.Folders {
		if xmlStruct.Folders[count].FolderID == "7682" {
			gameFolder := xmlStruct.Folders[count]
			for innerCounter := range gameFolder.InnerSections {
				section := gameFolder.InnerSections[innerCounter]
				_, err := stmt.Exec(
					section.SectionID,
					section.SectionName,
					section.CNTGoods)
				Utils.CheckError(err)
			}
		}
	}
	err = tx.Commit()
	Utils.CheckError(err)
}

//Мего функция которая обновляет данные в бд по играм, и с папками и без
func updateGamesInDB(xml *XMLParsers.SectionXML) {
	updateGameFoldersAndInnerSectionsInDB(xml)
	updateGameSectionsWithoutFoldersInDB(xml)
}

//Функция для обновления больших папок в бд
func updateBigFoldersInBD(xml *XMLParsers.SectionXML) {
	//Создаем транзакцию (tx) и припейрейд стейтмент который захерачим в бд. В постгрес для индексации параметров
	//в prepared statement используется $
	tx, err := db.Begin()

	stmt, err := tx.Prepare("insert into public.big_folders (id, name_folder,count_goods) values ($1,$2,$3) " +
		"ON CONFLICT (id) DO UPDATE SET name_folder = $2, count_goods = $3 ")
	Utils.CheckError(err)

	//Тут мы опять говорим типа после танцев с бубном верни обратно в пул конекшен к базе
	defer stmt.Close()

	//Кидаем в бдшку наши данные
	for count := range xml.Folders{
		_, err = stmt.Exec(
			xml.Folders[count].FolderID,
			xml.Folders[count].FolderName,
			xml.Folders[count].CNTGoods)
		Utils.CheckError(err)
	}
	//Коммитим нашу транзакцию, вообще можно если есть ошибка делать rollback но я хз пока, для тестов и тай сойдет
	err = tx.Commit()
	Utils.CheckError(err)
}

func UpdateBigFoldersAndGamesFromSections() {

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


	responseData, err := ioutil.ReadAll(resp.Body)
	Utils.CheckError(err)
	structs := XMLParsers.GetSectionsViaStructXML(responseData)

	//Зовем мего метод который обновляет в бд большие папки и другую функцию которая обновляет
	//все игры в бд
	updateBigFoldersInBD(structs)
	updateGamesInDB(structs)

}

//Достает все секции из бд в виде массива стринг и отдаем нам
func GetSectionsIDFromDB() []string {
	rows, err := db.Query("select id from games_sections")
	Utils.CheckError(err)

	defer rows.Close()

	var slice = make([]string,0)

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		Utils.CheckError(err)

		slice = append(slice,id)
	}

	//test

	fmt.Println("Got slice!")

	//test
	return slice
}

func GetAvailableGamesSectionsAndFoldersViaJson() *[]byte {
	rows, err := db.Query("select * from games_folders")
	Utils.CheckError(err)

	defer rows.Close()
	var slice = make([]JsonParsers.GamesSectionsAndFoldersJson,0)

	for rows.Next() {
		var folder JsonParsers.GamesSectionsAndFoldersJson
		err = rows.Scan(&folder.Id,&folder.Name,&folder.Count)
		folder.Type = "gamefolder"

		Utils.CheckError(err)

		slice = append(slice,folder)
	}
	rows, err = db.Query("select id,name_section,count_goods from games_sections where parent_folder = '0'")
	Utils.CheckError(err)

	for rows.Next() {
		var folder JsonParsers.GamesSectionsAndFoldersJson
		err = rows.Scan(&folder.Id,&folder.Name,&folder.Count)
		folder.Type = "gamesection"

		Utils.CheckError(err)

		slice = append(slice,folder)
	}

	jsone,_ := json.Marshal(slice)
	return &jsone

}

func GetSelectedByParentIDSections(parentId string) *[]byte {
	rows, err := db.Query("select * from games_folders")
	Utils.CheckError(err)

	defer rows.Close()
	var slice = make([]JsonParsers.GamesSectionsAndFoldersJson,0)

	rows, err =
		db.Query("select id,name_section,count_goods from games_sections where parent_folder = $1",parentId)
	Utils.CheckError(err)

	for rows.Next() {
		var folder JsonParsers.GamesSectionsAndFoldersJson
		err = rows.Scan(&folder.Id,&folder.Name,&folder.Count)
		folder.Type = "gamesection"

		Utils.CheckError(err)

		slice = append(slice,folder)
	}
	if len(slice) == 0 {
		return nil
	}
	jsone,_ := json.Marshal(slice)

	return &jsone

}