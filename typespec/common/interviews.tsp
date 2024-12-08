union InterviewState {
    Scheduled: "SCHEDULED",
    Completed: "COMPLETED",
    CandidateWithdrew: "CANDIDATE_WITHDREW",
    EmployerWithdrew: "EMPLOYER_WITHDREW",
}

union InterviewersDecision {
    StrongYes: "STRONG_YES",
    Yes: "YES",
    No: "NO",
    StrongNo: "STRONG_NO",
}
