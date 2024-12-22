import React, { useState, useEffect } from "react";
import {
  Table,
  Button,
  Tag,
  Space,
  message,
  Popconfirm,
  DatePicker,
  Select,
} from "antd";
import {
  PlusOutlined,
  EditOutlined,
  EyeOutlined,
  StopOutlined,
} from "@ant-design/icons";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import type { Dayjs } from "dayjs";
import { Opening, OpeningState } from "@/types/opening";

const { RangePicker } = DatePicker;
const { Option } = Select;

interface Filters {
  state: OpeningState[];
  dateRange: [Dayjs, Dayjs] | null;
}

const Openings: React.FC = () => {
  const [openings, setOpenings] = useState<Opening[]>([]);
  const [loading, setLoading] = useState(false);
  const [filters, setFilters] = useState<Filters>({
    state: [],
    dateRange: null,
  });
  const navigate = useNavigate();

  const fetchOpenings = async () => {
    try {
      setLoading(true);
      const response = await axios.post("/api/employer/filter-openings", {
        state: filters.state,
        from_date: filters.dateRange?.[0]?.format("YYYY-MM-DD"),
        to_date: filters.dateRange?.[1]?.format("YYYY-MM-DD"),
      });
      setOpenings(response.data);
    } catch (error) {
      message.error("Failed to fetch openings");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOpenings();
  }, [filters]);

  const handleStateChange = async (
    openingId: string,
    fromState: OpeningState,
    toState: OpeningState
  ) => {
    try {
      await axios.post("/api/employer/change-opening-state", {
        opening_id: openingId,
        from_state: fromState,
        to_state: toState,
      });
      message.success("Opening state updated successfully");
      fetchOpenings();
    } catch (error) {
      message.error("Failed to update opening state");
    }
  };

  const getStateColor = (state: OpeningState) => {
    switch (state) {
      case OpeningState.ACTIVE_OPENING:
        return "green";
      case OpeningState.DRAFT_OPENING:
        return "gold";
      case OpeningState.SUSPENDED_OPENING:
        return "red";
      case OpeningState.CLOSED_OPENING:
        return "default";
      default:
        return "default";
    }
  };

  const columns = [
    {
      title: "Title",
      dataIndex: "title",
      key: "title",
    },
    {
      title: "Positions",
      key: "positions",
      render: (record: Opening) =>
        `${record.filled_positions}/${record.positions}`,
    },
    {
      title: "Department",
      dataIndex: "cost_center_name",
      key: "cost_center_name",
    },
    {
      title: "Recruiter",
      dataIndex: ["recruiter", "name"],
      key: "recruiter",
    },
    {
      title: "Hiring Manager",
      dataIndex: ["hiring_manager", "name"],
      key: "hiring_manager",
    },
    {
      title: "Status",
      dataIndex: "state",
      key: "state",
      render: (state: OpeningState) => (
        <Tag color={getStateColor(state)}>{state.replace("_", " ")}</Tag>
      ),
    },
    {
      title: "Actions",
      key: "actions",
      render: (_: any, record: Opening) => (
        <Space>
          <Button
            icon={<EyeOutlined />}
            onClick={() => navigate(`/openings/${record.id}`)}
          />
          <Button
            icon={<EditOutlined />}
            onClick={() => navigate(`/openings/edit/${record.id}`)}
          />
          {record.state === OpeningState.ACTIVE_OPENING && (
            <Popconfirm
              title="Are you sure you want to suspend this opening?"
              onConfirm={() =>
                handleStateChange(
                  record.id,
                  record.state,
                  OpeningState.SUSPENDED_OPENING
                )
              }
              okText="Yes"
              cancelText="No"
            >
              <Button icon={<StopOutlined />} danger />
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: "flex", gap: 16 }}>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => navigate("/openings/create")}
        >
          Create Opening
        </Button>
        <Select
          mode="multiple"
          placeholder="Filter by status"
          style={{ width: 200 }}
          onChange={(values: OpeningState[]) =>
            setFilters((prev) => ({ ...prev, state: values }))
          }
        >
          {Object.values(OpeningState).map((state) => (
            <Option key={state} value={state}>
              {state.replace("_", " ")}
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
        dataSource={openings}
        loading={loading}
        rowKey="id"
      />
    </div>
  );
};

export default Openings;
