package internal

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"log"
	"regexp"
)

type XmlEntry struct {
	Id string `xml:"id"` // id of this XmlEntry
	Value xml.CharData `xml:",innerxml"` // un-marshalled data
}

// Unescape characters escaped by sharepoint
// https:://docs.microsoft.com/en-us/dotnet/api/system.xml.xmlconvert.encodename?view=netframework-4.7.2#System_Xml_XmlConvert_EncodeName_System_String_
func ReplaceEscapedXml(src []byte) []byte {
	re := regexp.MustCompile(`_x[0-f]{4}_`)
	return re.ReplaceAllFunc(src, xmlUnescape)
}

func xmlUnescape(escaped []byte) (unescaped []byte) {
	// deal with shitty stuff
	if string(escaped) == "_x000a_" {
		return []byte{'_'}
	}

	hexAsString := escaped[2:6]
	value, _ := hex.DecodeString(string(hexAsString))

	e := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	es, _, err := transform.Bytes(e.NewDecoder(), value)
	if err != nil {
		panic(err)
	}
	fmt.Printf("transformed %v to \"%s\"", escaped, es)
	return es
}



func ParseSharepointEntries(r io.Reader) (entries []XmlEntry, err error) {
	dec := xml.NewDecoder(r)

	for ; ;  {
		t, err := dec.Token()
		if t == nil || err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error decoding token ", err)
			return nil, err
		}

		switch ty := t.(type) {
		case xml.StartElement:
			if ty.Name.Local == "entry" {
				// If this is a start element named "entry", parse this element into a XmlEntry
				var entry XmlEntry
				if err = dec.DecodeElement(&entry, &ty); err != nil {
					log.Println("Error decoding item ", err)
					return nil, err
				}
				entries = append(entries, entry)
			}
		default:
		}
	}

	return entries, nil
	}