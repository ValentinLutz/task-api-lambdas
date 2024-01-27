package v1_put_task

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
		"put task", func(t *testing.T) {
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

			requestBody := testfunctional.ReadFile("./put_task_request_body.json")

			// when
			request, err := http.NewRequest(
				http.MethodPut,
				fmt.Sprintf("%s/%s", testfunctional.BaseUrl, postResponseBody.TaskId.String()),
				requestBody,
			)
			if err != nil {
				t.Fatal(err)
			}
			request.Header.Set("Content-Type", "application/json")
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()

			// then
			td.Cmp(t, response.StatusCode, 204)
		},
	)

	t.Run(
		"put task that does not exist", func(t *testing.T) {
			// given
			taskId := "00000000-0000-0000-0000-000000000000"
			requestBody := testfunctional.ReadFile("./put_task_request_body_not_found.json.json")

			// when
			request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", testfunctional.BaseUrl, taskId), requestBody)
			if err != nil {
				t.Fatal(err)
			}
			request.Header.Set("Content-Type", "application/json")
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()

			// then
			td.Cmp(t, response.StatusCode, 404)
		},
	)
}
