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
} from "antd";
import { PlusOutlined, EditOutlined, DeleteOutlined } from "@ant-design/icons";
import axios from "axios";
import { Location } from "@/types/location";

const { Option } = Select;

const Locations: React.FC = () => {
  const [locations, setLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [editingLocation, setEditingLocation] = useState<Location | null>(null);

  const fetchLocations = async () => {
    try {
      setLoading(true);
      const response = await axios.post("/api/employer/get-locations", {});
      setLocations(response.data);
    } catch (error) {
      message.error("Failed to fetch locations");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLocations();
  }, []);

  const handleAdd = () => {
    form.resetFields();
    setEditingLocation(null);
    setModalVisible(true);
  };

  const handleEdit = (record: Location) => {
    form.setFieldsValue(record);
    setEditingLocation(record);
    setModalVisible(true);
  };

  const handleDelete = async (title: string) => {
    try {
      await axios.post("/api/employer/defunct-location", { title });
      message.success("Location deleted successfully");
      fetchLocations();
    } catch (error) {
      message.error("Failed to delete location");
    }
  };

  const handleSubmit = async (values: any) => {
    try {
      if (editingLocation) {
        await axios.post("/api/employer/update-location", values);
        message.success("Location updated successfully");
      } else {
        await axios.post("/api/employer/add-location", values);
        message.success("Location added successfully");
      }
      setModalVisible(false);
      fetchLocations();
    } catch (error) {
      message.error("Failed to save location");
    }
  };

  const columns = [
    {
      title: "Title",
      dataIndex: "title",
      key: "title",
    },
    {
      title: "Country",
      dataIndex: "country_code",
      key: "country_code",
    },
    {
      title: "Address",
      dataIndex: "postal_address",
      key: "postal_address",
    },
    {
      title: "Postal Code",
      dataIndex: "postal_code",
      key: "postal_code",
    },
    {
      title: "Actions",
      key: "actions",
      render: (_: any, record: Location) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => handleEdit(record)} />
          <Popconfirm
            title="Are you sure you want to delete this location?"
            onConfirm={() => handleDelete(record.title)}
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
          Add Location
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={locations}
        loading={loading}
        rowKey="title"
      />

      <Modal
        title={editingLocation ? "Edit Location" : "Add Location"}
        open={modalVisible}
        onOk={form.submit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="title"
            label="Title"
            rules={[{ required: true, message: "Please enter location title" }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="country_code"
            label="Country Code"
            rules={[{ required: true, message: "Please select country code" }]}
          >
            <Select>
              <Option value="USA">USA</Option>
              <Option value="GBR">GBR</Option>
              <Option value="IND">IND</Option>
              {/* Add more country options */}
            </Select>
          </Form.Item>
          <Form.Item
            name="postal_address"
            label="Address"
            rules={[{ required: true, message: "Please enter address" }]}
          >
            <Input.TextArea />
          </Form.Item>
          <Form.Item
            name="postal_code"
            label="Postal Code"
            rules={[{ required: true, message: "Please enter postal code" }]}
          >
            <Input />
          </Form.Item>
          <Form.Item name="openstreetmap_url" label="OpenStreetMap URL">
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Locations;
