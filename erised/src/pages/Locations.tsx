import React, { useState } from "react";
import { Select } from "antd";
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

const LocationSelector: React.FC = () => {
  const [selectedCountry, setSelectedCountry] = useState<string>("");
  const [selectedState, setSelectedState] = useState<string>("");
  const [states, setStates] = useState<State[]>([]);
  const [cities, setCities] = useState<City[]>([]);

  const handleCountryChange = (value: string) => {
    const country = countriesStatesCities.find(
      (country: Country) => country.name === value
    );
    setSelectedCountry(value);
    if (country) {
      setStates(country.states);
      setCities([]);
      setSelectedState("");
    }
  };

  const handleStateChange = (value: string) => {
    const state = states.find((state: State) => state.name === value);
    setSelectedState(value);
    if (state) {
      setCities(state.cities);
    }
  };

  return (
    <div>
      <Select
        style={{ width: 200, marginRight: 8 }}
        placeholder="Select Country"
        onChange={handleCountryChange}
        value={selectedCountry}
      >
        {countriesStatesCities.map((country: Country) => (
          <Option key={country.name} value={country.name}>
            {country.name}
          </Option>
        ))}
      </Select>
      <Select
        style={{ width: 200, marginRight: 8 }}
        placeholder="Select State"
        onChange={handleStateChange}
        value={selectedState}
        disabled={!states.length}
      >
        {states.map((state: State) => (
          <Option key={state.name} value={state.name}>
            {state.name}
          </Option>
        ))}
      </Select>
      <Select
        style={{ width: 200 }}
        placeholder="Select City"
        disabled={!cities.length}
      >
        {cities.map((city: City) => (
          <Option key={city.name} value={city.name}>
            {city.name}
          </Option>
        ))}
      </Select>
    </div>
  );
};

export default LocationSelector;
