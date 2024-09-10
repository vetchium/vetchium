import { Button, Form, Switch } from "antd";
import { t } from "i18next";
import { formStyle } from "../Styles";

function InterviewCancelForm() {
  return (
    <Form layout="vertical" style={formStyle}>
      <Form.Item label={t("interviews.cancel_notice")}>
        <Switch defaultChecked />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit" danger>
          {t("interviews.cancel_interview")}
        </Button>
      </Form.Item>
    </Form>
  );
}

export default InterviewCancelForm;
