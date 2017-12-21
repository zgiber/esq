package esq

import (
	"encoding/json"
	"testing"

	"github.com/onsi/gomega"
)

func Test_term_search(t *testing.T) {
	gomega.RegisterTestingT(t)
	r := NewRequest(
		NewQuery().Should(
			Term("some value", "myKey.keyword"),
			Range(0, 10, "myNumberKey"),
		))

	q, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte(`
		{
			"query": {
				"bool": {
					"must": [
						{
							"term": {
								"myKey.keyword": "some value"
							}
						}
					]
				}
			}
		}`)

	gomega.Expect(q).Should(gomega.MatchJSON(expected))

}
