import { DeleteTwoTone } from "@ant-design/icons";
import { Button, Divider, Flex, Table } from "antd";
import React from "react";
import { tableStyle } from "../Styles";
import LocationCreator from "../components/LocationCreator";
import t from "../i18n/i18n";

const Locations: React.FC = () => {
  // Sample data for the table
  const data = [
    { id: 1, country: "China", state: "Beijing", city: "Beijing" },
    { id: 2, country: "China", state: "Shanghai", city: "Shanghai" },
    { id: 3, country: "Germany", state: "Bavaria", city: "Nürnberg" },
    { id: 4, country: "India", state: "Karnataka", city: "Bangalore" },
    { id: 5, country: "India", state: "Tamil Nadu", city: "Chennai" },
    { id: 6, country: "Russia", state: "Moscow", city: "Moscow" },
    { id: 7, country: "USA", state: "California", city: "Palo Alto" },
    { id: 8, country: "USA", state: "Utah", city: "Provo" },
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
          <Button icon={<DeleteTwoTone />} />
        </span>
      ),
    },
  ];

  return (
    <Flex wrap vertical>
      <Divider>{t("locations.add_location")}</Divider>
      <LocationCreator />
      <Divider />
      <Table dataSource={data} columns={columns} style={tableStyle} />
    </Flex>
  );
};

export default Locations;
