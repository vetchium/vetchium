import { DeleteTwoTone } from "@ant-design/icons";
import { Button, Divider, Flex, Form, Input, Table } from "antd";
import TextArea from "antd/lib/input/TextArea";
import React from "react";
import { formItemStyle, formStyle, tableStyle } from "../Styles";
import t from "../i18n/i18n";

const validateDepartmentName = (_: any, value: string) => {
  if (!value || value.length <= 100) {
    return Promise.resolve();
  }
  return Promise.reject(new Error(t("department_name_too_long")));
};

const validateNotes = (_: any, value: string) => {
  if (!value || value.length <= 500) {
    return Promise.resolve();
  }
  return Promise.reject(new Error(t("notes_too_long")));
};

const Departments: React.FC = () => {
  const data = [
    {
      key: 1,
      name: "APAC Sales",
      notes: "Covers for Japan, Singapore, Korea, Taiwan, China",
    },
    {
      key: 2,
      name: "Legal",
      notes: "Compliance, Audits, Litigation",
    },
    {
      key: 3,
      name: "Finance and HR",
      notes: "The best department in the company",
    },
    {
      key: 4,
      name: "Global Engineering",
    },
  ];

  const handleAddDepartment = (values: { name: string; notes: string }) => {};

  const columns = [
    { title: t("departments.department_name"), dataIndex: "name", key: "name" },
    { title: t("departments.notes"), dataIndex: "notes", key: "notes" },
    {
      title: t("actions"),
      key: "actions",
      render: (text: string, record: any) => (
        <span>
          <Button icon={<DeleteTwoTone />} />
        </span>
      ),
    },
  ];

  return (
    <Flex wrap>
      <Table dataSource={data} columns={columns} style={tableStyle} />
      <Divider>{t("departments.add_department")}</Divider>
      <Form
        onFinish={handleAddDepartment}
        initialValues={{ name: "", notes: "" }}
        style={formStyle}
        layout="vertical"
      >
        <Form.Item
          label={t("departments.department_name")}
          name="name"
          rules={[{ required: true }, { validator: validateDepartmentName }]}
          style={formItemStyle}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label={t("departments.department_notes")}
          name="notes"
          rules={[{ validator: validateNotes }]}
          style={formItemStyle}
        >
          <TextArea />
        </Form.Item>
        <Divider />
        <Flex gap="middle">
          <Form.Item style={formItemStyle}>
            <Button type="primary" htmlType="submit">
              {t("departments.add_department")}
            </Button>
          </Form.Item>
          <Form.Item style={formItemStyle}>
            <Button htmlType="reset">{t("common.reset")}</Button>
          </Form.Item>
        </Flex>
      </Form>
      <Divider />
    </Flex>
  );
};

export default Departments;
