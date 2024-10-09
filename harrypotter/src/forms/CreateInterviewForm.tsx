import { PlusCircleTwoTone } from "@ant-design/icons";
import { Button, DatePicker, Flex, Form, Input, Select, Switch } from "antd";
import TextArea from "antd/es/input/TextArea";
import {
  formInputStyle,
  formItemStyle,
  formSelectStyle,
  formStyle,
} from "../Styles";
import t from "../i18n/i18n";
import { useState } from "react";

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

function validateDatePicker(value: any): Promise<void> {
  const now = new Date();
  const maxDate = new Date();
  maxDate.setDate(now.getDate() + 45);
  if (
    !value ||
    new Date(value).getTime() < now.getTime() ||
    new Date(value).getTime() > maxDate.getTime()
  ) {
    return Promise.reject(
      new Error(
        "Please select a valid date and time in the format YYYY-MM-DD HH:mm."
      )
    );
  }
  return Promise.resolve();
}

function validateInterviewAppointmentBody(value: any): Promise<void> {
  if (!value || value.trim() === "") {
    return Promise.reject(new Error("Interview appointment body is required."));
  }
  return Promise.resolve();
}

export default function CreateInterviewForm() {
  const [appointmentBodyDisabled, setAppointmentBodyDisabled] = useState(true);

  return (
    <Form
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
      layout="vertical"
      style={formStyle}
      initialValues={{ duration: "1h" }} // Set default value for duration
    >
      <Form.Item
        label={t("create_interview.interviewers")}
        name="interviewers"
        rules={[{ required: true, validator: validateInterviewers }]}
        style={formItemStyle}
      >
        {/* In future should autocomplete from users via SSO etc. */}
        <Input style={formInputStyle} />
      </Form.Item>

      <Form.Item
        label={t("create_interview.interview_time")}
        name="interviewTime"
        rules={[{ required: true, validator: validateDatePicker }]}
        style={formItemStyle}
      >
        <DatePicker
          showTime
          format="YYYY-MM-DD h:mm A"
          style={formInputStyle}
          minuteStep={15}
          disabledDate={(currentDate) =>
            currentDate.startOf("day").isBefore(new Date().setHours(0, 0, 0, 0))
          }
        />
      </Form.Item>

      <Form.Item
        label={t("create_interview.duration")}
        name="duration"
        rules={[{ required: true }]}
        style={formItemStyle}
      >
        <Select value="1h" style={formSelectStyle}>
          <Select.Option value="15m">15 minutes</Select.Option>
          <Select.Option value="30m">30 minutes</Select.Option>
          <Select.Option value="45m">45 minutes</Select.Option>
          <Select.Option value="1h">1 hour</Select.Option>
          <Select.Option value="1h15m">1 hour 15 minutes</Select.Option>
          <Select.Option value="1h30m">1 hour 30 minutes</Select.Option>
          <Select.Option value="1h45m">1 hour 45 minutes</Select.Option>
          <Select.Option value="2h">2 hours</Select.Option>
        </Select>
      </Form.Item>

      <Flex>
        <Form.Item
          label={t("create_interview.send_appointment")}
          name="sendAppointment"
          style={formItemStyle}
        >
          <Switch
            onChange={(checked) => setAppointmentBodyDisabled(!checked)}
          />
        </Form.Item>

        <Form.Item
          label={t("create_interview.appointment_body")}
          name="interviewAppointmentBody"
          rules={[{ validator: validateInterviewAppointmentBody }]}
          style={formItemStyle}
        >
          <TextArea
            rows={12}
            style={formInputStyle}
            disabled={appointmentBodyDisabled}
          />
        </Form.Item>
      </Flex>

      <Form.Item>
        <Button type="primary" htmlType="submit" icon={<PlusCircleTwoTone />}>
          {t("create_interview.create_interview")}
        </Button>
      </Form.Item>
    </Form>
  );
}
