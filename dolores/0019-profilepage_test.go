package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = Describe("Profile Page", Ordered, func() {
	var db *pgxpool.Pool
	var hubToken1, hubToken2, hubToken3 string
	var sampleImageBytes []byte

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0019-profilepage-up.pgsql")

		// Read test image file
		var err error
		sampleImageBytes, err = os.ReadFile("avatar1.jpg")
		Expect(err).ShouldNot(HaveOccurred())

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(3)
		hubSigninAsync(
			"user1@profilepage-hub.example",
			"NewPassword123$",
			&hubToken1,
			&wg,
		)
		hubSigninAsync(
			"user2@profilepage-hub.example",
			"NewPassword123$",
			&hubToken2,
			&wg,
		)
		hubSigninAsync(
			"user3@profilepage-hub.example",
			"NewPassword123$",
			&hubToken3,
			&wg,
		)
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0019-profilepage-down.pgsql")
		db.Close()
	})

	Describe("Get Bio", func() {
		type getBioTestCase struct {
			description string
			token       string
			request     hub.GetBioRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various test cases correctly", func() {
			testCases := []getBioTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.GetBioRequest{
						Handle: "profilepage_user1",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.GetBioRequest{
						Handle: "profilepage_user1",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "get own bio",
					token:       hubToken1,
					request: hub.GetBioRequest{
						Handle: "profilepage_user1",
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var bio hub.Bio
						err := json.Unmarshal(resp, &bio)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(bio.Handle).Should(Equal("profilepage_user1"))
						Expect(
							bio.FullName,
						).Should(Equal("Profile Test User 1"))
						Expect(
							bio.ShortBio,
						).Should(Equal("Profile Test User 1 is experienced"))
						Expect(
							bio.LongBio,
						).Should(Equal("Profile Test User 1 was born in India and finished education at IIT Mumbai."))
					},
				},
				{
					description: "get another user's bio",
					token:       hubToken1,
					request: hub.GetBioRequest{
						Handle: "profilepage_user2",
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var bio hub.Bio
						err := json.Unmarshal(resp, &bio)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(bio.Handle).Should(Equal("profilepage_user2"))
						Expect(
							bio.FullName,
						).Should(Equal("Profile Test User 2"))
					},
				},
				{
					description: "get non-existent user's bio",
					token:       hubToken1,
					request: hub.GetBioRequest{
						Handle: "nonexistent_user",
					},
					wantStatus: http.StatusNotFound,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/get-bio",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("Update Bio", func() {
		type updateBioTestCase struct {
			description string
			token       string
			request     hub.UpdateBioRequest
			wantStatus  int
		}

		It("should handle various test cases correctly", func() {
			// Test data
			handle1 := "profilepage_user1"
			updatedName := "Updated Name"
			updatedShortBio := "Updated short bio"
			updatedLongBio := "Updated long bio"
			emptyStr := ""

			testCases := []updateBioTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.UpdateBioRequest{
						FullName: &updatedName,
						ShortBio: &updatedShortBio,
						LongBio:  &updatedLongBio,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "update bio with valid data",
					token:       hubToken1,
					request: hub.UpdateBioRequest{
						FullName: &updatedName,
						ShortBio: &updatedShortBio,
						LongBio:  &updatedLongBio,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "update with empty required fields",
					token:       hubToken1,
					request: hub.UpdateBioRequest{
						FullName: &emptyStr,
						ShortBio: &emptyStr,
						LongBio:  &emptyStr,
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/hub/update-bio",
					tc.wantStatus,
				)

				// If it was a successful update, verify the changes
				if tc.wantStatus == http.StatusOK {
					// Get the updated bio
					resp := testPOSTGetResp(
						tc.token,
						hub.GetBioRequest{
							Handle: handle1,
						},
						"/hub/get-bio",
						http.StatusOK,
					)

					var bio hub.Bio
					err := json.Unmarshal(resp.([]byte), &bio)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(bio.Handle).Should(Equal(handle1))
					Expect(bio.FullName).Should(Equal(*tc.request.FullName))
					Expect(bio.ShortBio).Should(Equal(*tc.request.ShortBio))
					Expect(bio.LongBio).Should(Equal(*tc.request.LongBio))
				}
			}
		})
	})

	Describe("Profile Picture Operations", func() {
		invalidImageBytes := []byte("invalid-base64")
		emptyImageBytes := []byte("")

		Describe("Upload Profile Picture", func() {
			type uploadProfilePictureTestCase struct {
				description string
				token       string
				imageBytes  []byte
				filename    string
				wantStatus  int
			}

			It("should handle various upload test cases correctly", func() {
				testCases := []uploadProfilePictureTestCase{
					{
						description: "without authentication",
						token:       "",
						imageBytes:  sampleImageBytes,
						filename:    "avatar1.jpg",
						wantStatus:  http.StatusUnauthorized,
					},
					{
						description: "upload valid image",
						token:       hubToken2,
						imageBytes:  sampleImageBytes,
						filename:    "avatar1.jpg",
						wantStatus:  http.StatusOK,
					},
					{
						description: "upload invalid image data",
						token:       hubToken2,
						imageBytes:  invalidImageBytes,
						filename:    "invalid.jpg",
						wantStatus:  http.StatusBadRequest,
					},
					{
						description: "upload empty image",
						token:       hubToken2,
						imageBytes:  emptyImageBytes,
						filename:    "empty.jpg",
						wantStatus:  http.StatusBadRequest,
					},
				}

				for _, tc := range testCases {
					fmt.Fprintf(
						GinkgoWriter,
						"### Testing: %s\n",
						tc.description,
					)

					// Create multipart form data
					body := &bytes.Buffer{}
					writer := multipart.NewWriter(body)

					// Create form file part with custom header for image/jpeg
					header := make(textproto.MIMEHeader)
					header.Set(
						"Content-Disposition",
						fmt.Sprintf(
							`form-data; name="image"; filename="%s"`,
							tc.filename,
						),
					)
					header.Set("Content-Type", "image/jpeg")
					var part io.Writer
					part, err := writer.CreatePart(header)
					Expect(err).ShouldNot(HaveOccurred())
					_, err = io.Copy(part, bytes.NewReader(tc.imageBytes))
					Expect(err).ShouldNot(HaveOccurred())

					// Close the multipart writer
					err = writer.Close()
					Expect(err).ShouldNot(HaveOccurred())

					// Create request
					req, err := http.NewRequest(
						http.MethodPost,
						serverURL+"/hub/upload-profile-picture",
						body,
					)
					Expect(err).ShouldNot(HaveOccurred())

					// Set headers - only Authorization and Content-Type with boundary
					if tc.token != "" {
						req.Header.Set("Authorization", "Bearer "+tc.token)
					}
					req.Header.Set("Content-Type", writer.FormDataContentType())

					// Send request
					resp, err := http.DefaultClient.Do(req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.StatusCode).Should(Equal(tc.wantStatus))
					resp.Body.Close()
				}
			})
		})

		Describe("Get Profile Picture", func() {
			type getProfilePictureTestCase struct {
				description string
				token       string
				handle      string
				wantStatus  int
			}

			It(
				"should handle various get picture test cases correctly",
				func() {
					testCases := []getProfilePictureTestCase{
						{
							description: "without authentication",
							token:       "",
							handle:      "profilepage_user1",
							wantStatus:  http.StatusUnauthorized,
						},
						{
							description: "get non-existent profile picture",
							token:       hubToken1,
							handle:      "profilepage_user1",
							wantStatus:  http.StatusNotFound,
						},
						{
							description: "get non-existent user's picture",
							token:       hubToken1,
							handle:      "nonexistent_user",
							wantStatus:  http.StatusNotFound,
						},
					}

					for _, tc := range testCases {
						fmt.Fprintf(
							GinkgoWriter,
							"### Testing: %s\n",
							tc.description,
						)
						req, err := http.NewRequest(
							http.MethodGet,
							fmt.Sprintf(
								"%s/hub/profile-picture/%s",
								serverURL,
								tc.handle,
							),
							nil,
						)
						Expect(err).ShouldNot(HaveOccurred())
						if tc.token != "" {
							req.Header.Set("Authorization", "Bearer "+tc.token)
						}

						resp, err := http.DefaultClient.Do(req)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.StatusCode).Should(Equal(tc.wantStatus))
						resp.Body.Close()
					}
				},
			)
		})

		Describe("End-to-End Profile Picture Flow", func() {
			It("should handle the complete profile picture lifecycle", func() {
				// First try to get non-existent picture
				fmt.Fprintf(
					GinkgoWriter,
					"get non-existent picture for profilepage_user3\n",
				)
				req, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf(
						"%s/hub/profile-picture/profilepage_user3?t=%d",
						serverURL,
						time.Now().UnixNano(),
					),
					nil,
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Authorization", "Bearer "+hubToken3)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))
				resp.Body.Close()

				// Upload a profile picture using multipart form
				fmt.Fprintf(
					GinkgoWriter,
					"upload profile picture for profilepage_user3\n",
				)
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				// Create form file part with custom header for image/jpeg
				header := make(textproto.MIMEHeader)
				header.Set(
					"Content-Disposition",
					fmt.Sprintf(
						`form-data; name="image"; filename="%s"`,
						"avatar1.jpg",
					),
				)
				header.Set("Content-Type", "image/jpeg")
				var part io.Writer
				part, err = writer.CreatePart(header)
				Expect(err).ShouldNot(HaveOccurred())
				_, err = io.Copy(part, bytes.NewReader(sampleImageBytes))
				Expect(err).ShouldNot(HaveOccurred())

				// Close the multipart writer
				err = writer.Close()
				Expect(err).ShouldNot(HaveOccurred())

				req, err = http.NewRequest(
					http.MethodPost,
					serverURL+"/hub/upload-profile-picture",
					body,
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Content-Type", writer.FormDataContentType())
				req.Header.Set("Authorization", "Bearer "+hubToken3)

				resp, err = http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))
				resp.Body.Close()

				// Get the uploaded picture
				fmt.Fprintf(
					GinkgoWriter,
					"get uploaded picture for profilepage_user3\n",
				)
				req, err = http.NewRequest(
					http.MethodGet,
					fmt.Sprintf(
						"%s/hub/profile-picture/profilepage_user3?t=%d",
						serverURL,
						time.Now().UnixNano(),
					),
					nil,
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Authorization", "Bearer "+hubToken3)

				resp, err = http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))
				responseBody, err := io.ReadAll(resp.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(responseBody)).Should(BeNumerically(">", 0))
				Expect(
					resp.Header.Get("Content-Type"),
				).Should(Equal("image/jpeg"))
				resp.Body.Close()

				// Remove the profile picture
				fmt.Fprintf(
					GinkgoWriter,
					"remove profile picture for profilepage_user3\n",
				)
				req, err = http.NewRequest(
					http.MethodPost,
					serverURL+"/hub/remove-profile-picture",
					bytes.NewBuffer([]byte("{}")),
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Authorization", "Bearer "+hubToken3)
				req.Header.Set("Content-Type", "application/json")

				resp, err = http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))
				resp.Body.Close()

				// Verify picture is gone
				fmt.Fprintf(
					GinkgoWriter,
					"verify picture is gone for profilepage_user3\n",
				)
				req, err = http.NewRequest(
					http.MethodGet,
					fmt.Sprintf(
						"%s/hub/profile-picture/profilepage_user3?t=%d",
						serverURL,
						time.Now().UnixNano(),
					),
					nil,
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Authorization", "Bearer "+hubToken3)

				resp, err = http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))
				resp.Body.Close()
			})
		})
	})
})
