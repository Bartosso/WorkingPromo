package XMLParsers

import (
	"encoding/xml"
	"WorkingPromo/Utils"
)

//Это наш хмл который мы кидаем с запросом XML
type GetSectionXML struct {
	XMLName   xml.Name `xml:"digiseller.request"`
	AgentId   string   `xml:"guid_agent"`
	IdCatalog string   `xml:"id_catalog"`
	Lang      string   `xml:"lang"`
	Encoding  string   `xml:"encoding"`
}

//Это папка, она может содержать в себе как и другие папки так и секции
//т.е. непосредственно игры и категории игр (допустим папку которая содержит секции по асасин криду и т.д.)
type SectionFolder struct {
	FolderID   string    `xml:"id,attr"`
	FolderName string    `xml:"name_folder"`
	CNTGoods   string    `xml:"cnt_goods"`
	InnerFolders  []SectionFolder `xml:"folder"`
	InnerSections []Section `xml:"section"`
}

//Это секция, она в себе содержит название секции (типа название игры допустим)
type Section struct {
	SectionID   string   `xml:"id,attr"`
	SectionName string   `xml:"name_section"`
	CNTGoods    string   `xml:"cnt_goods"`
}

//Главный xml документ, большой то есть, который держит в себе папки с разделами
//Игры, IP-Телефония и т.д.
type SectionXML struct {
	XMLName  xml.Name  `xml:"digiseller.response"`
	Folders  []SectionFolder `xml:"folder"`

}



//Рисует запрос к сервису требуя список разделов
func GetSectionsXML() ([]byte) {
	v := &GetSectionXML{AgentId: Utils.AgentID, IdCatalog: Utils.IDCatalogForSections,
	Lang: Utils.Lang, Encoding: "utf-8"}
	sectionXML, err := xml.MarshalIndent(&v, "", "  ")
	Utils.CheckError(err)
	return sectionXML
}

//Полученную дату переводит в структуру SectionXML, используется указатель ибо иначе анмаршел будет
//кидатся ошибками, честно я хз почему, да и имхо пофигу, чего память говном забивать
func GetSectionsViaStructXML(data []byte) (*SectionXML) {
	var sections = &SectionXML{}
	err := xml.Unmarshal(data, &sections)
	Utils.CheckError(err)
	return sections
}


