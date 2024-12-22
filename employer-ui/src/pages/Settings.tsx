import React, { useState } from "react";
import {
  Card,
  Form,
  Input,
  Button,
  Divider,
  message,
  Switch,
  Space,
} from "antd";
import { LockOutlined } from "@ant-design/icons";
import axios from "axios";
import { useAuth } from "@/hooks/useAuth";

const Settings: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const { user } = useAuth();
  const [form] = Form.useForm();

  const handlePasswordChange = async (values: any) => {
    try {
      setLoading(true);
      await axios.post("/api/employer/change-password", {
        current_password: values.currentPassword,
        new_password: values.newPassword,
      });
      message.success("Password changed successfully");
      form.resetFields();
    } catch (error) {
      message.error("Failed to change password");
    } finally {
      setLoading(false);
    }
  };

  const handleTFAToggle = async (enabled: boolean) => {
    try {
      setLoading(true);
      if (enabled) {
        const response = await axios.post("/api/employer/enable-tfa", {});
        // Show QR code modal or handle TFA setup
        message.success("Two-factor authentication enabled");
      } else {
        await axios.post("/api/employer/disable-tfa", {});
        message.success("Two-factor authentication disabled");
      }
    } catch (error) {
      message.error("Failed to update two-factor authentication");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <Card title="Account Settings" style={{ marginBottom: 24 }}>
        <div style={{ marginBottom: 24 }}>
          <h3>Profile Information</h3>
          <p>
            <strong>Name:</strong> {user?.name}
          </p>
          <p>
            <strong>Email:</strong> {user?.email}
          </p>
          <p>
            <strong>Roles:</strong> {user?.roles.join(", ")}
          </p>
        </div>

        <Divider />

        <div style={{ marginBottom: 24 }}>
          <h3>Change Password</h3>
          <Form form={form} layout="vertical" onFinish={handlePasswordChange}>
            <Form.Item
              name="currentPassword"
              label="Current Password"
              rules={[
                {
                  required: true,
                  message: "Please enter your current password",
                },
              ]}
            >
              <Input.Password prefix={<LockOutlined />} />
            </Form.Item>

            <Form.Item
              name="newPassword"
              label="New Password"
              rules={[
                { required: true, message: "Please enter your new password" },
                { min: 8, message: "Password must be at least 8 characters" },
              ]}
            >
              <Input.Password prefix={<LockOutlined />} />
            </Form.Item>

            <Form.Item
              name="confirmPassword"
              label="Confirm New Password"
              dependencies={["newPassword"]}
              rules={[
                { required: true, message: "Please confirm your new password" },
                ({ getFieldValue }) => ({
                  validator(_, value) {
                    if (!value || getFieldValue("newPassword") === value) {
                      return Promise.resolve();
                    }
                    return Promise.reject(
                      new Error("The two passwords do not match")
                    );
                  },
                }),
              ]}
            >
              <Input.Password prefix={<LockOutlined />} />
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loading}>
                Change Password
              </Button>
            </Form.Item>
          </Form>
        </div>

        <Divider />

        <div>
          <h3>Security Settings</h3>
          <Space direction="vertical" size="large" style={{ width: "100%" }}>
            <div>
              <Space>
                <span>Two-Factor Authentication:</span>
                <Switch
                  checked={user?.state === "TFA_ENABLED"}
                  onChange={handleTFAToggle}
                  loading={loading}
                />
              </Space>
            </div>
          </Space>
        </div>
      </Card>
    </div>
  );
};

export default Settings;
