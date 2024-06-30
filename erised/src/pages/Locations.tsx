import React from "react";
import { Table, Button, Flex, Divider } from "antd";
import LocationCreator from "../components/LocationCreator";
import { tableStyle } from "../Styles";
import t from "../i18n/i18n";

const Locations: React.FC = () => {
  // Sample data for the table
  const data = [
    { id: 1, country: "USA", state: "California", city: "Los Angeles" },
    { id: 2, country: "USA", state: "New York", city: "New York City" },
    { id: 3, country: "Canada", state: "Ontario", city: "Toronto" },
  ];

  // Columns configuration for the table
  const columns = [
    { title: "Country", dataIndex: "country", key: "country" },
    { title: "State", dataIndex: "state", key: "state" },
    { title: "City", dataIndex: "city", key: "city" },
    {
      title: "Actions",
      key: "actions",
      render: (text: string, record: any) => (
        <span>
          <Button type="link">Edit</Button>
          <Button type="link">Delete</Button>
        </span>
      ),
    },
  ];

  return (
    <Flex wrap vertical>
      <Table dataSource={data} columns={columns} style={tableStyle} />

      <Divider>{t("add_location")}</Divider>
      <LocationCreator />
    </Flex>
  );
};

export default Locations;
