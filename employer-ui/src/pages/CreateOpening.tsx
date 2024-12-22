import React, { useState, useEffect } from "react";
import {
  Form,
  Input,
  Button,
  Select,
  InputNumber,
  message,
  Card,
  Space,
} from "antd";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import {
  CreateOpeningRequest,
  OpeningType,
  EducationLevel,
  Salary,
} from "@/types/opening";
import { OrgUserShort } from "@/types/auth";
import { CostCenter } from "@/types/costCenter";
import { Location } from "@/types/location";

const { Option } = Select;
const { TextArea } = Input;

const CreateOpening: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [orgUsers, setOrgUsers] = useState<OrgUserShort[]>([]);
  const [departments, setDepartments] = useState<CostCenter[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [usersRes, deptsRes, locsRes] = await Promise.all([
          axios.post("/api/employer/get-org-users", {}),
          axios.post("/api/employer/get-cost-centers", {}),
          axios.post("/api/employer/get-locations", {}),
        ]);

        setOrgUsers(usersRes.data);
        setDepartments(deptsRes.data);
        setLocations(locsRes.data);
      } catch (error) {
        message.error("Failed to fetch required data");
      }
    };

    fetchData();
  }, []);

  const onFinish = async (values: any) => {
    try {
      setLoading(true);
      const openingData: CreateOpeningRequest = {
        ...values,
        salary: values.salary_enabled
          ? {
              min_amount: values.min_amount,
              max_amount: values.max_amount,
              currency: values.currency,
            }
          : undefined,
      };

      await axios.post("/api/employer/create-opening", openingData);
      message.success("Opening created successfully");
      navigate("/openings");
    } catch (error) {
      message.error("Failed to create opening");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card title="Create New Opening">
      <Form
        form={form}
        layout="vertical"
        onFinish={onFinish}
        initialValues={{
          opening_type: OpeningType.FULL_TIME_OPENING,
          min_education_level: EducationLevel.NOT_MATTERS_EDUCATION,
          yoe_min: 0,
          yoe_max: 5,
          salary_enabled: false,
        }}
      >
        <Form.Item
          name="title"
          label="Title"
          rules={[{ required: true, message: "Please enter the job title" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="positions"
          label="Number of Positions"
          rules={[
            { required: true, message: "Please enter number of positions" },
          ]}
        >
          <InputNumber min={1} />
        </Form.Item>

        <Form.Item
          name="jd"
          label="Job Description"
          rules={[{ required: true, message: "Please enter job description" }]}
        >
          <TextArea rows={6} />
        </Form.Item>

        <Space size="large" style={{ display: "flex" }}>
          <Form.Item
            name="recruiter"
            label="Recruiter"
            rules={[{ required: true, message: "Please select a recruiter" }]}
          >
            <Select style={{ width: 200 }}>
              {orgUsers.map((user) => (
                <Option key={user.id} value={user.id}>
                  {user.name}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="hiring_manager"
            label="Hiring Manager"
            rules={[
              { required: true, message: "Please select a hiring manager" },
            ]}
          >
            <Select style={{ width: 200 }}>
              {orgUsers.map((user) => (
                <Option key={user.id} value={user.id}>
                  {user.name}
                </Option>
              ))}
            </Select>
          </Form.Item>
        </Space>

        <Form.Item name="hiring_team_members" label="Hiring Team Members">
          <Select mode="multiple" style={{ width: "100%" }}>
            {orgUsers.map((user) => (
              <Option key={user.id} value={user.id}>
                {user.name}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="cost_center_name"
          label="Department"
          rules={[{ required: true, message: "Please select a department" }]}
        >
          <Select>
            {departments.map((dept) => (
              <Option key={dept.name} value={dept.name}>
                {dept.name}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item name="location_titles" label="Locations">
          <Select mode="multiple">
            {locations.map((loc) => (
              <Option key={loc.title} value={loc.title}>
                {loc.title}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Space size="large" style={{ display: "flex" }}>
          <Form.Item
            name="opening_type"
            label="Employment Type"
            rules={[
              { required: true, message: "Please select employment type" },
            ]}
          >
            <Select style={{ width: 200 }}>
              {Object.values(OpeningType).map((type) => (
                <Option key={type} value={type}>
                  {type.replace("_OPENING", "").replace("_", " ")}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="min_education_level"
            label="Minimum Education"
            rules={[
              { required: true, message: "Please select minimum education" },
            ]}
          >
            <Select style={{ width: 200 }}>
              {Object.values(EducationLevel).map((level) => (
                <Option key={level} value={level}>
                  {level.replace("_EDUCATION", "").replace("_", " ")}
                </Option>
              ))}
            </Select>
          </Form.Item>
        </Space>

        <Space size="large" style={{ display: "flex" }}>
          <Form.Item
            name="yoe_min"
            label="Minimum Years of Experience"
            rules={[{ required: true, message: "Please enter minimum YOE" }]}
          >
            <InputNumber min={0} />
          </Form.Item>

          <Form.Item
            name="yoe_max"
            label="Maximum Years of Experience"
            rules={[{ required: true, message: "Please enter maximum YOE" }]}
          >
            <InputNumber min={0} />
          </Form.Item>
        </Space>

        <Form.Item name="employer_notes" label="Internal Notes">
          <TextArea rows={4} />
        </Form.Item>

        <Form.Item name="salary_enabled" valuePropName="checked">
          <Form.Item
            noStyle
            shouldUpdate={(prevValues, currentValues) =>
              prevValues.salary_enabled !== currentValues.salary_enabled
            }
          >
            {({ getFieldValue }) =>
              getFieldValue("salary_enabled") && (
                <Space size="large" style={{ display: "flex" }}>
                  <Form.Item
                    name="min_amount"
                    label="Minimum Salary"
                    rules={[
                      {
                        required: true,
                        message: "Please enter minimum salary",
                      },
                    ]}
                  >
                    <InputNumber min={0} />
                  </Form.Item>

                  <Form.Item
                    name="max_amount"
                    label="Maximum Salary"
                    rules={[
                      {
                        required: true,
                        message: "Please enter maximum salary",
                      },
                    ]}
                  >
                    <InputNumber min={0} />
                  </Form.Item>

                  <Form.Item
                    name="currency"
                    label="Currency"
                    rules={[
                      {
                        required: true,
                        message: "Please select currency",
                      },
                    ]}
                  >
                    <Select style={{ width: 120 }}>
                      <Option value="USD">USD</Option>
                      <Option value="EUR">EUR</Option>
                      <Option value="GBP">GBP</Option>
                      <Option value="INR">INR</Option>
                    </Select>
                  </Form.Item>
                </Space>
              )
            }
          </Form.Item>
        </Form.Item>

        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit" loading={loading}>
              Create Opening
            </Button>
            <Button onClick={() => navigate("/openings")}>Cancel</Button>
          </Space>
        </Form.Item>
      </Form>
    </Card>
  );
};

export default CreateOpening;
