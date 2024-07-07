import { PlusCircleTwoTone, SaveTwoTone } from "@ant-design/icons";
import {
  Button,
  Divider,
  Flex,
  Form,
  Input,
  InputNumber,
  Radio,
  Select,
  Slider,
  Switch,
} from "antd";
import TextArea from "antd/es/input/TextArea";
import React, { useState } from "react";
import t from "../i18n/i18n";
import countriesData from "../static/countries-states-cities.json";
import { timezones } from "../static/timezones";
import {
  createOpeningFormStyle,
  formInputStyle,
  formItemStyle,
  formSelectStyle,
} from "../Styles";

function CreateOpeningForm() {
  const [isTimezoneSwitchOn, setIsTimezoneSwitchOn] = useState(false);
  const [isCountrySwitchOn, setIsCountrySwitchOn] = useState(false);

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

  function validateOnSiteLocations(rule: any, value: string) {
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
    <Form
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
      style={createOpeningFormStyle}
    >
      <Form.Item
        label={t("job_title")}
        name="title"
        rules={[{ required: true, validator: validateTitle }]}
        style={formItemStyle}
      >
        <Input style={formInputStyle} />
      </Form.Item>
      <Form.Item
        label={t("positions")}
        name="positions"
        initialValue={1}
        rules={[{ required: true, validator: validatePositions }]}
        style={formItemStyle}
      >
        <InputNumber min={1} max={25} style={formInputStyle} />
      </Form.Item>
      <Form.Item
        label={t("jd")}
        name="jd"
        rules={[{ required: true, validator: validateJD }]}
        style={formItemStyle}
      >
        <TextArea
          placeholder="Job Description"
          rows={10}
          style={formInputStyle}
        />
      </Form.Item>

      <Form.Item label={t("job_type")} name="jobType" style={formItemStyle}>
        <Radio.Group defaultValue={"full_time"} buttonStyle="solid">
          <Radio.Button value="full_time">{t("full_time")}</Radio.Button>
          <Radio.Button value="part_time">{t("part_time")}</Radio.Button>
          <Radio.Button value="contract">{t("contract")}</Radio.Button>
          <Radio.Button value="internship">{t("internship")}</Radio.Button>
          <Radio.Button value="unspecified">{t("unspecified")}</Radio.Button>
        </Radio.Group>
      </Form.Item>

      {/* <!--------- Location Fields ---------!> */}
      <Divider>{t("locations")}</Divider>
      <Form.Item
        label={t("onsite_locations")}
        rules={[{ required: true, validator: validateOnSiteLocations }]}
        style={formItemStyle}
      >
        <Select
          mode="tags"
          placeholder={t("locations")}
          style={formSelectStyle}
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

      <Form.Item
        label={t("remote_locations_countries")}
        name="remoteLocationsCountries"
        style={formItemStyle}
      >
        <Flex gap="small" vertical>
          <Switch
            style={{ maxWidth: "40px" }}
            checked={isCountrySwitchOn}
            onChange={(checked) => setIsCountrySwitchOn(checked)}
          />
          <Select
            mode="tags"
            placeholder={t("remote_locations_countries")}
            style={formSelectStyle}
            disabled={!isCountrySwitchOn}
          >
            {Array.isArray(countriesData) &&
              countriesData.map((country: any) => (
                <Select.Option key={country.name} value={country.name}>
                  {country.name}{" "}
                  {country.native === country.name ? "" : `(${country.native})`}
                </Select.Option>
              ))}
          </Select>
        </Flex>
      </Form.Item>

      <Form.Item
        label={t("remote_locations_timezones")}
        name="remoteLocationsTimezones"
        style={formItemStyle}
      >
        <Flex gap="small" vertical>
          <Switch
            style={{ maxWidth: "40px" }}
            checked={isTimezoneSwitchOn}
            onChange={(checked) => setIsTimezoneSwitchOn(checked)}
          />
          <Select
            mode="tags"
            placeholder={t("remote_locations_timezones")}
            style={formSelectStyle}
            disabled={!isTimezoneSwitchOn}
          >
            {timezones.map((timezone) => (
              <Select.Option key={timezone} value={timezone}>
                {timezone}
              </Select.Option>
            ))}
          </Select>
        </Flex>
      </Form.Item>

      {/* <!--------- Optional Fields ---------!> */}
      <Divider>{t("optional_fields")}</Divider>
      <Form.Item
        label={t("yoe")}
        name="yoe"
        rules={[{ validator: validateYOE }]}
        style={formItemStyle}
      >
        <Slider
          min={0}
          max={80}
          step={5}
          range={true}
          defaultValue={[0, 80]}
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
        label={t("educational_qualification")}
        name="educationalQualification"
        style={formItemStyle}
      >
        <Radio.Group defaultValue="unspecified" buttonStyle="solid">
          <Radio.Button value="bachelors">{t("bachelors")}</Radio.Button>
          <Radio.Button value="masters">{t("masters")}</Radio.Button>
          <Radio.Button value="phd">{t("phd")}</Radio.Button>
          <Radio.Button value="doesnt_matter">
            {t("doesnt_matter")}
          </Radio.Button>
          <Radio.Button value="unspecified">{t("unspecified")}</Radio.Button>
        </Radio.Group>
      </Form.Item>
      <Form.Item
        label={t("hiring_manager")}
        name="hiringManager"
        rules={[{ validator: validateHiringManager }]}
        style={formItemStyle}
      >
        <Input style={formInputStyle} />
      </Form.Item>
      <Form.Item
        label={t("currency")}
        name="currency"
        rules={[{ validator: validateCurrency }]}
        style={formItemStyle}
      >
        {/* Should fetch from API based on the job location */}
        <Select style={formSelectStyle}>
          <Select.Option value="USD">USD</Select.Option>
          <Select.Option value="INR">INR</Select.Option>
          <Select.Option value="EUR">EUR</Select.Option>
        </Select>
      </Form.Item>
      <Form.Item
        label={t("salary_min")}
        name="salarymin"
        rules={[{ validator: validateSalaryMin }]}
        style={formItemStyle}
      >
        <InputNumber style={formInputStyle} />
      </Form.Item>
      <Form.Item
        label={t("salary_max")}
        name="salarymax"
        rules={[{ validator: validateSalaryMax }]}
        style={formItemStyle}
      >
        <InputNumber style={formInputStyle} />
      </Form.Item>
      <Divider>{t("private_fields")}</Divider>
      <Form.Item
        label={t("department")}
        name="department"
        rules={[{ validator: validateDepartment }]}
        style={formItemStyle}
      >
        <Input style={formInputStyle} />
      </Form.Item>
      <Form.Item label={t("notes")} name="notes" style={formItemStyle}>
        <TextArea style={formInputStyle} />
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
