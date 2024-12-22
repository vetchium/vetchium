import React, { useState, useEffect } from "react";
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Select,
  Space,
  message,
  Popconfirm,
  Tag,
} from "antd";
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  LockOutlined,
} from "@ant-design/icons";
import axios from "axios";
import { OrgUser } from "@/types/auth";

const { Option } = Select;

const OrgUsers: React.FC = () => {
  const [users, setUsers] = useState<OrgUser[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [editingUser, setEditingUser] = useState<OrgUser | null>(null);

  const fetchUsers = async () => {
    try {
      setLoading(true);
      const response = await axios.post("/api/employer/get-org-users", {});
      setUsers(response.data);
    } catch (error) {
      message.error("Failed to fetch users");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleAdd = () => {
    form.resetFields();
    setEditingUser(null);
    setModalVisible(true);
  };

  const handleEdit = (record: OrgUser) => {
    form.setFieldsValue({
      ...record,
      password: undefined,
    });
    setEditingUser(record);
    setModalVisible(true);
  };

  const handleDelete = async (email: string) => {
    try {
      await axios.post("/api/employer/defunct-org-user", { email });
      message.success("User deleted successfully");
      fetchUsers();
    } catch (error) {
      message.error("Failed to delete user");
    }
  };

  const handleSubmit = async (values: any) => {
    try {
      if (editingUser) {
        await axios.post("/api/employer/update-org-user", {
          ...values,
          email: editingUser.email,
        });
        message.success("User updated successfully");
      } else {
        await axios.post("/api/employer/add-org-user", values);
        message.success("User added successfully");
      }
      setModalVisible(false);
      fetchUsers();
    } catch (error) {
      message.error("Failed to save user");
    }
  };

  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
    },
    {
      title: "Email",
      dataIndex: "email",
      key: "email",
    },
    {
      title: "Roles",
      dataIndex: "roles",
      key: "roles",
      render: (roles: string[]) => (
        <Space>
          {roles.map((role) => (
            <Tag key={role} color="blue">
              {role}
            </Tag>
          ))}
        </Space>
      ),
    },
    {
      title: "Status",
      dataIndex: "state",
      key: "state",
      render: (state: string) => (
        <Tag color={state === "ACTIVE_ORG_USER" ? "green" : "red"}>
          {state.replace("_ORG_USER", "").replace("_", " ")}
        </Tag>
      ),
    },
    {
      title: "Actions",
      key: "actions",
      render: (_: any, record: OrgUser) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => handleEdit(record)} />
          <Popconfirm
            title="Are you sure you want to delete this user?"
            onConfirm={() => handleDelete(record.email)}
            okText="Yes"
            cancelText="No"
          >
            <Button icon={<DeleteOutlined />} danger />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          Add User
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={users}
        loading={loading}
        rowKey="email"
      />

      <Modal
        title={editingUser ? "Edit User" : "Add User"}
        open={modalVisible}
        onOk={form.submit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="name"
            label="Name"
            rules={[{ required: true, message: "Please enter name" }]}
          >
            <Input />
          </Form.Item>

          {!editingUser && (
            <Form.Item
              name="email"
              label="Email"
              rules={[
                { required: true, message: "Please enter email" },
                { type: "email", message: "Please enter a valid email" },
              ]}
            >
              <Input />
            </Form.Item>
          )}

          {!editingUser && (
            <Form.Item
              name="password"
              label="Password"
              rules={[
                { required: true, message: "Please enter password" },
                { min: 8, message: "Password must be at least 8 characters" },
              ]}
            >
              <Input.Password prefix={<LockOutlined />} />
            </Form.Item>
          )}

          <Form.Item
            name="roles"
            label="Roles"
            rules={[
              { required: true, message: "Please select at least one role" },
            ]}
          >
            <Select mode="multiple">
              <Option value="ADMIN">Admin</Option>
              <Option value="RECRUITER">Recruiter</Option>
              <Option value="HIRING_MANAGER">Hiring Manager</Option>
              <Option value="INTERVIEWER">Interviewer</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default OrgUsers;
