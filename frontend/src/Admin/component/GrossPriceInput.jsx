import React, { useEffect, useMemo } from "react";
import { useFormContext } from "react-hook-form";
import { TextInput } from "react-admin";
import Decimal from "decimal.js";
import PropTypes from "prop-types";

const GrossPriceInput = ({
  netSource,
  vatSource,
  options = {},
  label = "Gross price",
}) => {
  const { watch, setValue } = useFormContext();

  const netPrice = watch(netSource) || "0";
  const vatRate = watch(vatSource) || "0";
  const grossSource = `${netSource}_gross`;
  const maximumFractionDigits = options.maximumFractionDigits || 2;

  const grossPrice = useMemo(() => {
    try {
      return new Decimal(netPrice)
        .mul(new Decimal(1).plus(new Decimal(vatRate).div(100)))
        .toFixed(maximumFractionDigits);
    } catch (error) {
      console.error("Error while determing the grossPrice: " + error);
      return new Decimal(0).toFixed(maximumFractionDigits);
    }
  }, [netPrice, vatRate, maximumFractionDigits]);

  useEffect(() => {
    setValue(grossSource, grossPrice);
  }, [grossPrice, grossSource, setValue]);

  return (
    <TextInput label={label} source={grossSource} type="number" disabled />
  );
};

GrossPriceInput.propTypes = {
  netSource: PropTypes.string,
  vatSource: PropTypes.string,
  label: PropTypes.string,
  options: PropTypes.object,
};

export default GrossPriceInput;
