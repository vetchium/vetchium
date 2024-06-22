import { PlusCircleTwoTone, SaveTwoTone } from "@ant-design/icons";
import { Button, Flex, Form, Input, InputNumber } from "antd";
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

  function validateDepartment(rule: any, value: string) {
    return new Promise<void>((resolve, reject) => {
      resolve();
    });
  }

  function validateHiringManager(rule: any, value: string) {
    return new Promise<void>((resolve, reject) => {
      resolve();
    });
  }

  function validatePositions(rule: any, value: number) {
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
        label={t("department")}
        name="department"
        rules={[{ validator: validateDepartment }]}
      >
        <Input />
      </Form.Item>
      <Form.Item
        label={t("hiring_manager")}
        name="hiringManager"
        rules={[{ validator: validateHiringManager }]}
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
