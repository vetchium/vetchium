import { Form, Input } from "antd";
import { formInputStyle, formItemStyle } from "../Styles";
import t from "../i18n/i18n";

function onFinish(values: any) {
  console.log("Received values:", values);
}

function onFinishFailed(errorInfo: any): void {
  console.log("Form validation failed:", errorInfo);
}

function validateInterviewers() {
  // In future, should get an array of email addresses and validate each ?
  return new Promise<void>((resolve, reject) => {
    resolve();
  });
}

export default function CreateInterviewForm() {
  return (
    <Form onFinish={onFinish} onFinishFailed={onFinishFailed} layout="vertical">
      <Form.Item
        label={t("create_interview.interviewers")}
        name="interviewers"
        rules={[{ required: true, validator: validateInterviewers }]}
        style={formItemStyle}
      >
        {/* In future should autocomplete from users via SSO etc. */}
        <Input style={formInputStyle} />
      </Form.Item>
    </Form>
  );
}
