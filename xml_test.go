package main

import "testing"

func TestConvertSelfClosingTags(t *testing.T) {
	// <all></all> incorrectly as a child of <greeting>
	// <ours> with and without <recDesc> tag
	// <same/> tag which is already self-closing
	// <hello></hello> tag which is not defined in the path set
	input := []byte(`<epp>
  <greeting>
    <all></all>
    <dcp>
      <access>
        <all></all>
		<hello></hello>
      </access>
      <statement>
        <recipient>
		  <ours></ours>
          <same/>
          <ours>
            <recDesc>message</recDesc>
          </ours>
        </recipient>
      </statement>
    </dcp>
  </greeting>
</epp>`)

	output := ConvertSelfClosingTags(input)

	expectedOutput := []byte(`<epp>
  <greeting>
    <all></all>
    <dcp>
      <access>
        <all/>
        <hello></hello>
      </access>
      <statement>
        <recipient>
          <ours/>
          <same/>
          <ours>
            <recDesc>message</recDesc>
          </ours>
        </recipient>
      </statement>
    </dcp>
  </greeting>
</epp>`)

	if string(output) != string(expectedOutput) {
		t.Errorf("Expected %s, but got %s", expectedOutput, output)
	}
}
