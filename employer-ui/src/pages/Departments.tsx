import React, { useState, useEffect } from "react";
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Space,
  message,
  Popconfirm,
} from "antd";
import { PlusOutlined, EditOutlined, DeleteOutlined } from "@ant-design/icons";
import axios from "axios";
import { CostCenter } from "@/types/costCenter";

const Departments: React.FC = () => {
  const [departments, setDepartments] = useState<CostCenter[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [editingDepartment, setEditingDepartment] = useState<CostCenter | null>(
    null
  );

  const fetchDepartments = async () => {
    try {
      setLoading(true);
      const response = await axios.post("/api/employer/get-cost-centers", {});
      setDepartments(response.data);
    } catch (error) {
      message.error("Failed to fetch departments");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDepartments();
  }, []);

  const handleAdd = () => {
    form.resetFields();
    setEditingDepartment(null);
    setModalVisible(true);
  };

  const handleEdit = (record: CostCenter) => {
    form.setFieldsValue(record);
    setEditingDepartment(record);
    setModalVisible(true);
  };

  const handleDelete = async (name: string) => {
    try {
      await axios.post("/api/employer/defunct-cost-center", { name });
      message.success("Department deleted successfully");
      fetchDepartments();
    } catch (error) {
      message.error("Failed to delete department");
    }
  };

  const handleSubmit = async (values: any) => {
    try {
      if (editingDepartment) {
        await axios.post("/api/employer/update-cost-center", values);
        message.success("Department updated successfully");
      } else {
        await axios.post("/api/employer/add-cost-center", values);
        message.success("Department added successfully");
      }
      setModalVisible(false);
      fetchDepartments();
    } catch (error) {
      message.error("Failed to save department");
    }
  };

  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
    },
    {
      title: "Notes",
      dataIndex: "notes",
      key: "notes",
    },
    {
      title: "Status",
      dataIndex: "state",
      key: "state",
    },
    {
      title: "Actions",
      key: "actions",
      render: (_: any, record: CostCenter) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => handleEdit(record)} />
          <Popconfirm
            title="Are you sure you want to delete this department?"
            onConfirm={() => handleDelete(record.name)}
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
          Add Department
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={departments}
        loading={loading}
        rowKey="name"
      />

      <Modal
        title={editingDepartment ? "Edit Department" : "Add Department"}
        open={modalVisible}
        onOk={form.submit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="name"
            label="Name"
            rules={[
              { required: true, message: "Please enter department name" },
              { min: 3, message: "Name must be at least 3 characters" },
              { max: 64, message: "Name cannot exceed 64 characters" },
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="notes"
            label="Notes"
            rules={[
              { max: 1024, message: "Notes cannot exceed 1024 characters" },
            ]}
          >
            <Input.TextArea />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Departments;
