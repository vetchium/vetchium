import { Button, Form, Input } from "antd";
import { useNavigate } from "react-router-dom";
import Router from "../components/Router";
import t from "../i18n/i18n";
import { formStyle } from "../Styles";

function SignInForm() {
  const navigate = useNavigate();

  function onFinish(values: any) {
    console.log("Received values:", values);
    navigate("/openings");
  }

  function validateDomain(rule: any, value: string) {
    // Add your domain validation logic here
    // Example: Check if the domain is valid
    return new Promise<void>((resolve, reject) => {
      if (!value || value !== "example.com") {
        reject(t("Invalid domain"));
      } else {
        resolve();
      }
    });
  }

  function validateEmail(rule: any, value: string) {
    // Add your email validation logic here
    return new Promise<void>((resolve, reject) => {
      if (!value || !value.includes("@")) {
        reject(t("Invalid email"));
      } else {
        resolve();
      }
    });
  }

  function validatePassword(rule: any, value: string) {
    // Add your password validation logic here
    // Example: Check if the password meets certain criteria
    return new Promise<void>((resolve, reject) => {
      if (!value || value.length < 8) {
        reject(t("Password must be at least 8 characters long"));
      } else {
        resolve();
      }
    });
  }

  function onFinishFailed(errorInfo: any): void {
    console.log("Form validation failed:", errorInfo);
  }

  return (
    <>
      <Router />
      <Form
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        style={formStyle}
      >
        <Form.Item
          label={t("company_domain")}
          name="domain"
          rules={[{ required: true, validator: validateDomain }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label={t("email_address")}
          name="email"
          rules={[{ required: true, validator: validateEmail, type: "email" }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label={t("password")}
          name="password"
          rules={[{ required: true, validator: validatePassword }]}
        >
          <Input.Password />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">
            {t("sign_in")}
          </Button>
        </Form.Item>
      </Form>
    </>
  );
}

export default SignInForm;
