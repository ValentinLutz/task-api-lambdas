package v1_delete_task

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
		"delete task", func(t *testing.T) {
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
			req, err := http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("%s/%s", testfunctional.BaseUrl, postResponseBody.TaskId.String()),
				nil,
			)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// then
			td.Cmp(t, resp.StatusCode, 204)
		},
	)

	t.Run(
		"delete task that does not exist", func(t *testing.T) {
			// given
			taskId := "00000000-0000-0000-0000-000000000000"

			// when
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", testfunctional.BaseUrl, taskId), nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// then
			td.Cmp(t, resp.StatusCode, 404)
		},
	)
}
