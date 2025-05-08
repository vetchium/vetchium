package db

import (
	"errors"
	"fmt"
)

// Ideally should be a const, but go doesn't support const errors.
var (
	ErrInternal               = errors.New("internal error")
	ErrNoEmployer             = errors.New("employer not found")
	ErrInviteTokenNotFound    = errors.New("invite token not found")
	ErrOrgUserAlreadyExists   = errors.New("org user already exists")
	ErrNoOrgUser              = errors.New("org user not found")
	ErrLastActiveAdmin        = errors.New("cannot disable last active admin")
	ErrOrgUserAlreadyDisabled = errors.New("org user already disabled")
	ErrOrgUserNotDisabled     = errors.New("org user not in disabled state")

	ErrDupCostCenterName = errors.New("duplicate cost center name")
	ErrNoCostCenter      = errors.New("cost center not found")

	ErrDupLocationName = errors.New("location name already exists")
	ErrNoLocation      = errors.New("location not found")

	ErrNoOpening       = errors.New("opening not found")
	ErrTooManyWatchers = errors.New("too many watchers")

	ErrNoRecruiter          = errors.New("recruiter not found")
	ErrNoHiringManager      = errors.New("hiring manager not found")
	ErrNoStateChangeWaiting = errors.New("no state change waiting")
	ErrInvalidRecruiter     = errors.New(
		"one or more invalid recruiter emails specified",
	)
	ErrInvalidHiringTeam = errors.New(
		"one or more invalid hiring team member emails specified or members not in active state",
	)

	ErrStateMismatch = errors.New(
		"current state does not match expected state",
	)
	ErrInvalidPasswordResetToken    = errors.New("invalid password reset token")
	ErrNoHubUser                    = errors.New("hub user not found")
	ErrDupHandle                    = errors.New("handle already in use")
	ErrBadResume                    = errors.New("bad resume")
	ErrNoApplication                = errors.New("application not found")
	ErrApplicationStateInCompatible = errors.New("state incompatible")
	ErrUnauthorizedComment          = errors.New(
		"user not authorized to comment on candidacy",
	)
	ErrInvalidCandidacyState = errors.New(
		"candidacy not in valid state for comments",
	)
	ErrNoInterview           = errors.New("interview not found")
	ErrInvalidInterviewState = errors.New("interview not in valid state")
	ErrNoCandidacy           = errors.New("candidacy not found")
	ErrInterviewerNotActive  = errors.New("interviewer is not in active state")
	ErrNotAnInterviewer      = errors.New(
		"user is not an interviewer for this interview",
	)
	ErrInvalidPaginationKey    = fmt.Errorf("invalid pagination key")
	ErrNoWorkHistory           = errors.New("work history not found")
	ErrDuplicateOfficialEmail  = errors.New("official email already exists")
	ErrTooManyOfficialEmails   = errors.New("too many official emails")
	ErrOfficialEmailNotFound   = errors.New("official email not found")
	ErrInvalidVerificationCode = errors.New("invalid verification code")
	ErrNotColleaguable         = errors.New(
		"cannot send colleague request to this user",
	)
	ErrNoConnection     = errors.New("no active colleague connection found")
	ErrNotColleague     = errors.New("one or more endorsers are not colleagues")
	ErrTooManyEndorsers = errors.New("too many endorsers specified")

	// ErrInviteNotNeeded is returned when the invite is not needed
	// because the user is already a hub user or the invite is already
	// sent recently or the user does not want to receive invites.
	// TODO: Implement support for users to block invites !?
	ErrInviteNotNeeded            = errors.New("invite not needed")
	ErrUserNotFound               = errors.New("user not found")
	ErrInviteNotFound             = errors.New("invite not found")
	ErrDomainNotApprovedForSignup = errors.New("domain not approved for signup")

	ErrNoEducation         = errors.New("education not found")
	ErrNoAchievement       = errors.New("achievement not found")
	ErrCannotApply         = errors.New("user cannot apply to this opening")
	ErrUnpaidHubUser       = errors.New("user is not a paid hub user")
	ErrNoPost              = errors.New("post not found")
	ErrNonVoteableUserPost = errors.New("user cannot vote for this post")
)
