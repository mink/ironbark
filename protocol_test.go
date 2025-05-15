package main

import (
	"encoding/xml"
	"fmt"
	"testing"
	"time"
)

func TestGreetingXML(t *testing.T) {
	date := time.Now().UTC().Format("2006-01-02T15:04:05.0Z")

	epp := EPP{
		Xmlns: "urn:ietf:params:xml:ns:epp-1.0",
		Greeting: Greeting{
			SvID:   "Ironbark EPP Server",
			SvDate: date,
			SvcMenu: ServiceMenu{
				Version: []string{"1.0"},
				Lang:    []string{"en"},
				ObjURI: []string{
					"urn:ietf:params:xml:ns:obj1",
					"urn:ietf:params:xml:ns:obj2",
					"urn:ietf:params:xml:ns:obj3",
				},
				SvcExtension: &SvcExtension{
					ExtURI: []string{"http://custom/obj1ext-1.0"},
				},
			},
			DCP: DCP{
				Access: Access{All: &struct{}{}},
				Statement: []DCPStatement{
					{
						Purpose: DCPPurpose{
							Admin: &struct{}{},
							Prov:  &struct{}{}},
						Recipient: DCPRecipient{
							Ours:   &Ours{},
							Public: &struct{}{},
						},
						Retention: DCPRetention{
							Stated: &struct{}{},
						},
					},
				},
				Expiry: nil,
			},
		},
	}

	bytes, err := xml.MarshalIndent(epp, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}
	greeting := string(bytes)

	// todo - self closing tags for DCP

	expectedGreeting := fmt.Sprintf(`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <greeting>
    <svID>Ironbark EPP Server</svID>
    <svDate>%s</svDate>
    <svcMenu>
      <version>1.0</version>
      <lang>en</lang>
      <objURI>urn:ietf:params:xml:ns:obj1</objURI>
      <objURI>urn:ietf:params:xml:ns:obj2</objURI>
      <objURI>urn:ietf:params:xml:ns:obj3</objURI>
      <svcExtension>
        <extURI>http://custom/obj1ext-1.0</extURI>
      </svcExtension>
    </svcMenu>
    <dcp>
      <access>
        <all></all>
      </access>
      <statement>
        <purpose>
          <admin></admin>
          <prov></prov>
        </purpose>
        <recipient>
          <ours></ours>
          <public></public>
        </recipient>
        <retention>
          <stated></stated>
        </retention>
      </statement>
    </dcp>
  </greeting>
</epp>`, date)

	if greeting != expectedGreeting {
		t.Errorf("Expected XML:\n%s\nGot:\n%s", expectedGreeting, greeting)
	}
}
