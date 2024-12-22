import React, { useState, useEffect } from "react";
import {
  Table,
  Tag,
  Space,
  Button,
  message,
  Select,
  DatePicker,
  Popconfirm,
} from "antd";
import {
  EyeOutlined,
  CheckOutlined,
  CloseOutlined,
  MessageOutlined,
} from "@ant-design/icons";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import type { Dayjs } from "dayjs";
import { Application, ApplicationState } from "@/types/application";

const { RangePicker } = DatePicker;
const { Option } = Select;

interface Filters {
  state: ApplicationState[];
  dateRange: [Dayjs, Dayjs] | null;
  opening_id?: string;
}

const Applications: React.FC = () => {
  const [applications, setApplications] = useState<Application[]>([]);
  const [loading, setLoading] = useState(false);
  const [filters, setFilters] = useState<Filters>({
    state: [],
    dateRange: null,
  });
  const navigate = useNavigate();

  const fetchApplications = async () => {
    try {
      setLoading(true);
      const response = await axios.post("/api/employer/filter-applications", {
        state: filters.state,
        from_date: filters.dateRange?.[0]?.format("YYYY-MM-DD"),
        to_date: filters.dateRange?.[1]?.format("YYYY-MM-DD"),
        opening_id: filters.opening_id,
      });
      setApplications(response.data);
    } catch (error) {
      message.error("Failed to fetch applications");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchApplications();
  }, [filters]);

  const handleStateChange = async (
    applicationId: string,
    fromState: ApplicationState,
    toState: ApplicationState
  ) => {
    try {
      await axios.post("/api/employer/change-application-state", {
        application_id: applicationId,
        from_state: fromState,
        to_state: toState,
      });
      message.success("Application state updated successfully");
      fetchApplications();
    } catch (error) {
      message.error("Failed to update application state");
    }
  };

  const getStateColor = (state: ApplicationState) => {
    switch (state) {
      case ApplicationState.SUBMITTED_APPLICATION:
        return "blue";
      case ApplicationState.SHORTLISTED_APPLICATION:
        return "green";
      case ApplicationState.REJECTED_APPLICATION:
        return "red";
      case ApplicationState.WITHDRAWN_APPLICATION:
        return "default";
      default:
        return "default";
    }
  };

  const columns = [
    {
      title: "Candidate",
      dataIndex: ["hub_user", "name"],
      key: "candidate",
    },
    {
      title: "Opening",
      dataIndex: ["opening", "title"],
      key: "opening",
    },
    {
      title: "Department",
      dataIndex: ["opening", "cost_center_name"],
      key: "department",
    },
    {
      title: "Applied On",
      dataIndex: "created_at",
      key: "created_at",
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: "Status",
      dataIndex: "state",
      key: "state",
      render: (state: ApplicationState) => (
        <Tag color={getStateColor(state)}>{state.replace("_", " ")}</Tag>
      ),
    },
    {
      title: "Actions",
      key: "actions",
      render: (_: any, record: Application) => (
        <Space>
          <Button
            icon={<EyeOutlined />}
            onClick={() => navigate(`/applications/${record.id}`)}
          />
          {record.state === ApplicationState.SUBMITTED_APPLICATION && (
            <>
              <Popconfirm
                title="Are you sure you want to shortlist this application?"
                onConfirm={() =>
                  handleStateChange(
                    record.id,
                    record.state,
                    ApplicationState.SHORTLISTED_APPLICATION
                  )
                }
                okText="Yes"
                cancelText="No"
              >
                <Button icon={<CheckOutlined />} type="primary" />
              </Popconfirm>
              <Popconfirm
                title="Are you sure you want to reject this application?"
                onConfirm={() =>
                  handleStateChange(
                    record.id,
                    record.state,
                    ApplicationState.REJECTED_APPLICATION
                  )
                }
                okText="Yes"
                cancelText="No"
              >
                <Button icon={<CloseOutlined />} danger />
              </Popconfirm>
            </>
          )}
          <Button
            icon={<MessageOutlined />}
            onClick={() => navigate(`/applications/${record.id}/messages`)}
          />
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: "flex", gap: 16 }}>
        <Select
          mode="multiple"
          placeholder="Filter by status"
          style={{ width: 200 }}
          onChange={(values: ApplicationState[]) =>
            setFilters((prev) => ({ ...prev, state: values }))
          }
        >
          {Object.values(ApplicationState).map((state) => (
            <Option key={state} value={state}>
              {state.replace("_APPLICATION", "").replace("_", " ")}
            </Option>
          ))}
        </Select>
        <RangePicker
          onChange={(dates) =>
            setFilters((prev) => ({
              ...prev,
              dateRange: dates as [Dayjs, Dayjs] | null,
            }))
          }
        />
      </div>

      <Table
        columns={columns}
        dataSource={applications}
        loading={loading}
        rowKey="id"
      />
    </div>
  );
};

export default Applications;
