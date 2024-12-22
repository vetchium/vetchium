import React, { useState, useEffect } from "react";
import {
  Table,
  Card,
  Select,
  Button,
  Space,
  Tag,
  Modal,
  message,
  Typography,
  Timeline,
  Descriptions,
  Badge,
  Tooltip,
  DatePicker,
} from "antd";
import {
  ClockCircleOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  CalendarOutlined,
  ClearOutlined,
} from "@ant-design/icons";
import axios from "axios";
import styled from "styled-components";
import dayjs from "dayjs";
import {
  Candidacy,
  CandidacyState,
  CandidacyEvent,
  Interview,
  InterviewState,
  InterviewRating,
} from "@/types/candidacy";

const { Title, Text } = Typography;
const { Option } = Select;
const { RangePicker } = DatePicker;

const StyledCard = styled(Card)`
  margin-bottom: 24px;
`;

const TimelineContainer = styled.div`
  max-height: 400px;
  overflow-y: auto;
  padding: 16px;
`;

const Candidacies: React.FC = () => {
  const [candidacies, setCandidacies] = useState<Candidacy[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedCandidacy, setSelectedCandidacy] = useState<Candidacy | null>(
    null
  );
  const [detailsModalVisible, setDetailsModalVisible] = useState(false);
  const [filters, setFilters] = useState({
    state: [] as CandidacyState[],
    from_date: undefined as string | undefined,
    to_date: undefined as string | undefined,
  });

  const fetchCandidacies = async () => {
    setLoading(true);
    try {
      const response = await axios.get<{ items: Candidacy[] }>(
        "/api/hub/candidacies",
        { params: filters }
      );
      setCandidacies(response.data.items);
    } catch (error) {
      message.error("Failed to fetch candidacies");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCandidacies();
  }, [filters]);

  const handleWithdraw = async (candidacyId: string) => {
    try {
      await axios.post(`/api/hub/candidacies/${candidacyId}/withdraw`);
      message.success("Candidacy withdrawn successfully");
      fetchCandidacies();
    } catch (error) {
      message.error("Failed to withdraw candidacy");
    }
  };

  const handleRespondToOffer = async (candidacyId: string, accept: boolean) => {
    try {
      await axios.post(`/api/hub/candidacies/${candidacyId}/respond-offer`, {
        accept,
      });
      message.success(`Offer ${accept ? "accepted" : "declined"} successfully`);
      fetchCandidacies();
    } catch (error) {
      message.error(`Failed to ${accept ? "accept" : "decline"} offer`);
    }
  };

  const getTimelineIcon = (eventType: string) => {
    switch (eventType) {
      case "STATE_CHANGE":
        return <ClockCircleOutlined />;
      case "INTERVIEW_SCHEDULED":
      case "INTERVIEW_COMPLETED":
        return <CalendarOutlined />;
      case "OFFER_MADE":
      case "OFFER_ACCEPTED":
        return <CheckCircleOutlined />;
      case "OFFER_DECLINED":
      case "CANDIDACY_WITHDRAWN":
        return <CloseCircleOutlined />;
      default:
        return null;
    }
  };

  const getStatusColor = (state: CandidacyState) => {
    switch (state) {
      case CandidacyState.SCREENING_CANDIDACY:
      case CandidacyState.INTERVIEWING_CANDIDACY:
        return "processing";
      case CandidacyState.OFFERED_CANDIDACY:
        return "warning";
      case CandidacyState.ACCEPTED_CANDIDACY:
        return "success";
      case CandidacyState.DECLINED_CANDIDACY:
      case CandidacyState.REJECTED_CANDIDACY:
      case CandidacyState.WITHDRAWN_CANDIDACY:
        return "error";
      default:
        return "default";
    }
  };

  const columns = [
    {
      title: "Job Title",
      dataIndex: ["application", "opening", "title"],
      key: "title",
    },
    {
      title: "Employer",
      dataIndex: ["application", "opening", "employer_name"],
      key: "employer",
    },
    {
      title: "Status",
      dataIndex: "state",
      key: "state",
      render: (state: CandidacyState) => (
        <Badge
          status={getStatusColor(state)}
          text={state.replace("_CANDIDACY", "")}
        />
      ),
    },
    {
      title: "Interviews",
      key: "interviews",
      render: (_: any, record: Candidacy) => (
        <Space>
          {record.interviews.map((interview) => (
            <Tag
              key={interview.id}
              color={
                interview.state === InterviewState.COMPLETED_INTERVIEW
                  ? "green"
                  : interview.state === InterviewState.CANCELLED_INTERVIEW
                  ? "red"
                  : "blue"
              }
            >
              Round {interview.round}
            </Tag>
          ))}
        </Space>
      ),
    },
    {
      title: "Actions",
      key: "actions",
      render: (_: any, record: Candidacy) => (
        <Space>
          <Button
            type="link"
            onClick={() => {
              setSelectedCandidacy(record);
              setDetailsModalVisible(true);
            }}
          >
            View Details
          </Button>
          {record.state === CandidacyState.OFFERED_CANDIDACY && (
            <Space>
              <Button
                type="primary"
                onClick={() => handleRespondToOffer(record.id, true)}
              >
                Accept Offer
              </Button>
              <Button
                danger
                onClick={() => handleRespondToOffer(record.id, false)}
              >
                Decline Offer
              </Button>
            </Space>
          )}
          {(record.state === CandidacyState.SCREENING_CANDIDACY ||
            record.state === CandidacyState.INTERVIEWING_CANDIDACY) && (
            <Button
              danger
              onClick={() => {
                Modal.confirm({
                  title: "Withdraw Candidacy",
                  content:
                    "Are you sure you want to withdraw from this process?",
                  onOk: () => handleWithdraw(record.id),
                });
              }}
            >
              Withdraw
            </Button>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Title level={2}>My Candidacies</Title>

      <StyledCard>
        <Space wrap>
          <Select
            mode="multiple"
            value={filters.state}
            onChange={(value) => setFilters({ ...filters, state: value })}
            placeholder="Filter by status"
            style={{ width: 200 }}
            allowClear
          >
            {Object.values(CandidacyState).map((state) => (
              <Option key={state} value={state}>
                {state.replace("_CANDIDACY", "")}
              </Option>
            ))}
          </Select>

          <RangePicker
            onChange={(dates) => {
              if (dates) {
                setFilters({
                  ...filters,
                  from_date: dates[0]?.toISOString(),
                  to_date: dates[1]?.toISOString(),
                });
              } else {
                setFilters({
                  ...filters,
                  from_date: undefined,
                  to_date: undefined,
                });
              }
            }}
          />

          <Button
            icon={<ClearOutlined />}
            onClick={() =>
              setFilters({
                state: [],
                from_date: undefined,
                to_date: undefined,
              })
            }
          >
            Clear Filters
          </Button>
        </Space>
      </StyledCard>

      <Table
        columns={columns}
        dataSource={candidacies}
        loading={loading}
        rowKey="id"
      />

      <Modal
        title="Candidacy Details"
        open={detailsModalVisible}
        onCancel={() => {
          setDetailsModalVisible(false);
          setSelectedCandidacy(null);
        }}
        footer={null}
        width={800}
      >
        {selectedCandidacy && (
          <>
            <Descriptions bordered column={2}>
              <Descriptions.Item label="Job Title">
                {selectedCandidacy.application.opening.title}
              </Descriptions.Item>
              <Descriptions.Item label="Employer">
                {selectedCandidacy.application.opening.employer_name}
              </Descriptions.Item>
              <Descriptions.Item label="Status">
                <Badge
                  status={getStatusColor(selectedCandidacy.state)}
                  text={selectedCandidacy.state.replace("_CANDIDACY", "")}
                />
              </Descriptions.Item>
              <Descriptions.Item label="Applied On">
                {dayjs(selectedCandidacy.created_at).format("MMM D, YYYY")}
              </Descriptions.Item>
            </Descriptions>

            <Title level={4} style={{ marginTop: 24 }}>
              Interview Schedule
            </Title>
            {selectedCandidacy.interviews.map((interview) => (
              <Card
                key={interview.id}
                size="small"
                style={{ marginBottom: 16 }}
              >
                <Descriptions size="small" column={2}>
                  <Descriptions.Item label="Round">
                    Round {interview.round}
                  </Descriptions.Item>
                  <Descriptions.Item label="Status">
                    <Tag
                      color={
                        interview.state === InterviewState.COMPLETED_INTERVIEW
                          ? "green"
                          : interview.state ===
                            InterviewState.CANCELLED_INTERVIEW
                          ? "red"
                          : "blue"
                      }
                    >
                      {interview.state.replace("_INTERVIEW", "")}
                    </Tag>
                  </Descriptions.Item>
                  <Descriptions.Item label="Date & Time">
                    {dayjs(interview.scheduled_at).format("MMM D, YYYY HH:mm")}
                  </Descriptions.Item>
                  <Descriptions.Item label="Duration">
                    {interview.duration_minutes} minutes
                  </Descriptions.Item>
                  {interview.meeting_link && (
                    <Descriptions.Item label="Meeting Link" span={2}>
                      <a
                        href={interview.meeting_link}
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        Join Meeting
                      </a>
                    </Descriptions.Item>
                  )}
                  {interview.feedback &&
                    interview.state === InterviewState.COMPLETED_INTERVIEW && (
                      <>
                        <Descriptions.Item label="Feedback" span={2}>
                          {interview.feedback.feedback_to_candidate}
                        </Descriptions.Item>
                      </>
                    )}
                </Descriptions>
              </Card>
            ))}

            <Title level={4}>Timeline</Title>
            <TimelineContainer>
              <Timeline>
                {selectedCandidacy.timeline.map((event) => (
                  <Timeline.Item
                    key={event.id}
                    dot={getTimelineIcon(event.event_type)}
                  >
                    <Text strong>{event.event_type.replace(/_/g, " ")}</Text>
                    <br />
                    <Text>{event.description}</Text>
                    <br />
                    <Text type="secondary">
                      {dayjs(event.created_at).format("MMM D, YYYY HH:mm")}
                    </Text>
                  </Timeline.Item>
                ))}
              </Timeline>
            </TimelineContainer>
          </>
        )}
      </Modal>
    </div>
  );
};

export default Candidacies;
