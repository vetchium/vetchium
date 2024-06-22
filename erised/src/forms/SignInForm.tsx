import { Button, Form, Input } from "antd";
import t from "../i18n/i18n";
import { formStyle } from "../Styles";

function SignInForm({ onSignIn }: { onSignIn: () => void }) {
  function onFinish() {
    onSignIn();
  }

  function validateDomain(rule: any, value: string) {
    return new Promise<void>((resolve, reject) => {
      if (!value || value !== "example.com") {
        reject(t("invalid_field"));
      }

      resolve();
    });
  }

  function validateEmail(rule: any, value: string) {
    return new Promise<void>((resolve, reject) => {
      if (!value || !value.includes("@")) {
        reject(t("invalid_field"));
      }

      resolve();
    });
  }

  function validatePassword(rule: any, value: string) {
    return new Promise<void>((resolve, reject) => {
      if (!value || value.length < 8) {
        reject(t("Password must be at least 8 characters long"));
      }

      resolve();
    });
  }

  function onFinishFailed(errorInfo: any): void {
    console.log("Form validation failed:", errorInfo);
  }

  return (
    <Form onFinish={onFinish} onFinishFailed={onFinishFailed} style={formStyle}>
      <Form.Item
        label={t("company_domain")}
        name="domain"
        rules={[{ required: true, validator: validateDomain }]}
        initialValue="example.com"
      >
        <Input />
      </Form.Item>
      <Form.Item
        label={t("email_address")}
        name="email"
        rules={[{ required: true, validator: validateEmail, type: "email" }]}
        initialValue="master@example.com"
      >
        <Input />
      </Form.Item>
      <Form.Item
        label={t("password")}
        name="password"
        rules={[{ required: true, validator: validatePassword }]}
        initialValue="Password123$"
      >
        <Input.Password />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          {t("sign_in")}
        </Button>
      </Form.Item>
    </Form>
  );
}

export default SignInForm;
