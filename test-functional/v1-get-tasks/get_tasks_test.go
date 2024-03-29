package v1_get_tasks

import (
	"encoding/json"
	"net/http"
	testfunctional "root/test-functional"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func Test(t *testing.T) {
	t.Run(
		"get tasks", func(t *testing.T) {
			// given
			postResponse, err := http.Post(
				testfunctional.BaseUrl,
				"application/json",
				testfunctional.ReadFile("./post_tasks_request_body.json"),
			)
			if err != nil {
				t.Fatal(err)
			}
			var postResponseBody testfunctional.TaskResponse
			err = json.NewDecoder(postResponse.Body).Decode(&postResponseBody)
			if err != nil {
				t.Fatal(err)
			}

			// when
			resp, err := http.Get(testfunctional.BaseUrl)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// then
			td.Cmp(t, resp.StatusCode, 200)

			var respBody testfunctional.TasksResponse
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			td.CmpNoError(t, err)

			td.CmpContains(
				t, respBody, postResponseBody, []any{},
			)
		},
	)
}
