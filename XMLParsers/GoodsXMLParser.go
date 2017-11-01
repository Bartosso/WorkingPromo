package XMLParsers

import (
	"encoding/xml"
)

//Тут типа сами товары будут, еще толком не занимался этим
type GoodsInner struct {
	guid_agent string `xml:"guid_agent"`

}

type GoodsXML struct {
	XMLName  xml.Name `xml:"digiseller.request"`
	Document GoodsInner `xml:"document"`
}

