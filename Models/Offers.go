package models

import (
	"WorkingPromo/XMLParsers"
	"WorkingPromo/Utils"
	"database/sql"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
)

//Берет на вход так же транзакцию, кидает в бд инфу по офферу не комитит, ничего интересного
func updateGamesOffersInDB(offersXML *XMLParsers.OffersXML, tx *sql.Tx) {

	stmt, err := tx.Prepare("insert into public.games_offers (id, parent_section_id, offer_name," +
		"price, currency, discount, gift, reward, id_seller) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)")
	Utils.CheckError(err)
	defer stmt.Close()

	offerRows := offersXML.Rows.OfferRows

	for counter := range offerRows{
		_, err = stmt.Exec(
			offerRows[counter].OfferID,
			offersXML.IdSection,
			// ^ Айди родителя, типа секция где находится данный товар, важно пздц, было бы замечательно в будущем
			//пометить в дб что данный столб нот нал
			offerRows[counter].OffersName,
			offerRows[counter].Price,
			offerRows[counter].Currency,
			offerRows[counter].Discount,
			offerRows[counter].Gift,
			offerRows[counter].Reward,
			offerRows[counter].IdSeller)
		Utils.CheckError(err)

		fmt.Println("added new offer! " + offerRows[counter].OffersName)
	}
}

//Берет на вход транзакцию, кидает в бд инфу по продавцу, не коммитит транзакцию
func updateGamesOffersSellersStats(offersXML *XMLParsers.OffersXML, tx *sql.Tx)  {
	stmt, err := tx.Prepare("insert into public.games_sellers_stats (seller_id, seller_name, rating," +
		"summpay, cnt_sell, cnt_return, cnt_goodresponses, cnt_badresponses) values ($1,$2,$3,$4,$5,$6,$7,$8) " +
			"ON CONFLICT (seller_id) DO NOTHING")
	Utils.CheckError(err)
	defer stmt.Close()

	offerRows := offersXML.Rows.OfferRows

	for counter := range offerRows{
		_, err = stmt.Exec(
			offerRows[counter].IdSeller,
			offerRows[counter].NameSeller,
			offerRows[counter].Rating,
			offerRows[counter].Summpay,
			offerRows[counter].Statistics.CountSell,
			offerRows[counter].Statistics.CountReturn,
			offerRows[counter].Statistics.CountGoodResponses,
			offerRows[counter].Statistics.CountBadResponses)
		Utils.CheckError(err)

		fmt.Println("added new seller! " + offerRows[counter].NameSeller)
	}

}
//Todo облагородить всю эту херню, мб где то накосячил
//Делает партнерский запрос по секции, переводит в структуры и возвращает, ничего интересного
//изначально хотел оставить в нижнем методе но понял что будет косяк с defer resp.body.close
func getOffersFromSectionFromWeb(sectionId string) *XMLParsers.OffersXML {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://plati.io/xml/goods.asp",
		bytes.NewBuffer(XMLParsers.GetOffersInSectionXML(sectionId,"500","RUR")))
	Utils.CheckError(err)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	Utils.CheckError(err)

	responseData, err := ioutil.ReadAll(resp.Body)
	Utils.CheckError(err)
	structs := XMLParsers.GetOffersViaStructXML(responseData)
	//

	fmt.Println("got structs!!!!")

	//
	return structs
}

//Обновляет офферы и продавцов, тупо форыч
func UpdateGamesOffersAndSellersStats() {
	sectionsIds := GetSectionsIDFromDB()

	tx, err := db.Begin()
	//Зачищаем табличку перед проведением обряда, в дальнейшем я думаю роллбак делать если что то не так
	tx.Exec("truncate table games_offers; truncate TABLE games_sellers_stats")
	//

	fmt.Println("truncate done!")

	//
	Utils.CheckError(err)

	for counter := range sectionsIds {
		offersBySectionXML := getOffersFromSectionFromWeb(sectionsIds[counter])
		updateGamesOffersInDB(offersBySectionXML,tx)
		updateGamesOffersSellersStats(offersBySectionXML,tx)
	}
	err = tx.Commit()
	Utils.CheckError(err)


	fmt.Println("Update Done!")
}

