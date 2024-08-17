import { Button, Divider, Flex, Input, List, Select, Typography } from "antd";
import { useState } from "react";
import t from "../i18n/i18n";
import { formInputStyle, formSelectStyle } from "../Styles";

const { Text, Link } = Typography;

const data = [
  {
    name: "John Doe",
    lastPosition: "Software Engineer",
    lastCompany: "Company A",
    vouchedBy: [
      { name: "Jane Smith", position: "Manager", company: "Company A" },
      { name: "Bob Johnson", position: "Team Lead", company: "Company X" },
    ],
    referredBy: "Sankar P <sankar@example.com>",
  },
  {
    name: "Alice Brown",
    lastPosition: "Product Manager",
    lastCompany: "Company B",
    vouchedBy: [
      { name: "Charlie Davis", position: "Director", company: "Company B" },
      { name: "Eve White", position: "CEO", company: "Company Y" },
    ],
  },
  {
    name: "Michael Green",
    lastPosition: "Data Scientist",
    lastCompany: "Company C",
    vouchedBy: [
      { name: "Fiona Black", position: "CTO", company: "Company C" },
      {
        name: "George Blue",
        position: "Lead Data Scientist",
        company: "Company Z",
      },
    ],
  },
  {
    name: "Sankarasivasubramanian Pasupathilingam",
    lastPosition: "UX Designer",
    lastCompany: "Company D",
    vouchedBy: [
      { name: "Hannah Red", position: "Head of Design", company: "Company D" },
      { name: "Ian Yellow", position: "Senior Designer", company: "Company W" },
    ],
  },
];

export default function Applications() {
  const [filterText, setFilterText] = useState("");
  const [filterColor, setFilterColor] = useState([]);

  const filteredData = data.filter((item) => {
    const matchesText =
      item.name.toLowerCase().includes(filterText.toLowerCase()) ||
      (item.lastPosition &&
        item.lastPosition.toLowerCase().includes(filterText.toLowerCase()));
    const matchesColor =
      filterColor.length === 0 ||
      filterColor.some((color) =>
        item.vouchedBy.some((vouched) =>
          vouched.position.toLowerCase().includes(color)
        )
      );
    return matchesText && matchesColor;
  });

  return (
    <Flex vertical>
      <Flex
        style={{ width: "100%", margin: "16px 0", padding: "10px" }}
        gap="large"
      >
        <Input
          type="text"
          placeholder={t("applications.filter_by_name_or_employer")}
          value={filterText}
          onChange={(e) => setFilterText(e.target.value)}
          style={formInputStyle}
        />
        <Select
          mode="multiple"
          placeholder="Filter by color"
          value={filterColor}
          onChange={setFilterColor}
          style={formSelectStyle}
          options={[
            {
              value: "red",
              label: (
                <div
                  style={{
                    width: "20px",
                    height: "20px",
                    backgroundColor: "red",
                    borderRadius: "50%",
                    display: "inline-block",
                  }}
                ></div>
              ),
            },
            {
              value: "blue",
              label: (
                <div
                  style={{
                    width: "20px",
                    height: "20px",
                    backgroundColor: "blue",
                    borderRadius: "50%",
                    display: "inline-block",
                  }}
                ></div>
              ),
            },
            {
              value: "green",
              label: (
                <div
                  style={{
                    width: "20px",
                    height: "20px",
                    backgroundColor: "green",
                    borderRadius: "50%",
                    display: "inline-block",
                  }}
                ></div>
              ),
            },
            {
              value: "purple",
              label: (
                <div
                  style={{
                    width: "20px",
                    height: "20px",
                    backgroundColor: "purple",
                    borderRadius: "50%",
                    display: "inline-block",
                  }}
                ></div>
              ),
            },
          ]}
        />
      </Flex>
      <Divider />
      <List
        itemLayout="vertical"
        dataSource={filteredData}
        renderItem={(item) => (
          <List.Item style={{ marginBottom: "36px" }}>
            <Flex vertical gap="large">
              <Flex justify="space-between">
                <Flex vertical>
                  <Typography.Title level={3} style={{ margin: 0 }}>
                    {item.name}
                  </Typography.Title>
                  <Text>
                    {item.lastPosition}, {item.lastCompany}
                  </Text>
                </Flex>

                <Flex vertical gap="small">
                  <Text strong>{t("applications.vouched_by")}</Text>
                  {item.vouchedBy.map((vouched, index) => (
                    <Link href="#">
                      {vouched.name} ({vouched.position}, {vouched.company})
                    </Link>
                  ))}
                  {item.referredBy && (
                    <>
                      <Text strong>{t("applications.referred_by")}</Text>
                      <Text>{item.referredBy}</Text>
                    </>
                  )}
                </Flex>

                <Link href="#">{t("applications.resume")}</Link>
                <Flex vertical gap="small">
                  <Select
                    mode="multiple"
                    placeholder={t("applications.add_tags")}
                    options={[
                      {
                        value: "red",
                        label: (
                          <div
                            style={{
                              width: "20px",
                              height: "20px",
                              backgroundColor: "red",
                              borderRadius: "50%",
                              display: "inline-block",
                            }}
                          ></div>
                        ),
                      },
                      {
                        value: "blue",
                        label: (
                          <div
                            style={{
                              width: "20px",
                              height: "20px",
                              backgroundColor: "blue",
                              borderRadius: "50%",
                              display: "inline-block",
                            }}
                          ></div>
                        ),
                      },
                      {
                        value: "green",
                        label: (
                          <div
                            style={{
                              width: "20px",
                              height: "20px",
                              backgroundColor: "green",
                              borderRadius: "50%",
                              display: "inline-block",
                            }}
                          ></div>
                        ),
                      },
                      {
                        value: "purple",
                        label: (
                          <div
                            style={{
                              width: "20px",
                              height: "20px",
                              backgroundColor: "purple",
                              borderRadius: "50%",
                              display: "inline-block",
                            }}
                          ></div>
                        ),
                      },
                    ]}
                    getPopupContainer={(triggerNode) =>
                      triggerNode.parentElement
                    }
                  />
                  <Button type="primary">{t("applications.shortlist")}</Button>
                </Flex>
              </Flex>
              <Flex>
                <Button danger>{t("applications.reject")}</Button>
              </Flex>
            </Flex>
          </List.Item>
        )}
      />
    </Flex>
  );
}
