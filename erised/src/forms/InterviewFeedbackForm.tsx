import { Button, Divider, Form, Radio } from "antd";
import { formItemStyle, formRadioStyle, formStyle } from "../Styles";
import t from "../i18n/i18n";
import TextArea from "antd/lib/input/TextArea";

function onFinish(values: any) {
  console.log(values);
}

function onFinishFailed(errorInfo: any) {
  console.log(errorInfo);
}

function validateResult(rule: any, value: any) {
  if (value === "") {
    return Promise.reject("Please select a result");
  }
  return Promise.resolve();
}

function validatePositives(rule: any, value: any) {
  if (value === "") {
    return Promise.reject("Please enter positives");
  }
  return Promise.resolve();
}

function validateNegatives(rule: any, value: any) {
  if (value === "") {
    return Promise.reject("Please enter negatives");
  }
  return Promise.resolve();
}

function validateSummary(rule: any, value: any) {
  if (value === "") {
    return Promise.reject("Please enter summary");
  }
  return Promise.resolve();
}

function InterviewFeedbackForm() {
  return (
    <Form
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
      style={formStyle}
      layout="vertical"
    >
      <Form.Item
        label={t("interview_feedback_form.positives")}
        name="positives"
        rules={[{ required: true, validator: validatePositives }]}
        style={formItemStyle}
      >
        <TextArea
          rows={4}
          placeholder={t("interview_feedback_form.positives_placeholder")}
        />
      </Form.Item>
      <Form.Item
        label={t("interview_feedback_form.negatives")}
        name="negatives"
        rules={[{ required: true, validator: validateNegatives }]}
        style={formItemStyle}
      >
        <TextArea
          rows={4}
          placeholder={t("interview_feedback_form.negatives_placeholder")}
        />
      </Form.Item>
      <Form.Item
        label={t("interview_feedback_form.summary")}
        name="summary"
        rules={[{ required: true, validator: validateSummary }]}
        style={formItemStyle}
      >
        <TextArea
          rows={4}
          placeholder={t("interview_feedback_form.summary_placeholder")}
        />
      </Form.Item>
      <Form.Item
        label={t("interview_feedback_form.result")}
        name="result"
        rules={[{ required: true, validator: validateResult }]}
        style={formItemStyle}
      >
        <Radio.Group
          defaultValue="YES"
          buttonStyle="solid"
          style={formRadioStyle}
        >
          <Radio.Button value="STRONG_YES">
            {t("interview_feedback_form.strong_yes")}
          </Radio.Button>
          <Radio.Button value="YES">
            {t("interview_feedback_form.yes")}
          </Radio.Button>
          <Radio.Button value="NO">
            {t("interview_feedback_form.no")}
          </Radio.Button>
          <Radio.Button value="STRONG_NO">
            {t("interview_feedback_form.strong_no")}
          </Radio.Button>
        </Radio.Group>
      </Form.Item>
      <Divider />

      <Form.Item
        label={t("interview_feedback_form.candidate_feedback")}
        name="candidateFeedback"
        rules={[]}
        style={formItemStyle}
      >
        <TextArea
          rows={4}
          placeholder={t(
            "interview_feedback_form.candidate_feedback_placeholder"
          )}
        />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit">
          {t("interview_feedback_form.submit_feedback")}
        </Button>
      </Form.Item>
    </Form>
  );
}

export default InterviewFeedbackForm;
