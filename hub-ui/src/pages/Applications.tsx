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
  Tooltip,
  DatePicker,
  Input,
} from "antd";
import {
  EyeOutlined,
  MessageOutlined,
  StopOutlined,
  FilterOutlined,
  ClearOutlined,
} from "@ant-design/icons";
import axios from "axios";
import styled from "styled-components";
import dayjs from "dayjs";
import {
  Application,
  ApplicationState,
  ApplicationMessage,
  SendMessageRequest,
} from "@/types/application";

const { Title, Text } = Typography;
const { Option } = Select;
const { RangePicker } = DatePicker;
const { TextArea } = Input;

const StyledCard = styled(Card)`
  margin-bottom: 24px;
`;

const MessageList = styled.div`
  max-height: 300px;
  overflow-y: auto;
  margin-bottom: 16px;
`;

const MessageItem = styled.div<{ isSelf: boolean }>`
  margin: 8px 0;
  padding: 8px 12px;
  border-radius: 8px;
  background-color: ${(props) => (props.isSelf ? "#e6f7ff" : "#f0f0f0")};
  align-self: ${(props) => (props.isSelf ? "flex-end" : "flex-start")};
  max-width: 80%;
`;

const Applications: React.FC = () => {
  const [applications, setApplications] = useState<Application[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedApplication, setSelectedApplication] =
    useState<Application | null>(null);
  const [messageModalVisible, setMessageModalVisible] = useState(false);
  const [messages, setMessages] = useState<ApplicationMessage[]>([]);
  const [newMessage, setNewMessage] = useState("");
  const [filters, setFilters] = useState({
    state: [] as ApplicationState[],
    from_date: undefined as string | undefined,
    to_date: undefined as string | undefined,
  });

  const fetchApplications = async () => {
    setLoading(true);
    try {
      const response = await axios.get<{ items: Application[] }>(
        "/api/hub/applications",
        { params: filters }
      );
      setApplications(response.data.items);
    } catch (error) {
      message.error("Failed to fetch applications");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchApplications();
  }, [filters]);

  const fetchMessages = async (applicationId: string) => {
    try {
      const response = await axios.get<{ items: ApplicationMessage[] }>(
        `/api/hub/applications/${applicationId}/messages`
      );
      setMessages(response.data.items);
    } catch (error) {
      message.error("Failed to fetch messages");
    }
  };

  const handleWithdraw = async (applicationId: string) => {
    try {
      await axios.post(`/api/hub/applications/${applicationId}/withdraw`);
      message.success("Application withdrawn successfully");
      fetchApplications();
    } catch (error) {
      message.error("Failed to withdraw application");
    }
  };

  const handleSendMessage = async () => {
    if (!selectedApplication || !newMessage.trim()) return;

    try {
      await axios.post(
        `/api/hub/applications/${selectedApplication.id}/messages`,
        { message: newMessage }
      );
      setNewMessage("");
      fetchMessages(selectedApplication.id);
      message.success("Message sent successfully");
    } catch (error) {
      message.error("Failed to send message");
    }
  };

  const columns = [
    {
      title: "Job Title",
      dataIndex: ["opening", "title"],
      key: "title",
    },
    {
      title: "Employer",
      dataIndex: ["opening", "employer_name"],
      key: "employer",
    },
    {
      title: "Status",
      dataIndex: "state",
      key: "state",
      render: (state: ApplicationState) => {
        const color =
          state === ApplicationState.SHORTLISTED_APPLICATION
            ? "green"
            : state === ApplicationState.REJECTED_APPLICATION
            ? "red"
            : state === ApplicationState.WITHDRAWN_APPLICATION
            ? "gray"
            : "blue";
        return <Tag color={color}>{state.replace("_APPLICATION", "")}</Tag>;
      },
    },
    {
      title: "Applied On",
      dataIndex: "created_at",
      key: "created_at",
      render: (date: string) => dayjs(date).format("MMM D, YYYY"),
    },
    {
      title: "Actions",
      key: "actions",
      render: (_: any, record: Application) => (
        <Space>
          <Tooltip title="View Details">
            <Button
              icon={<EyeOutlined />}
              onClick={() => {
                // TODO: Implement view details
              }}
            />
          </Tooltip>
          <Tooltip title="Messages">
            <Button
              icon={<MessageOutlined />}
              onClick={() => {
                setSelectedApplication(record);
                setMessageModalVisible(true);
                fetchMessages(record.id);
              }}
            />
          </Tooltip>
          {record.state === ApplicationState.SUBMITTED_APPLICATION && (
            <Tooltip title="Withdraw">
              <Button
                icon={<StopOutlined />}
                danger
                onClick={() => {
                  Modal.confirm({
                    title: "Withdraw Application",
                    content:
                      "Are you sure you want to withdraw this application?",
                    onOk: () => handleWithdraw(record.id),
                  });
                }}
              />
            </Tooltip>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Title level={2}>My Applications</Title>

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
            {Object.values(ApplicationState).map((state) => (
              <Option key={state} value={state}>
                {state.replace("_APPLICATION", "")}
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
        dataSource={applications}
        loading={loading}
        rowKey="id"
      />

      <Modal
        title="Messages"
        open={messageModalVisible}
        onCancel={() => {
          setMessageModalVisible(false);
          setSelectedApplication(null);
          setMessages([]);
          setNewMessage("");
        }}
        footer={null}
        width={600}
      >
        <MessageList>
          {messages.map((msg) => (
            <MessageItem key={msg.id} isSelf={msg.sender_type === "HUB_USER"}>
              <Text strong>
                {msg.sender_type === "HUB_USER" ? "You" : "Employer"}
              </Text>
              <div>{msg.message}</div>
              <Text type="secondary" style={{ fontSize: "12px" }}>
                {dayjs(msg.created_at).format("MMM D, YYYY HH:mm")}
              </Text>
            </MessageItem>
          ))}
        </MessageList>

        <Space.Compact style={{ width: "100%" }}>
          <TextArea
            value={newMessage}
            onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
              setNewMessage(e.target.value)
            }
            placeholder="Type your message..."
            autoSize={{ minRows: 2, maxRows: 4 }}
          />
          <Button type="primary" onClick={handleSendMessage}>
            Send
          </Button>
        </Space.Compact>
      </Modal>
    </div>
  );
};

export default Applications;
