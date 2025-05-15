package main

import "encoding/xml"

type EPP struct {
	XMLName  xml.Name `xml:"epp"`
	Xmlns    string   `xml:"xmlns,attr"`
	Greeting Greeting `xml:"greeting"`
}

type Greeting struct {
	SvID    string      `xml:"svID"`
	SvDate  string      `xml:"svDate"` // ISO 8601, UTC
	SvcMenu ServiceMenu `xml:"svcMenu"`
	DCP     DCP         `xml:"dcp,omitempty"`
}

type ServiceMenu struct {
	Version      []string      `xml:"version"`
	Lang         []string      `xml:"lang"`
	ObjURI       []string      `xml:"objURI"`
	SvcExtension *SvcExtension `xml:"svcExtension,omitempty"`
}

type SvcExtension struct {
	ExtURI []string `xml:"extURI"`
}

type DCP struct {
	Access    Access         `xml:"access"`
	Statement []DCPStatement `xml:"statement"`
	Expiry    *DCPExpiry     `xml:"expiry,omitempty"`
}

type Access struct {
	All              *struct{} `xml:"all,omitempty"`
	None             *struct{} `xml:"none,omitempty"`
	Null             *struct{} `xml:"null,omitempty"`
	Personal         *struct{} `xml:"personal,omitempty"`
	PersonalAndOther *struct{} `xml:"personalAndOther,omitempty"`
	Other            *struct{} `xml:"other,omitempty"`
}

type DCPStatement struct {
	Purpose   DCPPurpose   `xml:"purpose"`
	Recipient DCPRecipient `xml:"recipient"`
	Retention DCPRetention `xml:"retention"`
}

type DCPPurpose struct {
	Admin   *struct{} `xml:"admin,omitempty"`
	Contact *struct{} `xml:"contact,omitempty"`
	Prov    *struct{} `xml:"prov,omitempty"`
	Other   *struct{} `xml:"other,omitempty"`
}

type DCPRecipient struct {
	Other     *struct{} `xml:"other,omitempty"`
	Ours      *Ours     `xml:"ours,omitempty"`
	Public    *struct{} `xml:"public,omitempty"`
	Same      *struct{} `xml:"same,omitempty"`
	Unrelated *struct{} `xml:"unrelated,omitempty"`
}

type Ours struct {
	RecDesc *string `xml:"recDesc,omitempty"`
}

type DCPRetention struct {
	Business   *struct{} `xml:"business,omitempty"`
	Indefinite *struct{} `xml:"indefinite,omitempty"`
	Legal      *struct{} `xml:"legal,omitempty"`
	None       *struct{} `xml:"none,omitempty"`
	Stated     *struct{} `xml:"stated,omitempty"`
}

type DCPExpiry struct {
	Absolute *AbsoluteExpiry `xml:"absolute,omitempty"`
	Relative *RelativeExpiry `xml:"relative,omitempty"`
}

type AbsoluteExpiry struct {
	DateTime string `xml:",chardata"` // ISO 8601
}

type RelativeExpiry struct {
	Duration string `xml:",chardata"` // ISO 8601
}
