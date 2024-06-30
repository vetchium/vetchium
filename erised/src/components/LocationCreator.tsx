import React, { useState, useEffect } from "react";
import { Button, Flex, Select } from "antd";
import { PlusCircleTwoTone } from "@ant-design/icons"; // Import PlusCircleTwoTone from @ant-design/icons
import countriesStatesCitiesData from "../static/countries-states-cities.json";

const { Option } = Select;

interface City {
  name: string;
}

interface State {
  name: string;
  cities: City[];
}

interface Country {
  name: string;
  states: State[];
}

const countriesStatesCities: Country[] = countriesStatesCitiesData as Country[];

const LocationCreator: React.FC = () => {
  const [selectedCountry, setSelectedCountry] = useState<string>("");
  const [selectedState, setSelectedState] = useState<string>("");
  const [states, setStates] = useState<State[]>([]);
  const [cities, setCities] = useState<City[]>([]);

  useEffect(() => {
    const country = countriesStatesCities.find(
      (c) => c.name === selectedCountry
    );
    if (country) {
      setStates(country.states);
    } else {
      setStates([]);
    }
    setSelectedState("");
    setCities([]);
  }, [selectedCountry]);

  useEffect(() => {
    const state = states.find((s) => s.name === selectedState);
    if (state) {
      setCities(state.cities);
    } else {
      setCities([]);
    }
  }, [selectedState, states]);

  return (
    <Flex wrap justify="center">
      <Select
        showSearch
        style={{ width: 200, marginRight: 8 }}
        placeholder="Select Country"
        onChange={(value) => setSelectedCountry(value)}
        value={selectedCountry}
        filterOption={(input, option) =>
          option?.value
            ? option.value
                .toString()
                .toLowerCase()
                .indexOf(input.toLowerCase()) >= 0
            : false
        }
      >
        {countriesStatesCities.map((country) => (
          <Option key={country.name} value={country.name}>
            {country.name}
          </Option>
        ))}
      </Select>
      <Select
        showSearch
        style={{ width: 200, marginRight: 8 }}
        placeholder="Select State"
        onChange={(value) => setSelectedState(value)}
        value={selectedState}
        disabled={!selectedCountry}
      >
        {states.map((state) => (
          <Option key={state.name} value={state.name}>
            {state.name}
          </Option>
        ))}
      </Select>
      <Select
        showSearch
        style={{ width: 200 }}
        placeholder="Select City"
        disabled={!selectedState}
      >
        {cities.map((city) => (
          <Option key={city.name} value={city.name}>
            {city.name}
          </Option>
        ))}
      </Select>
      <Button
        type="primary"
        icon={<PlusCircleTwoTone />}
        style={{ marginLeft: 8 }}
      >
        Add Location
      </Button>
    </Flex>
  );
};

export default LocationCreator;
