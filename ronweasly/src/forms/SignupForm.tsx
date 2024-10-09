import React from "react";
import { Form, Input, Button } from "antd";
import { useTranslation } from "react-i18next";

const SignupForm: React.FC = () => {
  const { t } = useTranslation();
  const [form] = Form.useForm();

  const onFinish = (values: any) => {
    console.log("Signup form values:", values);
    // Handle signup logic here
  };

  return (
    <Form form={form} name="signup" onFinish={onFinish} layout="vertical">
      <Form.Item
        name="email"
        label="Email"
        rules={[
          {
            required: true,
            type: "email",
            message: "Please enter a valid email",
          },
        ]}
      >
        <Input />
      </Form.Item>
      <Form.Item
        name="password"
        label="Password"
        rules={[{ required: true, message: "Please enter your password" }]}
      >
        <Input.Password />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          {t("signup.submit")}
        </Button>
      </Form.Item>
    </Form>
  );
};

export default SignupForm;
