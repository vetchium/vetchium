import React, { useState } from "react";
import {
  Card,
  Form,
  Input,
  Button,
  message,
  Typography,
  Divider,
  Upload,
  Space,
  Switch,
  Select,
} from "antd";
import {
  UserOutlined,
  PhoneOutlined,
  LinkedinOutlined,
  GithubOutlined,
  GlobalOutlined,
  UploadOutlined,
  LockOutlined,
} from "@ant-design/icons";
import type { UploadProps } from "antd";
import axios from "axios";
import styled from "styled-components";
import { useAuth } from "@/hooks/useAuth";
import { UpdateProfileRequest, ChangePasswordRequest } from "@/types/auth";

const { Title } = Typography;
const { Option } = Select;

const StyledCard = styled(Card)`
  max-width: 600px;
  margin: 0 auto;
`;

const Settings: React.FC = () => {
  const { user } = useAuth();
  const [profileForm] = Form.useForm();
  const [passwordForm] = Form.useForm();
  const [preferencesForm] = Form.useForm();
  const [loading, setLoading] = useState(false);

  const handleUpdateProfile = async (values: UpdateProfileRequest) => {
    setLoading(true);
    try {
      await axios.put("/api/hub/profile", values);
      message.success("Profile updated successfully");
    } catch (error) {
      message.error("Failed to update profile");
    } finally {
      setLoading(false);
    }
  };

  const handleChangePassword = async (values: ChangePasswordRequest) => {
    setLoading(true);
    try {
      await axios.put("/api/hub/change-password", values);
      message.success("Password changed successfully");
      passwordForm.resetFields();
    } catch (error) {
      message.error("Failed to change password");
    } finally {
      setLoading(false);
    }
  };

  const handleUpdatePreferences = async (values: {
    theme: "light" | "dark";
    language: string;
  }) => {
    setLoading(true);
    try {
      await axios.put("/api/hub/preferences", values);
      message.success("Preferences updated successfully");
    } catch (error) {
      message.error("Failed to update preferences");
    } finally {
      setLoading(false);
    }
  };

  const uploadProps: UploadProps = {
    name: "resume",
    action: "/api/hub/upload",
    onChange(info) {
      if (info.file.status === "done") {
        profileForm.setFieldsValue({
          resume_url: info.file.response.url,
        });
        message.success("Resume uploaded successfully");
      } else if (info.file.status === "error") {
        message.error("Resume upload failed");
      }
    },
  };

  return (
    <div>
      <Title level={2}>Settings</Title>

      <StyledCard>
        <Title level={4}>Profile Information</Title>
        <Form
          form={profileForm}
          layout="vertical"
          initialValues={{
            name: user?.name,
            phone: user?.phone,
            linkedin_url: user?.linkedin_url,
            github_url: user?.github_url,
            portfolio_url: user?.portfolio_url,
            resume_url: user?.resume_url,
          }}
          onFinish={handleUpdateProfile}
        >
          <Form.Item
            name="name"
            label="Full Name"
            rules={[{ required: true, message: "Please enter your name" }]}
          >
            <Input prefix={<UserOutlined />} placeholder="Your full name" />
          </Form.Item>

          <Form.Item name="phone" label="Phone Number">
            <Input prefix={<PhoneOutlined />} placeholder="Your phone number" />
          </Form.Item>

          <Form.Item name="linkedin_url" label="LinkedIn Profile">
            <Input
              prefix={<LinkedinOutlined />}
              placeholder="LinkedIn profile URL"
            />
          </Form.Item>

          <Form.Item name="github_url" label="GitHub Profile">
            <Input
              prefix={<GithubOutlined />}
              placeholder="GitHub profile URL"
            />
          </Form.Item>

          <Form.Item name="portfolio_url" label="Portfolio Website">
            <Input
              prefix={<GlobalOutlined />}
              placeholder="Portfolio website URL"
            />
          </Form.Item>

          <Form.Item name="resume_url" label="Resume">
            <Space direction="vertical" style={{ width: "100%" }}>
              {user?.resume_url && (
                <a
                  href={user.resume_url}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Current Resume
                </a>
              )}
              <Upload {...uploadProps}>
                <Button icon={<UploadOutlined />}>Upload New Resume</Button>
              </Upload>
            </Space>
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              Update Profile
            </Button>
          </Form.Item>
        </Form>

        <Divider />

        <Title level={4}>Change Password</Title>
        <Form
          form={passwordForm}
          layout="vertical"
          onFinish={handleChangePassword}
        >
          <Form.Item
            name="current_password"
            label="Current Password"
            rules={[
              { required: true, message: "Please enter your current password" },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="Current password"
            />
          </Form.Item>

          <Form.Item
            name="new_password"
            label="New Password"
            rules={[
              { required: true, message: "Please enter your new password" },
              { min: 8, message: "Password must be at least 8 characters" },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="New password"
            />
          </Form.Item>

          <Form.Item
            name="confirm_password"
            label="Confirm New Password"
            dependencies={["new_password"]}
            rules={[
              { required: true, message: "Please confirm your new password" },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue("new_password") === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(
                    new Error("The two passwords do not match")
                  );
                },
              }),
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="Confirm new password"
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              Change Password
            </Button>
          </Form.Item>
        </Form>

        <Divider />

        <Title level={4}>Preferences</Title>
        <Form
          form={preferencesForm}
          layout="vertical"
          initialValues={{
            theme: "light",
            language: "en",
          }}
          onFinish={handleUpdatePreferences}
        >
          <Form.Item name="theme" label="Theme" rules={[{ required: true }]}>
            <Space>
              <Switch
                checkedChildren="Dark"
                unCheckedChildren="Light"
                onChange={(checked) =>
                  preferencesForm.setFieldsValue({
                    theme: checked ? "dark" : "light",
                  })
                }
              />
              <span>
                {preferencesForm.getFieldValue("theme") === "dark"
                  ? "Dark Mode"
                  : "Light Mode"}
              </span>
            </Space>
          </Form.Item>

          <Form.Item
            name="language"
            label="Language"
            rules={[{ required: true }]}
          >
            <Select style={{ width: 200 }}>
              <Option value="en">English</Option>
              <Option value="es">Español</Option>
              <Option value="fr">Français</Option>
              <Option value="de">Deutsch</Option>
              <Option value="zh">中文</Option>
              <Option value="ja">日本語</Option>
              <Option value="ko">한국어</Option>
            </Select>
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              Update Preferences
            </Button>
          </Form.Item>
        </Form>
      </StyledCard>
    </div>
  );
};

export default Settings;
