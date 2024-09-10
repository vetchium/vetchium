import { Divider, Flex, Space, Table, Typography } from "antd";
import { t } from "i18next";
import { Link } from "react-router-dom";
import InterviewCancelForm from "../forms/InterviewCancelForm";
import InterviewFeedbackForm from "../forms/InterviewFeedbackForm";

const { Text } = Typography;

const data = {
  interviewId: "INT-001",
  interviewStatus: "SCHEDULED",
  openings: [
    {
      id: "JAN14-1",
      title: "Distinguished Engineer",
      hiringManager: { name: "Urs Holzle", email: "ursholzle@example.com" },
    },
    {
      id: "JAN14-2",
      title: "Fellow",
      hiringManager: { name: "Larry Page", email: "larrypage@example.com" },
    },
  ],
  candidacyId: "cand0123123",
  date: "2024-03-14T10:00:00Z",
  durationInMins: 60,
  candidateName: "Jeff Dean",
  candidateCurrentCompany: "Google",
  interviewers: [
    { name: "Larry Page", email: "larrypage@example.com" },
    { name: "Urs Holzle", email: "ursholzle@example.com" },
  ],
};

export default function Interview() {
  const columns = [
    {
      title: "Field Name",
      dataIndex: "field",
      key: "field",
    },
    {
      title: "Value",
      dataIndex: "value",
      key: "value",
    },
  ];

  const dataSource = [
    { key: "1", field: t("interviews.status"), value: data.interviewStatus },
    { key: "2", field: t("interviews.candidate"), value: data.candidateName },
    {
      key: "3",
      field: t("interviews.candidate_current_company"),
      value: data.candidateCurrentCompany,
    },
    {
      key: "4",
      field: t("interviews.openings"),
      value: data.openings.map((opening) => (
        <>
          <Space key={opening.id}>
            <Link to={`/opening/${opening.id}`}>{opening.title}</Link>
            <Text>
              {t("interviews.for")} {opening.hiringManager.name} (
              {opening.hiringManager.email})
            </Text>
          </Space>
          <br />
        </>
      )),
    },
    {
      key: "5",
      field: t("interviews.candidacy_id"),
      value: (
        <Link to={`/candidacy/${data.candidacyId}`}>{data.candidacyId}</Link>
      ),
    },
    {
      key: "6",
      field: t("interviews.date"),
      value: new Date(data.date).toLocaleString(),
    },
    { key: "7", field: t("interviews.duration"), value: data.durationInMins },
    {
      key: "8",
      field: t("interviews.interviewers"),
      value: data.interviewers.map((i) => `${i.name} (${i.email})`).join(", "),
    },
  ];

  return (
    <Flex vertical style={{ margin: "2rem" }}>
      <Table
        dataSource={dataSource}
        columns={columns}
        pagination={false}
        showHeader={false}
        style={{ border: "none", margin: "3rem" }}
      />
      <Divider />
      <InterviewFeedbackForm />
      <Divider />

      {/* Should be shown only for recruiters and admins */}
      <InterviewCancelForm />
    </Flex>
  );
}
