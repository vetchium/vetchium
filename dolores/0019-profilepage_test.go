package dolores

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/psankar/vetchi/typespec/hub"
)

var _ = FDescribe("Profile Page", Ordered, func() {
	var db *pgxpool.Pool
	var hubToken1, hubToken2 string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0019-profilepage-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(2)
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
			handle2 := "profilepage_user2" // Already used by another user
			newHandle := "new_unique_handle"
			updatedName := "Updated Name"
			updatedShortBio := "Updated short bio"
			updatedLongBio := "Updated long bio"
			invalidHandle := "inv@lid_handle"
			emptyStr := ""

			testCases := []updateBioTestCase{
				// {
				// 	description: "without authentication",
				// 	token:       "",
				// 	request: hub.UpdateBioRequest{
				// 		Handle:   &handle1,
				// 		FullName: &updatedName,
				// 		ShortBio: &updatedShortBio,
				// 		LongBio:  &updatedLongBio,
				// 	},
				// 	wantStatus: http.StatusUnauthorized,
				// },
				{
					description: "update bio with new unique handle",
					token:       hubToken1,
					request: hub.UpdateBioRequest{
						Handle:   &newHandle,
						FullName: &updatedName,
						ShortBio: &updatedShortBio,
						LongBio:  &updatedLongBio,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "update bio with duplicate handle",
					token:       hubToken1,
					request: hub.UpdateBioRequest{
						Handle:   &handle2, // Try to use handle that's already taken by user2
						FullName: &updatedName,
						ShortBio: &updatedShortBio,
						LongBio:  &updatedLongBio,
					},
					wantStatus: http.StatusConflict,
				},
				{
					description: "update bio without changing handle",
					token:       hubToken2,
					request: hub.UpdateBioRequest{
						Handle:   &handle2, // Keep same handle
						FullName: &updatedName,
						ShortBio: &updatedShortBio,
						LongBio:  &updatedLongBio,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "update with invalid handle format",
					token:       hubToken1,
					request: hub.UpdateBioRequest{
						Handle:   &invalidHandle,
						FullName: &updatedName,
						ShortBio: &updatedShortBio,
						LongBio:  &updatedLongBio,
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "update with empty required fields",
					token:       hubToken1,
					request: hub.UpdateBioRequest{
						Handle:   &handle1,
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
							Handle: *tc.request.Handle,
						},
						"/hub/get-bio",
						http.StatusOK,
					)

					var bio hub.Bio
					err := json.Unmarshal(resp.([]byte), &bio)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(bio.Handle).Should(Equal(*tc.request.Handle))
					Expect(bio.FullName).Should(Equal(*tc.request.FullName))
					Expect(bio.ShortBio).Should(Equal(*tc.request.ShortBio))
					Expect(bio.LongBio).Should(Equal(*tc.request.LongBio))
				}
			}
		})
	})

	Describe("Profile Picture Operations", func() {
		sampleImageBytes := []byte(
			"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==",
		)
		invalidImageBytes := []byte("invalid-base64")
		emptyImageBytes := []byte("")

		Describe("Upload Profile Picture", func() {
			type uploadProfilePictureTestCase struct {
				description string
				token       string
				request     hub.UploadProfilePictureRequest
				wantStatus  int
			}

			It("should handle various upload test cases correctly", func() {
				testCases := []uploadProfilePictureTestCase{
					{
						description: "without authentication",
						token:       "",
						request: hub.UploadProfilePictureRequest{
							Image: sampleImageBytes,
						},
						wantStatus: http.StatusUnauthorized,
					},
					{
						description: "upload valid image",
						token:       hubToken2,
						request: hub.UploadProfilePictureRequest{
							Image: sampleImageBytes,
						},
						wantStatus: http.StatusOK,
					},
					{
						description: "upload invalid base64 data",
						token:       hubToken2,
						request: hub.UploadProfilePictureRequest{
							Image: invalidImageBytes,
						},
						wantStatus: http.StatusBadRequest,
					},
					{
						description: "upload empty image",
						token:       hubToken2,
						request: hub.UploadProfilePictureRequest{
							Image: emptyImageBytes,
						},
						wantStatus: http.StatusBadRequest,
					},
				}

				for _, tc := range testCases {
					fmt.Fprintf(
						GinkgoWriter,
						"### Testing: %s\n",
						tc.description,
					)
					testPOST(
						tc.token,
						tc.request,
						"/hub/upload-profile-picture",
						tc.wantStatus,
					)
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
				req, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf(
						"%s/hub/profile-picture/profilepage_user1",
						serverURL,
					),
					nil,
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Authorization", "Bearer "+hubToken1)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))
				resp.Body.Close()

				// Upload a profile picture
				testPOST(
					hubToken1,
					hub.UploadProfilePictureRequest{
						Image: sampleImageBytes,
					},
					"/hub/upload-profile-picture",
					http.StatusOK,
				)

				// Get the uploaded picture
				req, err = http.NewRequest(
					http.MethodGet,
					fmt.Sprintf(
						"%s/hub/profile-picture/profilepage_user1",
						serverURL,
					),
					nil,
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Authorization", "Bearer "+hubToken1)

				resp, err = http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))
				body, err := io.ReadAll(resp.Body)
				Expect(err).ShouldNot(HaveOccurred())
				_, err = base64.StdEncoding.DecodeString(string(body))
				Expect(err).ShouldNot(HaveOccurred())
				resp.Body.Close()

				// Remove the profile picture
				req, err = http.NewRequest(
					http.MethodPost,
					serverURL+"/hub/remove-profile-picture",
					bytes.NewBuffer([]byte("{}")),
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Authorization", "Bearer "+hubToken1)
				req.Header.Set("Content-Type", "application/json")

				resp, err = http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))
				resp.Body.Close()

				// Verify picture is gone
				req, err = http.NewRequest(
					http.MethodGet,
					fmt.Sprintf(
						"%s/hub/profile-picture/profilepage_user1",
						serverURL,
					),
					nil,
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Authorization", "Bearer "+hubToken1)

				resp, err = http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))
				resp.Body.Close()
			})
		})
	})
})
