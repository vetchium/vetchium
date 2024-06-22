import { Button, Flex, Form, Input, InputNumber } from "antd";
import { PlusCircleTwoTone } from "@ant-design/icons";
import t from "../i18n/i18n";
import { formStyle } from "../Styles";

function CreateOpeningForm() {
  return (
    <Form style={formStyle}>
      <Form.Item label={t("job_title")} name="title">
        <Input />
      </Form.Item>
      <Form.Item label={t("department")} name="department">
        <Input />
      </Form.Item>
      <Form.Item label={t("hiring_manager")} name="hiringManager">
        <Input />
      </Form.Item>
      <Form.Item label={t("positions")} name="positions">
        <InputNumber min={1} max={25} defaultValue={1} />
      </Form.Item>
      <Flex gap="middle">
        <Form.Item>
          <Button type="primary" htmlType="submit" icon={<PlusCircleTwoTone />}>
            {t("create_opening")}
          </Button>
        </Form.Item>
        <Flex gap="middle" justify="flex-end">
          <Form.Item>
            <Button>{t("cancel")}</Button>
          </Form.Item>
          <Form.Item>
            <Button>{t("save_draft")}</Button>
          </Form.Item>
        </Flex>
      </Flex>
    </Form>
  );
}

export default CreateOpeningForm;
