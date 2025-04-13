package dolores

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = FDescribe("Follow Operations", Ordered, func() {
	var db *pgxpool.Pool
	var followUser1Token, followUser2Token, preexistingFollowToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0023-follows-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(3) // 3 hub users to sign in

		hubSigninAsync(
			"follow-user1@0023-follow.example.com",
			"NewPassword123$",
			&followUser1Token,
			&wg,
		)
		hubSigninAsync(
			"follow-user2@0023-follow.example.com",
			"NewPassword123$",
			&followUser2Token,
			&wg,
		)
		hubSigninAsync(
			"preexisting@0023-follow.example.com",
			"NewPassword123$",
			&preexistingFollowToken,
			&wg,
		)
		wg.Wait()
	})

	AfterAll(func() {
		// Clean up the database using the down migration
		seedDatabase(db, "0023-follows-down.pgsql")
		db.Close()
	})

	Describe("Follow User", func() {
		type followUserTestCase struct {
			description string
			token       string
			request     hub.FollowUserRequest
			wantStatus  int
			verify      func()
		}

		It("should handle various follow user scenarios correctly", func() {
			testCases := []followUserTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.FollowUserRequest{
						Handle: "follow-user2",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.FollowUserRequest{
						Handle: "follow-user2",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "follow a user successfully",
					token:       followUser1Token,
					request: hub.FollowUserRequest{
						Handle: "follow-user2",
					},
					wantStatus: http.StatusOK,
					verify: func() {
						// Verify the follow relationship was created in the database
						var count int
						err := db.QueryRow(
							context.Background(),
							`SELECT COUNT(*) FROM following_relationships 
							 WHERE consuming_hub_user_id = '12345678-0023-0023-0023-000000000001' 
							 AND producing_hub_user_id = '12345678-0023-0023-0023-000000000002'`,
						).Scan(&count)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(count).Should(Equal(1))
					},
				},
				{
					description: "follow the same user again (idempotent)",
					token:       followUser1Token,
					request: hub.FollowUserRequest{
						Handle: "follow-user2",
					},
					wantStatus: http.StatusOK,
					verify: func() {
						// Verify there's still only one follow relationship
						var count int
						err := db.QueryRow(
							context.Background(),
							`SELECT COUNT(*) FROM following_relationships 
							 WHERE consuming_hub_user_id = '12345678-0023-0023-0023-000000000001' 
							 AND producing_hub_user_id = '12345678-0023-0023-0023-000000000002'`,
						).Scan(&count)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(count).Should(Equal(1))
					},
				},
				{
					description: "follow with non-existent handle",
					token:       followUser1Token,
					request: hub.FollowUserRequest{
						Handle: "non-existent-user",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "follow user with deleted account",
					token:       followUser1Token,
					request: hub.FollowUserRequest{
						Handle: "deleted-user",
					},
					wantStatus: http.StatusNotFound, // Assuming we can't follow deleted users
				},
				{
					description: "follow self (according to spec: returns 200 without creating records)",
					token:       followUser1Token,
					request: hub.FollowUserRequest{
						Handle: "follow-user1", // Same as the token user
					},
					wantStatus: http.StatusOK, // According to spec, following yourself returns 200 OK
					verify: func() {
						// Verify no self-follow relationship was created in the database
						var count int
						err := db.QueryRow(
							context.Background(),
							`SELECT COUNT(*) FROM following_relationships 
							 WHERE consuming_hub_user_id = '12345678-0023-0023-0023-000000000001' 
							 AND producing_hub_user_id = '12345678-0023-0023-0023-000000000001'`,
						).Scan(&count)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(
							count,
						).Should(Equal(0), "No self-follow relationship should be created")
					},
				},
				{
					description: "follow with empty handle",
					token:       followUser1Token,
					request: hub.FollowUserRequest{
						Handle: "",
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing FollowUser: %s\n",
					tc.description,
				)

				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/follow-user",
					tc.wantStatus,
				)

				if tc.verify != nil && tc.wantStatus == http.StatusOK {
					tc.verify()
				}

				// For debugging purposes
				if tc.wantStatus != http.StatusOK {
					fmt.Fprintf(
						GinkgoWriter,
						"Response: %v\n",
						resp,
					)
				}
			}
		})
	})

	Describe("Unfollow User", func() {
		type unfollowUserTestCase struct {
			description string
			token       string
			request     hub.UnfollowUserRequest
			wantStatus  int
			verify      func()
		}

		It("should handle various unfollow user scenarios correctly", func() {
			// First, set up a follow relationship for unfollow testing
			// User 2 follows User 1 for this specific test
			_, err := db.Exec(
				context.Background(),
				`INSERT INTO following_relationships (consuming_hub_user_id, producing_hub_user_id) 
				 VALUES ('12345678-0023-0023-0023-000000000002', '12345678-0023-0023-0023-000000000001')`,
			)
			Expect(err).ShouldNot(HaveOccurred())

			testCases := []unfollowUserTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.UnfollowUserRequest{
						Handle: "follow-user1",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.UnfollowUserRequest{
						Handle: "follow-user1",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "unfollow a user successfully",
					token:       followUser2Token,
					request: hub.UnfollowUserRequest{
						Handle: "follow-user1",
					},
					wantStatus: http.StatusOK,
					verify: func() {
						// Verify the follow relationship was removed from the database
						var count int
						err := db.QueryRow(
							context.Background(),
							`SELECT COUNT(*) FROM following_relationships 
							 WHERE consuming_hub_user_id = '12345678-0023-0023-0023-000000000002' 
							 AND producing_hub_user_id = '12345678-0023-0023-0023-000000000001'`,
						).Scan(&count)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(count).Should(Equal(0))
					},
				},
				{
					description: "unfollow the same user again (idempotent)",
					token:       followUser2Token,
					request: hub.UnfollowUserRequest{
						Handle: "follow-user1",
					},
					wantStatus: http.StatusOK,
					verify: func() {
						// Verify the follow relationship is still gone
						var count int
						err := db.QueryRow(
							context.Background(),
							`SELECT COUNT(*) FROM following_relationships 
							 WHERE consuming_hub_user_id = '12345678-0023-0023-0023-000000000002' 
							 AND producing_hub_user_id = '12345678-0023-0023-0023-000000000001'`,
						).Scan(&count)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(count).Should(Equal(0))
					},
				},
				{
					description: "unfollow with non-existent handle",
					token:       followUser2Token,
					request: hub.UnfollowUserRequest{
						Handle: "non-existent-user",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "unfollow user with deleted account",
					token:       followUser2Token,
					request: hub.UnfollowUserRequest{
						Handle: "deleted-user",
					},
					// Still 200 because we're just ensuring they're not followed, which they're not
					wantStatus: http.StatusOK,
				},
				{
					description: "unfollow self (according to spec: returns 404)",
					token:       followUser2Token,
					request: hub.UnfollowUserRequest{
						Handle: "follow-user2", // Same as the token user
					},
					wantStatus: http.StatusNotFound, // According to spec "If a user attempts to unfollow themselves, a 404 status is returned"
				},
				{
					description: "unfollow with empty handle",
					token:       followUser2Token,
					request: hub.UnfollowUserRequest{
						Handle: "",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "unfollow a user that was never followed",
					token:       followUser2Token,
					request: hub.UnfollowUserRequest{
						Handle: "followee-user", // User 2 never followed this user
					},
					wantStatus: http.StatusOK, // Idempotent operation
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing UnfollowUser: %s\n",
					tc.description,
				)

				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/unfollow-user",
					tc.wantStatus,
				)

				if tc.verify != nil && tc.wantStatus == http.StatusOK {
					tc.verify()
				}

				// For debugging purposes
				if tc.wantStatus != http.StatusOK {
					fmt.Fprintf(
						GinkgoWriter,
						"Response: %v\n",
						resp,
					)
				}
			}
		})
	})

	Describe("Get Follow Status", func() {
		type getFollowStatusTestCase struct {
			description string
			token       string
			request     hub.GetFollowStatusRequest
			wantStatus  int
			validate    func([]byte)
		}

		It(
			"should handle various get follow status scenarios correctly",
			func() {
				testCases := []getFollowStatusTestCase{
					{
						description: "without authentication",
						token:       "",
						request: hub.GetFollowStatusRequest{
							Handle: "follow-user2",
						},
						wantStatus: http.StatusUnauthorized,
					},
					{
						description: "with invalid token",
						token:       "invalid-token",
						request: hub.GetFollowStatusRequest{
							Handle: "follow-user2",
						},
						wantStatus: http.StatusUnauthorized,
					},
					{
						description: "get status of user not following",
						token:       followUser1Token,
						request: hub.GetFollowStatusRequest{
							Handle: "preexisting-follow",
						},
						wantStatus: http.StatusOK,
						validate: func(respBody []byte) {
							var status hub.FollowStatus
							err := json.Unmarshal(respBody, &status)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(status.IsFollowing).Should(BeFalse())
							Expect(status.IsBlocked).Should(BeFalse())
							Expect(status.CanFollow).Should(BeTrue())
						},
					},
					{
						description: "get status of preexisting follow relationship",
						token:       preexistingFollowToken,
						request: hub.GetFollowStatusRequest{
							Handle: "followee-user",
						},
						wantStatus: http.StatusOK,
						validate: func(respBody []byte) {
							var status hub.FollowStatus
							err := json.Unmarshal(respBody, &status)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(status.IsFollowing).Should(BeTrue())
							Expect(status.IsBlocked).Should(BeFalse())
							Expect(
								status.CanFollow,
							).Should(BeFalse(), "canFollow should be false when already following a user")
						},
					},
					{
						description: "get status with non-existent handle",
						token:       followUser1Token,
						request: hub.GetFollowStatusRequest{
							Handle: "non-existent-user",
						},
						wantStatus: http.StatusNotFound,
					},
					{
						description: "get status of deleted user",
						token:       followUser1Token,
						request: hub.GetFollowStatusRequest{
							Handle: "deleted-user",
						},
						wantStatus: http.StatusOK,
						validate: func(respBody []byte) {
							var status hub.FollowStatus
							err := json.Unmarshal(respBody, &status)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(status.IsFollowing).Should(BeFalse())
							Expect(status.IsBlocked).Should(BeFalse())
							// Deleted users cannot be followed
							Expect(status.CanFollow).Should(BeFalse())
						},
					},
					{
						description: "get status of self (according to spec: isFollowing=true, isBlocked=false, canFollow=false)",
						token:       followUser1Token,
						request: hub.GetFollowStatusRequest{
							Handle: "follow-user1", // Same as the token user
						},
						wantStatus: http.StatusOK,
						validate: func(respBody []byte) {
							var status hub.FollowStatus
							err := json.Unmarshal(respBody, &status)
							Expect(err).ShouldNot(HaveOccurred())
							// According to spec: "When checking one's own status: isFollowing=true, isBlocked=false, canFollow=false"
							Expect(
								status.IsFollowing,
							).Should(BeTrue(), "isFollowing should be true when checking self status")
							Expect(
								status.IsBlocked,
							).Should(BeFalse(), "isBlocked should be false when checking self status")
							Expect(
								status.CanFollow,
							).Should(BeFalse(), "canFollow should be false when checking self status")
						},
					},
					{
						description: "get status with empty handle",
						token:       followUser1Token,
						request: hub.GetFollowStatusRequest{
							Handle: "",
						},
						wantStatus: http.StatusBadRequest,
					},
				}

				for _, tc := range testCases {
					fmt.Fprintf(
						GinkgoWriter,
						"### Testing GetFollowStatus: %s\n",
						tc.description,
					)

					resp := testPOSTGetResp(
						tc.token,
						tc.request,
						"/hub/get-follow-status",
						tc.wantStatus,
					)

					if tc.validate != nil && tc.wantStatus == http.StatusOK {
						tc.validate(resp.([]byte))
					}

					// For debugging purposes
					if tc.wantStatus != http.StatusOK &&
						tc.wantStatus != http.StatusUnauthorized {
						fmt.Fprintf(
							GinkgoWriter,
							"Response: %v\n",
							resp,
						)
					}
				}
			},
		)
	})
})
