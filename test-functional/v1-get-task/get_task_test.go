package v1_get_task

import (
	"encoding/json"
	"fmt"
	"net/http"
	testfunctional "root/test-functional"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func Test(t *testing.T) {
	t.Run(
		"get task", func(t *testing.T) {
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
			resp, err := http.Get(fmt.Sprintf("%s/%s", testfunctional.BaseUrl, postResponseBody.TaskId.String()))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// then
			td.Cmp(t, resp.StatusCode, 200)

			var respBody testfunctional.TaskResponse
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			td.CmpNoError(t, err)

			td.CmpJSON(
				t, respBody, "./get_task_response_body.json", []any{
					td.String(postResponseBody.TaskId.String()),
				},
			)
		},
	)

	t.Run(
		"get task that does not exist", func(t *testing.T) {
			// given
			taskId := "00000000-0000-0000-0000-000000000000"

			// when
			resp, err := http.Get(fmt.Sprintf("%s/%s", testfunctional.BaseUrl, taskId))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// then
			td.Cmp(t, resp.StatusCode, 404)
		},
	)
}
