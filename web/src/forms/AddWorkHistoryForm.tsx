import {
  AutoComplete,
  Button,
  DatePicker,
  Flex,
  Form,
  Input,
  Switch,
} from "antd";
import {
  formInputStyle,
  formItemStyle,
  formStyle,
  formSwitchStyle,
  formTextAreaStyle,
} from "../Styles";
import t from "../i18n/i18n";
import { useNavigate } from "react-router-dom";

function onFinish(values: any) {
  console.log(values);
}

function onFinishFailed(errorInfo: any): void {
  console.log("Form validation failed:", errorInfo);
}

function validateCompanyName(rule: any, value: any) {
  if (value.length < 3) {
    return Promise.reject(
      new Error("Company name must be at least 3 characters")
    );
  }
  return Promise.resolve();
}

function validateJobTitle(rule: any, value: any) {
  if (value.length < 3) {
    return Promise.reject(new Error("Job title must be at least 3 characters"));
  }
  return Promise.resolve();
}

function validateStartDate(rule: any, value: any) {
  if (!value) {
    return Promise.reject(new Error("Start date is required"));
  }
  return Promise.resolve();
}

function validateEndDate(rule: any, value: any) {
  if (value && value.isBefore(rule.startDate)) {
    return Promise.reject(new Error("End date must be after start date"));
  }
  return Promise.resolve();
}

function AddWorkHistoryForm() {
  const navigate = useNavigate();

  return (
    <Form
      layout="vertical"
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
      style={formStyle}
    >
      <Form.Item
        label={t("add_work_history_form.company_name")}
        name="company_name"
        rules={[{ required: true, validator: validateCompanyName }]}
        style={formItemStyle}
      >
        <AutoComplete
          style={formInputStyle}
          options={[
            {
              text: "Example Private Limited",
              value: "Example Private Limited",
            },
            { text: "Example Inc", value: "Example Inc" },
            { text: "Example Gmbh", value: "Example Gmbh" },
          ]}
          onSearch={(searchText) => {
            return [
              {
                text: "Example Private Limited",
                value: "Example Private Limited",
              },
              { text: "Example Inc", value: "Example Inc" },
              { text: "Example Gmbh", value: "Example Gmbh" },
            ].filter((option) =>
              option.text.toLowerCase().includes(searchText.toLowerCase())
            );
          }}
        />
      </Form.Item>

      <Form.Item
        label={t("add_work_history_form.job_title")}
        name="job_title"
        rules={[{ required: true, validator: validateJobTitle }]}
        style={formItemStyle}
      >
        <Input style={formInputStyle} />
      </Form.Item>

      <Form.Item
        label={t("add_work_history_form.start_date")}
        name="start_date"
        rules={[{ required: true, validator: validateStartDate }]}
        style={formItemStyle}
      >
        <DatePicker style={formInputStyle} />
      </Form.Item>

      <Flex gap="large">
        <Form.Item
          label={t("add_work_history_form.still_employed")}
          name="still_employed"
          rules={[{ required: false }]}
        >
          <Switch style={formSwitchStyle} />
        </Form.Item>
        <Form.Item
          label={t("add_work_history_form.end_date")}
          name="end_date"
          rules={[{ required: false, validator: validateEndDate }]}
        >
          <DatePicker style={formInputStyle} />
        </Form.Item>
      </Flex>

      <Form.Item
        label={t("add_work_history_form.description")}
        name="description"
        rules={[{ required: false }]}
        style={formItemStyle}
      >
        <Input.TextArea style={formTextAreaStyle} />
      </Form.Item>

      <Flex justify="space-between">
        <Form.Item>
          <Button type="primary" htmlType="submit">
            {t("add_work_history_form.submit")}
          </Button>
        </Form.Item>

        <Form.Item>
          <Button
            htmlType="button"
            onClick={() => {
              navigate("/home");
            }}
          >
            {t("add_work_history_form.cancel")}
          </Button>
        </Form.Item>
      </Flex>
    </Form>
  );
}

export default AddWorkHistoryForm;
