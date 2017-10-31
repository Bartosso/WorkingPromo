package XMLParsers

import (
	"encoding/xml"
	"WorkingPromo/Utils"
)

type SectionXML struct {
	XMLName   xml.Name `xml:"digiseller.request"`
	AgentId   string   `xml:"guid_agent"`
	IdCatalog string   `xml:"id_catalog"`
	Lang      string   `xml:"lang"`
	Encoding  string   `xml:"encoding"`
}




func GetSectionsXML() (string) {
	v := &SectionXML{AgentId: Utils.AgentID, IdCatalog: Utils.IDCatalogForSections,
	Lang: Utils.Lang, Encoding: "utf-8"}
	sectionXML, err := xml.MarshalIndent(&v, "", "  ")
	Utils.CheckError(err)
	return string(sectionXML)
}


