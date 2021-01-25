package main_test

import (
	"bytes"
	"net/http"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main=>", func() {
	Describe("Creating app to expose services=>", func() {
		Context("Returns Created app", func() {
			app, err := createAPI()
			Expect(app).NotTo(Equal(nil))
			Expect(err).To(BeNil())
		})

	})

	Describe("Hit secret scanner API=>", func() {
		Context("Update method=>", func() {
			Context("Returns 200 as reponse, everything is ok", func() {
				app := fiber.New()

				controllerObjectUpdateSecretFile = func(data *updateParams) (int, string) {
					return 200, "Success"
				}

				app.Post("/api/detectsecrets/update", updateSecretFile)

				content := `
					{
						"repo":"secret-scanner",
						"owner":"felipe-hernandez-globant",
						"changes":{
							"vars/aws_horizontal.yml": [
						{
							"hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
							"is_verified": true,
									"is_secret": false,
							"line_number": 9,
							"type": "Secret Keyword"
						}
								],
							"vars/aws_sdc.yml": [
						{
							"hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
							"is_verified": true,
									"is_secret": true,
							"line_number": 9,
							"type": "Secret Keyword"
						}, {}
								]
						}
					}
				`
				req, _ := http.NewRequest("POST", "/api/detectsecrets/update", bytes.NewReader([]byte(content)))
				req.Header.Set("Content-Type", "application/json")

				resp, _ := app.Test(req)
				Expect(resp.StatusCode).To(Equal(200))

			})

			Context("Returns 400 as reponse, something went wrong", func() {
				app := fiber.New()

				controllerObjectUpdateSecretFile = func(data *updateParams) (int, string) {
					return 400, "Error"
				}

				app.Post("/api/detectsecrets/update", updateSecretFile)

				content := `
					{
						"repo":"secret-scanner",
						"owner":"felipe-hernandez-globant",
						"changes":{
							"vars/aws_horizontal.yml": [
						{
							"hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
							"is_verified": true,
									"is_secret": false,
							"line_number": 9,
							"type": "Secret Keyword"
						}
								],
							"vars/aws_sdc.yml": [
						{
							"hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
							"is_verified": true,
									"is_secret": true,
							"line_number": 9,
							"type": "Secret Keyword"
						}
								]
						}
					}
				`
				req, _ := http.NewRequest("POST", "/api/detectsecrets/update", bytes.NewReader([]byte(content)))
				req.Header.Set("Content-Type", "application/json")

				resp, _ := app.Test(req)
				Expect(resp.StatusCode).To(Equal(400))

			})

			Context("Error parsing request json structure", func() {
				app := fiber.New()

				app.Post("/api/detectsecrets/update", updateSecretFile)

				content := `{"key": "wrong request json structure"}`
				req, _ := http.NewRequest("POST", "/api/detectsecrets/update", bytes.NewReader([]byte(content)))

				resp, _ := app.Test(req)
				Expect(resp.StatusCode).To(Equal(400))

			})
		})

		Context("Create method=>", func() {
			Context("Returns 200 as reponse, everything is ok", func() {
				app := fiber.New()

				controllerObjectCreateSecretFile = func(data *createParams) (int, string) {
					return 200, "Success"
				}

				app.Post("/api/detectsecrets/create", createSecretFile)

				content := `
				{
					"repo": "secret-scanner",
					"owner": "felipe-hernandez-globant",
					"Content":"{\r\n  \"custom_plugin_paths\": [],\r\n  \"exclude\": {\r\n    \"files\": null,\r\n    \"lines\": null\r\n  },\r\n  \"generated_at\": \"2020-12-24T19:40:49Z\",\r\n  \"plugins_used\": [\r\n    {\r\n      \"name\": \"AWSKeyDetector\"\r\n    },\r\n    {\r\n      \"name\": \"ArtifactoryDetector\"\r\n    },\r\n    {\r\n      \"base64_limit\": \"4.5\",\r\n      \"name\": \"Base64HighEntropyString\"\r\n    },\r\n    {\r\n      \"name\": \"BasicAuthDetector\"\r\n    },\r\n    {\r\n      \"name\": \"CloudantDetector\"\r\n    },\r\n    {\r\n      \"hex_limit\": \"3\",\r\n      \"name\": \"HexHighEntropyString\"\r\n    },\r\n    {\r\n      \"name\": \"IbmCloudIamDetector\"\r\n    },\r\n    {\r\n      \"name\": \"IbmCosHmacDetector\"\r\n    },\r\n    {\r\n      \"name\": \"JwtTokenDetector\"\r\n    },\r\n    {\r\n      \"keyword_exclude\": null,\r\n      \"name\": \"KeywordDetector\"\r\n    },\r\n    {\r\n      \"name\": \"MailchimpDetector\"\r\n    },\r\n    {\r\n      \"name\": \"PrivateKeyDetector\"\r\n    },\r\n    {\r\n      \"name\": \"SlackDetector\"\r\n    },\r\n    {\r\n      \"name\": \"SoftlayerDetector\"\r\n    },\r\n    {\r\n      \"name\": \"StripeDetector\"\r\n    },\r\n    {\r\n      \"name\": \"TwilioKeyDetector\"\r\n    }\r\n  ],\r\n  \"results\": {},\r\n  \"version\": \"0.14.3\",\r\n  \"word_list\": {\r\n    \"file\": null,\r\n    \"hash\": null\r\n  }\r\n} "
				}
				`
				req, _ := http.NewRequest("POST", "/api/detectsecrets/create", bytes.NewReader([]byte(content)))
				req.Header.Set("Content-Type", "application/json")

				resp, _ := app.Test(req)
				Expect(resp.StatusCode).To(Equal(200))

			})

			Context("Returns 400 as reponse, something went wrong", func() {
				app := fiber.New()

				controllerObjectCreateSecretFile = func(data *createParams) (int, string) {
					return 400, "Error"
				}

				app.Post("/api/detectsecrets/create", createSecretFile)

				content := `
					{
						"repo":"secret-scanner",
						"owner":"felipe-hernandez-globant",
						"changes":{
							"vars/aws_horizontal.yml": [
						{
							"hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
							"is_verified": true,
									"is_secret": false,
							"line_number": 9,
							"type": "Secret Keyword"
						}
								],
							"vars/aws_sdc.yml": [
						{
							"hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
							"is_verified": true,
									"is_secret": true,
							"line_number": 9,
							"type": "Secret Keyword"
						}
								]
						}
					}
				`
				req, _ := http.NewRequest("POST", "/api/detectsecrets/create", bytes.NewReader([]byte(content)))
				req.Header.Set("Content-Type", "application/json")

				resp, _ := app.Test(req)
				Expect(resp.StatusCode).To(Equal(400))

			})

			Context("Error parsing request json structure", func() {
				app := fiber.New()

				app.Post("/api/detectsecrets/create", createSecretFile)

				content := `{"key": "wrong request json structure"}`
				req, _ := http.NewRequest("POST", "/api/detectsecrets/create", bytes.NewReader([]byte(content)))

				resp, _ := app.Test(req)
				Expect(resp.StatusCode).To(Equal(400))

			})
		})

	})
})
