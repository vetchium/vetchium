import { PlusCircleTwoTone } from "@ant-design/icons";
import { Button, Flex, Modal, Table, Typography } from "antd";
import { useNavigate, useParams } from "react-router-dom";
import t from "../i18n/i18n";

const { Text, Link } = Typography;

const data = {
  name: "Jeff Dean",
  lastPosition: "Senior Fellow",
  lastCompany: "Google",
  shortlistedOpenings: [
    {
      id: "JAN14-1",
      hiringManager: "a@example.com",
      title: "Distinguished Engineer",
    },
    {
      id: "JAN14-2",
      hiringManager: "b@example.com",
      title: "Fellow",
    },
    {
      id: "JAN14-3",
      hiringManager: "c@example.com",
      title: "Senior Fellow",
    },
  ],
  interviews: [
    {
      id: "9m4e2mr0ui3e8a215n200",
      status: "CANCELLED",
    },
    {
      id: "9m4e2mr0ui3e8a3215n000",
      status: "COMPLETED",
      interviewers: ["A <a@example.com>", "B <b@example.com>"],
      at: "2024-03-14T10:00:00Z",
      evaluation_status: "EVALUATION_PENDING",
    },
    {
      id: "9m4e2mr0ui3e8a32215n000",
      status: "COMPLETED",
      interviewers: ["C <c@example.com>"],
      at: "2024-03-14T10:00:00Z",
      evaluation_status: "EVALUATION_COMPLETED",
      evaluation: {
        positives: "everything",
        negatives: "none",
        result: "STRONG_YES",
        summary: `will be a good fit
          
          Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
          `,
        feedback: "keep doing the good work",
      },
    },
    {
      id: "9m4e2mr0ui3e8a215n000",
      status: "SCHEDULED",
      interviewers: ["d@example.com"],
      at: "2024-03-14T10:00:00Z",
    },
    {
      id: "9m4e2mr0ui3e8a215n100",
      status: "PENDING_SCHEDULE",
    },
  ],
};

export default function Candidacy() {
  const navigate = useNavigate();
  const { candidacy_id } = useParams();

  return (
    <Flex vertical>
      <Typography.Title level={2}>{data.name}</Typography.Title>
      <Text>{data.lastPosition + ", " + data.lastCompany}</Text>

      <Typography.Title level={5}>Shortlisted for Opening(s)</Typography.Title>
      <Table
        columns={[
          {
            title: "Opening ID",
            dataIndex: "opening",
            key: "opening",
            render: (opening) => (
              <Link href={`/openings/${opening.id}`}>{opening.id}</Link>
            ),
          },
          {
            title: "Hiring Manager",
            dataIndex: "opening",
            key: "hiringManager",
            render: (opening) => opening.hiringManager,
          },
          {
            title: "Title",
            dataIndex: "opening",
            key: "title",
            render: (opening) => opening.title,
          },
        ]}
        dataSource={data.shortlistedOpenings.map((opening, index) => ({
          key: index,
          opening,
        }))}
        pagination={false}
      />

      <Typography.Title level={5}>Interviews</Typography.Title>
      <Table
        columns={[
          {
            title: "Interviewers",
            dataIndex: "interviewers",
            key: "interviewers",
            render: (interviewers) =>
              interviewers ? interviewers.join(", ") : "",
          },
          { title: "Status", dataIndex: "status", key: "status" },
          {
            title: "Evaluation Status",
            dataIndex: "evaluation_status",
            key: "evaluation_status",
          },
          {
            title: "Result",
            dataIndex: "evaluation",
            key: "evaluation",
            render: (evaluation) => (evaluation ? evaluation.result : ""),
          },
          {
            title: "Evaluation Report",
            key: "evaluation_report",
            render: (record) =>
              record.evaluation_status === "EVALUATION_COMPLETED" ? (
                <Link
                  onClick={() => {
                    Modal.info({
                      title: "Evaluation Report",
                      content: (
                        <div
                          style={{
                            maxHeight: "400px",
                            overflowY: "scroll",
                          }}
                        >
                          <Typography.Title level={4}>
                            Positives
                          </Typography.Title>
                          <p>{record.evaluation.positives}</p>
                          <Typography.Title level={4}>
                            Negatives
                          </Typography.Title>
                          <p>{record.evaluation.negatives}</p>
                          <Typography.Title level={4}>Summary</Typography.Title>
                          <p>{record.evaluation.summary}</p>
                          <Typography.Title level={4}>
                            Feedback
                          </Typography.Title>
                          <p>{record.evaluation.feedback}</p>
                        </div>
                      ),
                    });
                  }}
                >
                  Evaluation Report
                </Link>
              ) : null,
          },
        ]}
        dataSource={data.interviews}
        rowKey="id"
        pagination={false}
      />
      <Button
        type="primary"
        icon={<PlusCircleTwoTone />}
        onClick={() => navigate(`/create-interview/${candidacy_id}`)}
        style={{ marginTop: "1rem", marginLeft: "2rem", width: "fit-content" }}
      >
        {t("create_interview")}
      </Button>
    </Flex>
  );
}
