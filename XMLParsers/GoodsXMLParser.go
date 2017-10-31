package XMLParsers

import (
	"encoding/xml"
)

type GoodsInner struct {
	guid_agent string `xml:"guid_agent"`

}

type GoodsXML struct {
	XMLName  xml.Name `xml:"digiseller.request"`
	Document GoodsInner `xml:"document"`
}

