import { Button, Form, Switch } from "antd";
import { t } from "i18next";
import { formItemStyle, formStyle } from "../Styles";

function InterviewCancelForm() {
  return (
    <Form style={formStyle}>
      <Form.Item
        label={t("interviews.cancel_notice")}
        style={{ textAlign: "left" }}
      >
        <Switch defaultChecked />
      </Form.Item>
      <Form.Item style={formItemStyle}>
        <Button type="primary" htmlType="submit" danger>
          {t("interviews.cancel_interview")}
        </Button>
      </Form.Item>
    </Form>
  );
}

export default InterviewCancelForm;
