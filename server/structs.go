package main

import "encoding/xml"

//Tag for json structure
type Tag struct {
	Writable    bool              `json:"writable"`
	Path        string            `json:"path"`
	Group       string            `json:"group"`
	Description map[string]string `json:"description"`
	Type        string            `json:"type"`
}

//TagContainer of multiple tags for json
type TagContainer struct {
	Tags []Tag `json:"tags"`
}

//NewTagContainer is a TagContainer constructor
func NewTagContainer(tags []Tag) *TagContainer {
	return &TagContainer{Tags: tags}
}

//Table is a struct of string read from exif output
type Table struct {
	XMLName xml.Name `xml:"table"`
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name,attr"`
	G0      string   `xml:"g0,attr"`
	G1      string   `xml:"g1,attr"`
	G2      string   `xml:"g2,attr"`
	Desc    []struct {
		Text string `xml:",chardata"`
		Lang string `xml:"lang,attr"`
	} `xml:"desc"`
	Tag []struct {
		Text     string `xml:",chardata"`
		ID       string `xml:"id,attr"`
		Name     string `xml:"name,attr"`
		Type     string `xml:"type,attr"`
		Count    string `xml:"count,attr"`
		Writable bool   `xml:"writable,attr"`
		G2       string `xml:"g2,attr"`
		Desc     []struct {
			Text string `xml:",chardata"`
			Lang string `xml:"lang,attr"`
			Desc []struct {
				Text string `xml:",chardata"`
				Lang string `xml:"lang,attr"`
			} `xml:"desc"`
		} `xml:"desc"`
		Val []struct {
			Text string `xml:",chardata"`
			Lang string `xml:"lang,attr"`
		} `xml:"val"`
	} `xml:"tag"`
}
