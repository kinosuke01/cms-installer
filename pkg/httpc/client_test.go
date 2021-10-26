package httpc

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestClient_NewRequest(t *testing.T) {
	tt := []struct {
		name string

		config Config
		opts   RequestOptions

		expectedScheme   string
		expectedHost     string
		expectedPath     string
		expectedRawQuery string
		expectedMethod   string
		expectedHeaders  map[string]string
		expectedBody     string
		expectedError    string
	}{
		{
			name: "empty_options",
			config: Config{
				Scheme: "https",
				Host:   "localhost",
			},
			opts:             RequestOptions{},
			expectedScheme:   "https",
			expectedHost:     "localhost",
			expectedPath:     "",
			expectedRawQuery: "",
			expectedMethod:   http.MethodGet,
			expectedHeaders:  map[string]string{},
			expectedBody:     "",
			expectedError:    "",
		},
		{
			name: "get_request",
			config: Config{
				Scheme:   "https",
				Host:     "localhost",
				BasePath: "api",
				BaseHeaders: map[string]string{
					"User-Agent": "httpc",
				},
			},
			opts: RequestOptions{
				Method: http.MethodGet,
				Path:   "articles",
				Queries: map[string]string{
					"cat":     "wordpress",
					"keyword": "blog",
				},
			},
			expectedScheme:   "https",
			expectedHost:     "localhost",
			expectedPath:     "/api/articles",
			expectedRawQuery: "cat=wordpress&keyword=blog",
			expectedMethod:   http.MethodGet,
			expectedHeaders: map[string]string{
				"User-Agent": "httpc",
			},
			expectedBody:  "",
			expectedError: "",
		},
		{
			name: "post_request",
			config: Config{
				Scheme:   "https",
				Host:     "localhost",
				BasePath: "api",
				BaseHeaders: map[string]string{
					"User-Agent": "httpc",
				},
			},
			opts: RequestOptions{
				Method: http.MethodPost,
				Path:   "articles",
				Headers: map[string]string{
					"Authorization": "Bearer xxxxx",
				},
				BodyValues: url.Values{
					"title":   []string{"Today is nice weather"},
					"content": []string{"Hello world. Today is nice weather. So I feel good."},
				},
			},
			expectedScheme:   "https",
			expectedHost:     "localhost",
			expectedPath:     "/api/articles",
			expectedRawQuery: "",
			expectedMethod:   http.MethodPost,
			expectedHeaders: map[string]string{
				"User-Agent":    "httpc",
				"Authorization": "Bearer xxxxx",
			},
			expectedBody:  "content=Hello+world.+Today+is+nice+weather.+So+I+feel+good.&title=Today+is+nice+weather",
			expectedError: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := New(&tc.config)
			req, err := c.NewRequest(context.Background(), &tc.opts)

			if err != nil {
				if tc.expectedError != err.Error() {
					t.Fatalf("error message wrong. want=%+v, got=%+v", tc.expectedError, err.Error())
				}
				return
			}

			if tc.expectedError != "" {
				t.Fatalf("error message wrong. want=%+v, got=%+v", tc.expectedError, "")
				return
			}

			if tc.expectedScheme != req.URL.Scheme {
				t.Fatalf("scheme wrong. want=%+v, got=%+v", tc.expectedScheme, req.URL.Scheme)
			}

			if tc.expectedHost != req.URL.Host {
				t.Fatalf("host wrong. want=%+v, got=%+v", tc.expectedHost, req.URL.Host)
			}

			if tc.expectedPath != req.URL.Path {
				t.Fatalf("path wrong. want=%+v, got=%+v", tc.expectedPath, req.URL.Path)
			}

			if tc.expectedRawQuery != req.URL.RawQuery {
				t.Fatalf("query wrong. want=%+v, got=%+v", tc.expectedRawQuery, req.URL.RawQuery)
			}

			if tc.expectedMethod != req.Method {
				t.Fatalf("query wrong. want=%+v, got=%+v", tc.expectedMethod, req.Method)
			}

			for k, v := range tc.expectedHeaders {
				val, ok := req.Header[k]
				if !ok {
					t.Fatalf("header key not exists. want=%+v", k)
				} else if v != strings.Join(val, ",") {
					t.Fatalf("header value wrong. key=%+v, want=%+v, got=%+v", k, v, val[0])
				}
			}

			bodyBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Fatalf(err.Error())
			} else if tc.expectedBody != string(bodyBytes) {
				t.Fatalf("body wrong. want=%+v, got=%+v", tc.expectedBody, string(bodyBytes))
			}
		})
	}
}

func TestClient_DoRequest(t *testing.T) {
	tt := []struct {
		name string

		opts RequestOptions

		expectedReqPath       string
		expectedReqMethod     string
		expectedReqHeaders    map[string]string
		expectedReqForm       map[string]string
		expectedResStatusCode int
		expectedResBody       string
	}{
		{
			name: "post_request",
			opts: RequestOptions{
				Method: http.MethodPost,
				Path:   "api/articles",
				Queries: map[string]string{
					"id": "12345",
				},
				Headers: map[string]string{
					"User-Agent":    "httpc",
					"Authorization": "Bearer xxxxx",
				},
				BodyValues: url.Values{
					"title":   []string{"Today is nice weather"},
					"content": []string{"Hello world. Today is nice weather. So I feel good."},
				},
			},
			expectedReqPath:   "/api/articles",
			expectedReqMethod: http.MethodPost,
			expectedReqHeaders: map[string]string{
				"User-Agent":    "httpc",
				"Authorization": "Bearer xxxxx",
			},
			expectedReqForm: map[string]string{
				"id":      "12345",
				"title":   "Today is nice weather",
				"content": "Hello world. Today is nice weather. So I feel good.",
			},
			expectedResStatusCode: 201,
			expectedResBody:       "result:true",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(
					http.HandlerFunc(
						func(w http.ResponseWriter, req *http.Request) {
							if tc.expectedReqPath != req.URL.Path {
								t.Fatalf("URL wrong. want=%+v, got=%+v", tc.expectedReqPath, req.URL.Path)
							}
							if tc.expectedReqMethod != req.Method {
								t.Fatalf("method wrong. want=%+v, got=%+v", tc.expectedReqMethod, req.Method)
							}

							for k, v := range tc.expectedReqHeaders {
								val, ok := req.Header[k]
								if !ok {
									t.Fatalf("header key not exists. want=%+v", k)
								} else if v != strings.Join(val, ",") {
									t.Fatalf("header value wrong. key=%+v, want=%+v, got=%+v", k, v, val[0])
								}
							}

							for k, v := range tc.expectedReqForm {
								val := req.FormValue(k)
								if v != val {
									t.Fatalf("form value wrong. key=%+v, want=%+v, got=%+v", k, v, val)
								}
							}

							w.WriteHeader(tc.expectedResStatusCode)
							w.Write([]byte(tc.expectedResBody))
						},
					),
				),
			)
			defer server.Close()

			serverURL, err := url.Parse(server.URL)
			if err != nil {
				t.Fatalf("failed to get mock server URL: %s", err.Error())
			}

			c := &Client{
				Config: Config{
					Scheme: serverURL.Scheme,
					Host:   serverURL.Host,
				},
				client: server.Client(),
			}
			res, err := c.DoRequest(context.Background(), &tc.opts)

			if err != nil {
				t.Fatalf("response error should not be nil. got=%+v", err.Error())
			}

			if tc.expectedResStatusCode != res.StatusCode {
				t.Fatalf("path wrong. want=%+v, got=%+v", tc.expectedResStatusCode, res.StatusCode)
			}
			if tc.expectedResBody != string(res.BodyBytes) {
				t.Fatalf("path wrong. want=%+v, got=%+v", tc.expectedResBody, string(res.BodyBytes))
			}
		})
	}
}
