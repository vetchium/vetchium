import { PlusCircleTwoTone, SaveTwoTone } from "@ant-design/icons";
import {
  Button,
  Divider,
  Flex,
  Form,
  Input,
  InputNumber,
  Select,
  Slider,
} from "antd";
import TextArea from "antd/es/input/TextArea";
import t from "../i18n/i18n";
import { formStyle } from "../Styles";

function CreateOpeningForm() {
  function onFinish(values: any) {
    console.log("Received values:", values);
  }

  function onFinishFailed(errorInfo: any): void {
    console.log("Form validation failed:", errorInfo);
  }

  function validateTitle(rule: any, value: string) {
    return new Promise<void>((resolve, reject) => {
      if (!value || value.length < 3) {
        reject(t("invalid_field"));
      }

      resolve();
    });
  }

  function validatePositions(rule: any, value: number) {
    return new Promise<void>((resolve, reject) => {
      resolve();
    });
  }

  function validateJD(rule: any, value: string) {
    return new Promise<void>((resolve, reject) => {
      resolve();
    });
  }

  function validateLocations(rule: any, value: string) {
    return new Promise<void>((resolve, reject) => {
      resolve();
    });
  }

  function validateYOE(rule: any, value: number) {
    return new Promise<void>((resolve, reject) => {
      resolve();
    });
  }

  function validateHiringManager(rule: any, value: string) {
    return new Promise<void>((resolve, reject) => {
      resolve();
    });
  }

  function validateCurrency(rule: any, value: string) {
    return new Promise<void>((resolve, reject) => {
      resolve();
    });
  }

  function validateSalaryMin(rule: any, value: number) {
    return new Promise<void>((resolve, reject) => {
      resolve();
    });
  }

  function validateSalaryMax(rule: any, value: number) {
    return new Promise<void>((resolve, reject) => {
      resolve();
    });
  }

  function validateDepartment(rule: any, value: string) {
    return new Promise<void>((resolve, reject) => {
      resolve();
    });
  }

  return (
    <Form onFinish={onFinish} onFinishFailed={onFinishFailed} style={formStyle}>
      <Form.Item
        label={t("job_title")}
        name="title"
        rules={[{ required: true, validator: validateTitle }]}
      >
        <Input />
      </Form.Item>

      <Form.Item
        label={t("positions")}
        name="positions"
        initialValue={1}
        rules={[{ required: true, validator: validatePositions }]}
      >
        <InputNumber min={1} max={25} />
      </Form.Item>

      <Form.Item
        label={t("jd")}
        name="jd"
        rules={[{ required: true, validator: validateJD }]}
      >
        <TextArea placeholder="Job Description" rows={10} />
      </Form.Item>

      <Form.Item
        label={t("locations")}
        rules={[{ required: true, validator: validateLocations }]}
      >
        <Select
          mode="tags"
          placeholder={t("locations")}
          style={{ minWidth: "120px" }}
        >
          {/* Should fetch from API based on the company */}
          <Select.Option value="global">Global</Select.Option>
          <Select.Option value="bangalore">Bangalore</Select.Option>
          <Select.Option value="chennai">Chennai</Select.Option>
          <Select.Option value="san francisco">San Francisco</Select.Option>
          <Select.Option value="germany">Germany</Select.Option>
          <Select.Option value="europe remote">Europe Remote</Select.Option>
        </Select>
      </Form.Item>

      <Divider>{t("optional_fields")}</Divider>

      <Form.Item
        label={t("yoe")}
        name="yoe"
        rules={[{ validator: validateYOE }]}
      >
        <Slider
          min={0}
          max={80}
          step={5}
          range={true}
          defaultValue={[0, 10]}
          style={{ minWidth: "300px" }}
          marks={{
            0: "0",
            10: "10",
            20: "20",
            30: "30",
            40: "40",
            50: "50",
            60: "60",
            70: "70",
            80: "80",
          }}
        />
      </Form.Item>

      <Form.Item
        label={t("hiring_manager")}
        name="hiringManager"
        rules={[{ validator: validateHiringManager }]}
      >
        <Input />
      </Form.Item>

      <Form.Item
        label={t("currency")}
        name="currency"
        rules={[{ validator: validateCurrency }]}
      >
        {/* Should fetch from API based on the job location */}
        <Select>
          <Select.Option value="USD">USD</Select.Option>
          <Select.Option value="INR">INR</Select.Option>
          <Select.Option value="EUR">EUR</Select.Option>
        </Select>
      </Form.Item>

      <Form.Item
        label={t("salary_min")}
        name="salarymin"
        rules={[{ validator: validateSalaryMin }]}
      >
        <InputNumber></InputNumber>
      </Form.Item>

      <Form.Item
        label={t("salary_max")}
        name="salarymax"
        rules={[{ validator: validateSalaryMax }]}
      >
        <InputNumber></InputNumber>
      </Form.Item>

      <Divider>{t("private_fields")}</Divider>

      <Form.Item
        label={t("department")}
        name="department"
        rules={[{ validator: validateDepartment }]}
      >
        <Input />
      </Form.Item>

      <Form.Item label={t("notes")} name="notes">
        <TextArea />
      </Form.Item>

      <Divider />

      <Flex gap="middle">
        <Form.Item>
          <Button type="primary" icon={<PlusCircleTwoTone />} htmlType="submit">
            {t("create_opening")}
          </Button>
        </Form.Item>
        <Flex gap="middle" justify="flex-end">
          <Form.Item>
            <Button>{t("cancel")}</Button>
          </Form.Item>
          <Form.Item>
            <Button icon={<SaveTwoTone />}>{t("save_draft")}</Button>
          </Form.Item>
        </Flex>
      </Flex>
    </Form>
  );
}

export default CreateOpeningForm;
