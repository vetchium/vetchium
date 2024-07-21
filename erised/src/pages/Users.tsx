import { DeleteTwoTone } from "@ant-design/icons";
import {
  Button,
  Divider,
  Flex,
  Form,
  Input,
  Radio,
  Table,
  Typography,
} from "antd";
import React, { useState } from "react";
import { formItemStyle, formStyle, tableStyle } from "../Styles";
import t from "../i18n/i18n";

const Users: React.FC = () => {
  const [users, setUsers] = useState([
    { key: 1, email: "john.doe@example.com", roles: ["Admin"] },
    { key: 2, email: "jane.smith@example.com", roles: ["Recruiter"] },
    { key: 3, email: "alice.johnson@example.com", roles: ["Panelist"] },
    {
      key: 4,
      email: "bob.brown@example.com",
      roles: ["Recruiter", "Panelist"],
    },
  ]);

  const handleAddUser = (values: { email: string; roles: string[] }) => {
    const newUser = {
      key: users.length + 1,
      email: values.email,
      roles: values.roles,
    };
    setUsers([...users, newUser]);
  };

  const handleDeleteUser = (key: number) => {
    setUsers(users.filter((user) => user.key !== key));
  };

  const columns = [
    { title: t("email"), dataIndex: "email", key: "email" },
    {
      title: t("roles"),
      dataIndex: "roles",
      key: "roles",
      render: (roles: string[]) => roles.join(", "),
    },
    {
      title: t("actions"),
      key: "actions",
      render: (text: string, record: any) =>
        // TODO: Should enable the Delete icon only if the logged in user is and Admin
        record.roles.includes("Admin") ? null : (
          <span>
            <Button
              icon={<DeleteTwoTone />}
              onClick={() => handleDeleteUser(record.key)}
            />
          </span>
        ),
    },
  ];

  return (
    <Flex wrap>
      <Divider>{t("add_user")}</Divider>
      <Form
        onFinish={handleAddUser}
        initialValues={{ email: "", roles: [] }}
        style={formStyle}
      >
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
          <Radio.Group defaultValue={"panelist"} buttonStyle="solid">
            <Radio.Button value="admin">{t("admin")}</Radio.Button>
            <Radio.Button value="recruiter">{t("recruiter")}</Radio.Button>
            <Radio.Button value="panelist">{t("panelist")}</Radio.Button>
          </Radio.Group>
        </Form.Item>
        <Divider />
        <Flex gap="middle">
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
      <Divider />
      <Table dataSource={users} columns={columns} style={tableStyle} />
      <Typography>
        <Typography.Title level={3}>{t("role_descriptions")}</Typography.Title>
        <Typography.Paragraph>
          <strong>{t("admin")}:</strong> {t("admin_description")}
        </Typography.Paragraph>
        <Typography.Paragraph>
          <strong>{t("recruiter")}:</strong> {t("recruiter_description")}
        </Typography.Paragraph>
        <Typography.Paragraph>
          <strong>{t("panelist")}:</strong> {t("panelist_description")}
        </Typography.Paragraph>
      </Typography>
      <Divider />
    </Flex>
  );
};

export default Users;
