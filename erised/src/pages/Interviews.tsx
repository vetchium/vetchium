import {
  Button,
  Calendar,
  Flex,
  Radio,
  RadioChangeEvent,
  Table,
  Tag,
} from "antd";
import type { Dayjs } from "dayjs";
import { useState } from "react";
import { Link } from "react-router-dom";
import t from "../i18n/i18n";

const interviewData = [
  {
    interviewId: "INT-001",
    interviewStatus: "SCHEDULED",
    openingIds: ["JAN14-1", "JAN14-2"],
    candidacyId: "CAND-001",
    date: "2024-03-14T10:00:00Z",
    durationInMins: 60,
    candidateName: "Jeff Dean",
    candidateCurrentCompany: "Google",
    interviewers: [{ name: "Larry Page", email: "larrypage@example.com" }],
    hiringManager: { name: "Urs Holzle", email: "ursholzle@example.com" },
  },
  {
    interviewId: "INT-002",
    interviewStatus: "COMPLETED",
    evaluationStatus: "EVALUATION_PENDING",
    openingIds: ["JAN14-1", "JAN14-2"],
    candidacyId: "CAND-001",
    date: "2024-03-13T10:00:00Z",
    durationInMins: 30,
    candidateName: "Jeff Dean",
    candidateCurrentCompany: "Google",
    interviewers: [{ name: "Sergey Brin", email: "sergeybrin@example.com" }],
    hiringManager: { name: "Urs Holzle", email: "ursholzle@example.com" },
  },
  {
    interviewId: "INT-003",
    interviewStatus: "COMPLETED",
    evaluationStatus: "EVALUATION_COMPLETED",
    evaluationResult: "STRONG_YES",
    evaluationReport: "Strong Candidate. Will be lucky to hire him.",
    openingIds: ["JAN14-1", "JAN14-2"],
    candidacyId: "CAND-001",
    date: "2024-03-12T10:00:00Z",
    durationInMins: 30,
    candidateName: "Jeff Dean",
    candidateCurrentCompany: "Google",
    interviewers: [
      { name: "Urs Holzle", email: "urs@example.com" },
      { name: "Eric Schmidt", email: "eric@example.com" },
    ],
    hiringManager: { name: "Urs Holzle", email: "ursholzle@example.com" },
  },
];

export default function Interviews() {
  const [view, setView] = useState("list");

  const handleCellRender = (date: Dayjs) => {
    const dateString = date.format("YYYY-MM-DD");
    const interviewsForDate = interviewData.filter((interview) =>
      interview.date.startsWith(dateString)
    );

    return (
      <div style={{ maxHeight: "100px", overflowY: "scroll" }}>
        {interviewsForDate.map((interview) => (
          <Tag key={interview.interviewId} color="blue">
            {interview.candidateName}
          </Tag>
        ))}
      </div>
    );
  };

  const handleViewChange = (e: RadioChangeEvent) => {
    setView(e.target.value);
  };

  const columns = [
    {
      title: t("interviews.interview_id"),
      dataIndex: "interviewId",
      key: "interviewId",
      render: (interviewId: string) => (
        <Link to={`/interview/${interviewId}`}>{interviewId}</Link>
      ),
    },
    {
      title: t("interviews.status"),
      dataIndex: "interviewStatus",
      key: "interviewStatus",
    },
    {
      title: t("interviews.candidate"),
      dataIndex: "candidateName",
      key: "candidateName",
      render: (candidateName: string, record: any) =>
        `${candidateName}, ${record.candidateCurrentCompany}`,
    },
    {
      title: t("interviews.openings"),
      dataIndex: "openingIds",
      key: "openingIds",
      render: (openingIds: string[]) =>
        openingIds
          ? openingIds.map((id, index) => (
              <>
                <Link key={id} to={`/openings/${id}`}>
                  {id}
                </Link>
                {index < openingIds.length - 1 && ", "}
              </>
            ))
          : null,
    },
    {
      title: t("interviews.date"),
      dataIndex: "date",
      key: "date",
    },
    {
      title: t("interviews.interviewers"),
      dataIndex: "interviewers",
      key: "interviewers",
      render: (interviewers: any) =>
        interviewers.map((interviewer: any) => (
          <div key={interviewer.email}>
            {interviewer.name} ({interviewer.email})
          </div>
        )),
    },
    {
      title: t("interviews.evaluation_status"),
      dataIndex: "evaluationStatus",
      key: "evaluationStatus",
    },
    {
      title: t("interviews.evaluation_result"),
      dataIndex: "evaluationResult",
      key: "evaluationResult",
    },
    {
      title: t("interviews.evaluation_report"),
      dataIndex: "evaluationReport",
      key: "evaluationReport",
    },
  ];

  return (
    <Flex vertical gap="large">
      <Radio.Group
        defaultValue="list"
        onChange={handleViewChange}
        buttonStyle="solid"
      >
        <Radio.Button value="list">{t("interviews.view_list")}</Radio.Button>
        <Radio.Button value="calendar">
          {t("interviews.view_calendar")}
        </Radio.Button>
      </Radio.Group>

      {view === "calendar" ? (
        <Calendar cellRender={handleCellRender} mode="month" />
      ) : (
        <Table
          columns={columns}
          dataSource={interviewData}
          rowKey={(record) => record.interviewId}
          scroll={{ x: true }}
        />
      )}
    </Flex>
  );
}
