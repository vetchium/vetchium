import React from "react";
import {
  FormControl,
  Autocomplete,
  TextField,
  FormHelperText,
} from "@mui/material";
import countries from "@psankar/vetchi-typespec/common/countries.json";
import { useTranslation } from "@/hooks/useTranslation";

interface Country {
  country_code: string;
  en: string;
}

interface CountrySelectProps {
  value: string;
  onChange: (value: string) => void;
  error?: boolean;
  helperText?: string;
}

export const CountrySelect: React.FC<CountrySelectProps> = ({
  value,
  onChange,
  error,
  helperText,
}) => {
  const { t } = useTranslation();
  const selectedCountry = value
    ? countries.find((country) => country.country_code === value)
    : null;

  return (
    <FormControl fullWidth margin="normal" error={error}>
      <Autocomplete
        value={selectedCountry}
        onChange={(_, newValue) => {
          onChange(newValue?.country_code || "");
        }}
        options={countries}
        getOptionLabel={(option: Country) =>
          `${option.en} (${option.country_code})`
        }
        renderInput={(params) => (
          <TextField
            {...params}
            label={t("hubUserOnboarding.form.countryCode")}
            required
            error={error}
            placeholder={t("hubUserOnboarding.form.countryCodePlaceholder")}
          />
        )}
        isOptionEqualToValue={(option, value) =>
          option.country_code === value.country_code
        }
        blurOnSelect
        clearOnBlur
      />
      {helperText && <FormHelperText>{helperText}</FormHelperText>}
    </FormControl>
  );
};
