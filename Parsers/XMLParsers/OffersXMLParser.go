package XMLParsers

import (
	"encoding/xml"
	"WorkingPromo/Utils"
)

//Главный хмл по предложениям по секции который мы получаем
//что то возможно проигнорировал, сейчас не вспомню
type OffersXML struct {
	XMLName     xml.Name      `xml:"digiseller.response"`
	Retval      string        `xml:"retval"`
	RetDesc     string        `xml:"retdesc"`
	IdSection   string        `xml:"id_section"`
	NameSection string        `xml:"name_section"`
	Page        string        `xml:"page"`
	Order       string        `xml:"order"`
	CountGoods  string        `xml:"cnt_good"`
	Pages       string        `xml:"pages"`
	Rows        OfferRows     `xml:"rows"`
}

//Массив строк, есть счетчик сколько в ответе будет товаров
type OfferRows struct {
	RowsCount string        `xml:"cnt,attr"`
	OfferRows []OfferRowXML `xml:"row"`
}

//Статистика по продавцу, идет вместе с предложением
type OfferStatistic struct {
	CountSell          string `xml:"cnt_sell"`
	CountReturn        string `xml:"cnt_return"`
	CountGoodResponses string `xml:"cnt_goodresponses"`
	CountBadResponses  string `xml:"cnt_badresponses"`
}

//Строка по определенному предложению, вся инфа что есть
type OfferRowXML struct {
	OfferRowID string         `xml:"id,attr"`
	OfferID    string         `xml:"id_goods"`
	OffersName string         `xml:"name_goods"`
	Price      string         `xml:"price"`
	Currency   string         `xml:"currency"`
	Discount   string         `xml:"discount"`
	Gift       string         `xml:"gift"`
	Reward     string         `xml:"reward"`
	IdSeller   string         `xml:"id_seller"`
	NameSeller string         `xml:"name_seller"`
	Rating     string         `xml:"rating"`
	Summpay    string         `xml:"summpay"`
	Statistics OfferStatistic `xml:"statistics"`
}

//Этим мы требуем у партнерки список предложений по секции
//Важно! обязательно нужно указывать количество строк на страницу и валюту, ну и само собой
//агент айди
type GetOffersXML struct {
	XMLName   xml.Name `xml:"digiseller.request"`
	AgentId   string   `xml:"guid_agent"`
	IDSection string   `xml:"id_section"`
	Lang      string   `xml:"lang"`
	Encoding  string   `xml:"encoding"`
	Page      string   `xml:"page"`
	Rows      string   `xml:"rows"`
	Currency  string   `xml:"currency"`
	Order     string   `xml:"order"`

}

//Требуем в виде байт хмл для запроса предложений по секции
func GetOffersInSectionXML(section string, rowsNumber string, currency string) ([]byte) {
	v := &GetOffersXML{AgentId: Utils.AgentID, IDSection: section,
		Lang: Utils.Lang, Encoding: "utf-8", Rows: rowsNumber, Currency: currency}
	getOffersInSectionXML, err := xml.MarshalIndent(&v, "", "  ")
	Utils.CheckError(err)
	return getOffersInSectionXML
}

//Переводим байты в структуры, как и в секциях вообщем, ничего интересного
func GetOffersViaStructXML(data []byte) (*OffersXML) {
	var sections = &OffersXML{}
	err := xml.Unmarshal(data, &sections)
	Utils.CheckError(err)
	return sections
}

