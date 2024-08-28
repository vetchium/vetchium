import { DeleteTwoTone } from "@ant-design/icons";
import { Button, Flex, Form, Input, Radio, Table, Typography } from "antd";
import React, { useState } from "react";
import { formItemStyle, formStyle, tableStyle } from "../Styles";
import t from "../i18n/i18n";

const Users: React.FC = () => {
  const [users, setUsers] = useState([
    { key: 1, email: "john.doe@example.com", role: "Admin" },
    { key: 2, email: "jane.smith@example.com", role: "Recruiter" },
    { key: 3, email: "alice.johnson@example.com", role: "Interviewer" },
  ]);

  const handleAddUser = (values: { email: string; role: string }) => {
    const newUser = {
      key: users.length + 1,
      email: values.email,
      role: values.role,
    };
    setUsers([...users, newUser]);
  };

  const handleDeleteUser = (key: number) => {};

  const columns = [
    { title: t("email"), dataIndex: "email", key: "email" },
    {
      title: t("role"),
      dataIndex: "role",
      key: "role",
    },
    {
      title: t("actions"),
      key: "actions",
      render: (text: string, record: any) => (
        // TODO: Should enable the Delete icon only if the logged in user is and Admin
        <Button
          icon={<DeleteTwoTone />}
          onClick={() => handleDeleteUser(record.key)}
        />
      ),
    },
  ];

  return (
    <Flex wrap vertical>
      <Form onFinish={handleAddUser} style={formStyle} layout="vertical">
        <Form.Item
          label={t("email")}
          name="email"
          rules={[{ required: true, type: "email" }]}
          style={formItemStyle}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label={t("role")}
          name="role"
          rules={[{ required: true }]}
          style={formItemStyle}
        >
          {/* In future, we may have to accomodate roles that are not hierarchical */}
          <Radio.Group defaultValue={"interviewer"} buttonStyle="solid">
            <Radio.Button value="admin">{t("admin")}</Radio.Button>
            <Radio.Button value="recruiter">{t("recruiter")}</Radio.Button>
            <Radio.Button value="interviewer">{t("interviewer")}</Radio.Button>
          </Radio.Group>
          <Typography style={formItemStyle}>
            <Typography.Title level={5}>
              {t("role_descriptions")}
            </Typography.Title>
            <Typography.Paragraph>
              <strong>{t("admin")}:</strong> {t("admin_description")}
            </Typography.Paragraph>
            <Typography.Paragraph>
              <strong>{t("recruiter")}:</strong> {t("recruiter_description")}
            </Typography.Paragraph>
            <Typography.Paragraph>
              <strong>{t("interviewer")}:</strong>{" "}
              {t("interviewer_description")}
            </Typography.Paragraph>
          </Typography>
        </Form.Item>
        <Flex wrap gap="large">
          <Form.Item style={formItemStyle}>
            <Button type="primary" htmlType="submit">
              {t("add_user")}
            </Button>
          </Form.Item>
          <Form.Item style={formItemStyle}>
            <Button htmlType="reset">{t("reset")}</Button>
          </Form.Item>
        </Flex>
      </Form>
      <Table dataSource={users} columns={columns} style={tableStyle} />
    </Flex>
  );
};

export default Users;
