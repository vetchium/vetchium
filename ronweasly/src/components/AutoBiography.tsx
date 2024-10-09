import React, { useState } from "react";
import { Card, Typography, List, Button, Input, Space } from "antd";
import { EditOutlined, CheckOutlined, CloseOutlined } from "@ant-design/icons";
import t from "../i18n/i18n";

const { Text, Paragraph } = Typography;

interface UserInfo {
  displayName: string;
  profileUrl: string;
  aboutMe: string;
  websites: string[];
}

function AutoBiography() {
  const [isEditing, setIsEditing] = useState(false);
  const [userInfo, setUserInfo] = useState<UserInfo>({
    displayName: "John Doe",
    profileUrl: "johndoe",
    aboutMe: "I am a software developer.",
    websites: ["https://johndoe.com", "https://github.com/johndoe"],
  });
  const [tempInfo, setTempInfo] = useState<UserInfo>(userInfo);

  const handleEdit = () => {
    setIsEditing(true);
    setTempInfo(userInfo);
  };

  const handleSave = () => {
    setUserInfo(tempInfo);
    setIsEditing(false);
    // Add your callback function here
    console.log("Saved user info:", tempInfo);
  };

  const handleCancel = () => {
    setIsEditing(false);
  };

  const handleInputChange = (
    field: keyof UserInfo,
    value: string | string[]
  ) => {
    setTempInfo({ ...tempInfo, [field]: value });
  };

  const checkAvailability = () => {
    // Add your logic to check profile URL availability
    console.log("Checking availability for:", tempInfo.profileUrl);
  };

  const extra = isEditing ? (
    <Space>
      <Button icon={<CheckOutlined />} onClick={handleSave} />
      <Button icon={<CloseOutlined />} onClick={handleCancel} />
    </Space>
  ) : (
    <Button icon={<EditOutlined />} onClick={handleEdit} />
  );

  return (
    <Card title={t("myprofile.auto_biography")} extra={extra}>
      <Space direction="vertical" style={{ width: "100%" }}>
        <Text strong>{t("myprofile.display_name")}:</Text>
        {isEditing ? (
          <Input
            value={tempInfo.displayName}
            onChange={(e) => handleInputChange("displayName", e.target.value)}
          />
        ) : (
          <Paragraph>{userInfo.displayName}</Paragraph>
        )}

        <Text strong>{t("myprofile.profile_url")}:</Text>
        {isEditing ? (
          <Space>
            <Input
              value={tempInfo.profileUrl}
              onChange={(e) => handleInputChange("profileUrl", e.target.value)}
            />
            <Button onClick={checkAvailability}>
              {t("myprofile.check_availability")}
            </Button>
          </Space>
        ) : (
          <Paragraph>{userInfo.profileUrl}</Paragraph>
        )}

        <Text strong>{t("myprofile.about_me")}:</Text>
        {isEditing ? (
          <Input.TextArea
            value={tempInfo.aboutMe}
            rows={10}
            onChange={(e) => handleInputChange("aboutMe", e.target.value)}
          />
        ) : (
          <Paragraph>{userInfo.aboutMe}</Paragraph>
        )}

        <Text strong>{t("myprofile.websites")}:</Text>
        {isEditing ? (
          <Input.TextArea
            value={tempInfo.websites.join("\n")}
            onChange={(e) =>
              handleInputChange("websites", e.target.value.split("\n"))
            }
          />
        ) : (
          <List
            dataSource={userInfo.websites}
            renderItem={(item) => <List.Item>{item}</List.Item>}
          />
        )}
      </Space>
    </Card>
  );
}

export default AutoBiography;
