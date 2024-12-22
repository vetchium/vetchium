import React, { useState, useEffect } from "react";
import {
  Table,
  Card,
  Input,
  Select,
  Button,
  Form,
  Space,
  Tag,
  Modal,
  Upload,
  message,
  Typography,
} from "antd";
import {
  SearchOutlined,
  UploadOutlined,
  FilterOutlined,
  ClearOutlined,
} from "@ant-design/icons";
import type { UploadProps } from "antd";
import axios from "axios";
import styled from "styled-components";
import { Opening, OpeningType, EducationLevel } from "@/types/opening";
import { CreateApplicationRequest } from "@/types/application";

const { Title } = Typography;
const { Option } = Select;
const { TextArea } = Input;

const StyledCard = styled(Card)`
  margin-bottom: 24px;
`;

const Openings: React.FC = () => {
  const [openings, setOpenings] = useState<Opening[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedOpening, setSelectedOpening] = useState<Opening | null>(null);
  const [applyModalVisible, setApplyModalVisible] = useState(false);
  const [applyForm] = Form.useForm();
  const [selectedCountry, setSelectedCountry] = useState<string>("");
  const [locations, setLocations] = useState<string[]>([]);
  const [filters, setFilters] = useState({
    employer_name: "",
    opening_type: [] as OpeningType[],
    min_education_level: [] as EducationLevel[],
    yoe_min: undefined as number | undefined,
    yoe_max: undefined as number | undefined,
    country_code: undefined as string | undefined,
    location_titles: [] as string[],
  });

  const fetchOpenings = async () => {
    setLoading(true);
    try {
      const response = await axios.get<{ items: Opening[] }>(
        "/api/hub/openings",
        {
          params: filters,
        }
      );
      setOpenings(response.data.items);
    } catch (error) {
      message.error("Failed to fetch openings");
    } finally {
      setLoading(false);
    }
  };

  const fetchLocations = async (countryCode: string) => {
    try {
      const response = await axios.get<{ items: string[] }>(
        `/api/hub/locations/${countryCode}`
      );
      setLocations(response.data.items);
    } catch (error) {
      message.error("Failed to fetch locations");
    }
  };

  useEffect(() => {
    if (selectedCountry) {
      fetchLocations(selectedCountry);
      setFilters((prev) => ({
        ...prev,
        country_code: selectedCountry,
        location_titles: [],
      }));
    }
  }, [selectedCountry]);

  useEffect(() => {
    fetchOpenings();
  }, [filters]);

  const handleApply = async (values: CreateApplicationRequest) => {
    try {
      await axios.post("/api/hub/applications", values);
      message.success("Application submitted successfully");
      setApplyModalVisible(false);
      applyForm.resetFields();
    } catch (error) {
      message.error("Failed to submit application");
    }
  };

  const uploadProps: UploadProps = {
    name: "resume",
    action: "/api/hub/upload",
    onChange(info) {
      if (info.file.status === "done") {
        applyForm.setFieldsValue({
          resume_url: info.file.response.url,
        });
        message.success("Resume uploaded successfully");
      } else if (info.file.status === "error") {
        message.error("Resume upload failed");
      }
    },
  };

  const columns = [
    {
      title: "Title",
      dataIndex: "title",
      key: "title",
      render: (text: string, record: Opening) => (
        <a
          onClick={() => {
            setSelectedOpening(record);
            setApplyModalVisible(true);
          }}
        >
          {text}
        </a>
      ),
    },
    {
      title: "Employer",
      dataIndex: "employer_name",
      key: "employer_name",
    },
    {
      title: "Type",
      dataIndex: "opening_type",
      key: "opening_type",
      render: (type: OpeningType) => (
        <Tag color="blue">{type.replace("_OPENING", "")}</Tag>
      ),
    },
    {
      title: "Location",
      key: "location",
      render: (_: any, record: Opening) => (
        <>
          {record.location_titles?.map((location) => (
            <Tag key={location}>{location}</Tag>
          ))}
          {record.remote_country_codes?.length ? (
            <Tag color="green">REMOTE</Tag>
          ) : null}
        </>
      ),
    },
    {
      title: "Experience",
      key: "experience",
      render: (_: any, record: Opening) =>
        `${record.yoe_min}-${record.yoe_max} years`,
    },
    {
      title: "Actions",
      key: "actions",
      render: (_: any, record: Opening) => (
        <Button
          type="primary"
          onClick={() => {
            setSelectedOpening(record);
            setApplyModalVisible(true);
          }}
        >
          Apply
        </Button>
      ),
    },
  ];

  return (
    <div>
      <Title level={2}>Job Openings</Title>

      <StyledCard>
        <Form layout="vertical">
          <Space wrap>
            <Form.Item label="Search by employer">
              <Input
                prefix={<SearchOutlined />}
                value={filters.employer_name}
                onChange={(e) =>
                  setFilters({ ...filters, employer_name: e.target.value })
                }
                placeholder="Search employers"
                allowClear
              />
            </Form.Item>

            <Form.Item label="Country">
              <Select
                value={selectedCountry}
                onChange={setSelectedCountry}
                placeholder="Select country"
                style={{ width: 200 }}
                allowClear
                onClear={() => {
                  setLocations([]);
                  setFilters((prev) => ({
                    ...prev,
                    country_code: undefined,
                    location_titles: [],
                  }));
                }}
              >
                <Option value="US">United States</Option>
                <Option value="GB">United Kingdom</Option>
                <Option value="CA">Canada</Option>
                <Option value="AU">Australia</Option>
                <Option value="IN">India</Option>
                <Option value="SG">Singapore</Option>
                {/* Add more countries as needed */}
              </Select>
            </Form.Item>

            <Form.Item label="Locations">
              <Select
                mode="multiple"
                value={filters.location_titles}
                onChange={(value) =>
                  setFilters({ ...filters, location_titles: value })
                }
                placeholder="Select locations"
                style={{ width: 300 }}
                allowClear
                disabled={!selectedCountry}
              >
                {locations.map((location) => (
                  <Option key={location} value={location}>
                    {location}
                  </Option>
                ))}
              </Select>
            </Form.Item>

            <Form.Item label="Job Type">
              <Select
                mode="multiple"
                value={filters.opening_type}
                onChange={(value) =>
                  setFilters({ ...filters, opening_type: value })
                }
                placeholder="Select job types"
                style={{ width: 200 }}
                allowClear
              >
                {Object.values(OpeningType).map((type) => (
                  <Option key={type} value={type}>
                    {type.replace("_OPENING", "")}
                  </Option>
                ))}
              </Select>
            </Form.Item>

            <Form.Item label="Education Level">
              <Select
                mode="multiple"
                value={filters.min_education_level}
                onChange={(value) =>
                  setFilters({ ...filters, min_education_level: value })
                }
                placeholder="Select education level"
                style={{ width: 200 }}
                allowClear
              >
                {Object.values(EducationLevel).map((level) => (
                  <Option key={level} value={level}>
                    {level.replace("_EDUCATION", "")}
                  </Option>
                ))}
              </Select>
            </Form.Item>

            <Form.Item label="Experience (Years)">
              <Space>
                <Input
                  type="number"
                  value={filters.yoe_min}
                  onChange={(e) =>
                    setFilters({
                      ...filters,
                      yoe_min: e.target.value
                        ? Number(e.target.value)
                        : undefined,
                    })
                  }
                  placeholder="Min"
                  style={{ width: 100 }}
                />
                <Input
                  type="number"
                  value={filters.yoe_max}
                  onChange={(e) =>
                    setFilters({
                      ...filters,
                      yoe_max: e.target.value
                        ? Number(e.target.value)
                        : undefined,
                    })
                  }
                  placeholder="Max"
                  style={{ width: 100 }}
                />
              </Space>
            </Form.Item>

            <Form.Item label=" ">
              <Button
                icon={<ClearOutlined />}
                onClick={() =>
                  setFilters({
                    employer_name: "",
                    opening_type: [],
                    min_education_level: [],
                    yoe_min: undefined,
                    yoe_max: undefined,
                    country_code: undefined,
                    location_titles: [],
                  })
                }
              >
                Clear Filters
              </Button>
            </Form.Item>
          </Space>
        </Form>
      </StyledCard>

      <Table
        columns={columns}
        dataSource={openings}
        loading={loading}
        rowKey="id"
      />

      <Modal
        title="Apply for Position"
        open={applyModalVisible}
        onCancel={() => {
          setApplyModalVisible(false);
          setSelectedOpening(null);
          applyForm.resetFields();
        }}
        footer={null}
      >
        {selectedOpening && (
          <Form form={applyForm} onFinish={handleApply} layout="vertical">
            <Form.Item
              name="opening_id"
              initialValue={selectedOpening.id}
              hidden
            >
              <Input />
            </Form.Item>

            <Form.Item
              name="resume_url"
              label="Resume"
              rules={[
                { required: true, message: "Please upload your resume!" },
              ]}
            >
              <Upload {...uploadProps}>
                <Button icon={<UploadOutlined />}>Upload Resume</Button>
              </Upload>
            </Form.Item>

            <Form.Item name="cover_letter" label="Cover Letter">
              <TextArea
                placeholder="Write a cover letter (optional)"
                rows={4}
              />
            </Form.Item>

            <Form.Item>
              <Space>
                <Button type="primary" htmlType="submit">
                  Submit Application
                </Button>
                <Button
                  onClick={() => {
                    setApplyModalVisible(false);
                    setSelectedOpening(null);
                    applyForm.resetFields();
                  }}
                >
                  Cancel
                </Button>
              </Space>
            </Form.Item>
          </Form>
        )}
      </Modal>
    </div>
  );
};

export default Openings;
