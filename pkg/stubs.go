package stubs

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ListenAndServe(cfg *Config) error {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
	)

	fmt.Printf("Running the stub server on port %d\n", cfg.Port)

	for _, service := range cfg.Services {
		for _, endpoint := range service.Endpoints {
			service := service
			endpoint := endpoint

			router.MethodFunc(endpoint.Method, "/"+service.Prefix+endpoint.Name, func(w http.ResponseWriter, r *http.Request) {
				log.Println("request:", endpoint.Method, "/"+service.Prefix+endpoint.Name)

				for k, v := range cfg.Header {
					w.Header().Set(k, v)
				}

				if len(endpoint.Matches) != 0 {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						log.Fatal(err)
						return
					}

					for _, match := range endpoint.Matches {
						reqBody, err := os.ReadFile(cfg.ResponseDir + "/" + match.RequestBody)
						if err != nil {
							http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
							log.Fatal(err)
							return
						}

						if bytes.Contains(body, reqBody) {
							for k, v := range match.Response.Header {
								w.Header().Set(k, v)
							}
							w.WriteHeader(match.Response.Status)

							responseBody, err := os.ReadFile(cfg.ResponseDir + "/" + match.Response.Body)
							if err != nil {
								http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
								log.Fatal(err)
								return
							}

							_, err = w.Write(responseBody)
							if err != nil {
								fmt.Println(err)
							}

							return
						}
					}
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				} else {
					for k, v := range endpoint.Response.Header {
						w.Header().Set(k, v)
					}
					w.WriteHeader(endpoint.Response.Status)

					body, err := os.ReadFile(cfg.ResponseDir + "/" + endpoint.Response.Body)
					if err != nil {
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						log.Fatal(err)
						return
					}

					_, err = w.Write(body)
					if err != nil {
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						log.Fatal(err)
					}
				}
			})
		}
	}
	return http.ListenAndServe(net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)), router)
}
