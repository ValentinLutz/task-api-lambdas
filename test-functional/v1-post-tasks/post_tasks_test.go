package v1_post_tasks

import (
	"encoding/json"
	"net/http"
	testfunctional "root/test-functional"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func Test(t *testing.T) {
	t.Run(
		"post tasks", func(t *testing.T) {
			// given
			reqBody := testfunctional.ReadFile("./post_tasks_request_body.json")

			// when
			resp, err := http.Post(testfunctional.BaseUrl, "application/json", reqBody)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// then
			td.Cmp(t, resp.StatusCode, 200)

			var respBody testfunctional.TasksResponse
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			td.CmpNoError(t, err)

			td.CmpJSON(
				t, respBody, "./post_tasks_response_body.json", []any{},
			)
		},
	)
}
