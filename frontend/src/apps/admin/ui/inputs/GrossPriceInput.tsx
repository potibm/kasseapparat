import React, { useEffect } from "react";
import { useFormContext, useWatch } from "react-hook-form";
import { NumberInput, NumberInputProps } from "react-admin";
import Decimal from "decimal.js";
import PropTypes from "prop-types";

interface GrossPriceInputProps extends NumberInputProps {
  netSource: string;
  vatSource: string;
  options?: Intl.NumberFormatOptions;
}

const GrossPriceInput: React.FC<GrossPriceInputProps> = ({
  netSource,
  vatSource,
  options,
  ...props
}) => {
  const { setValue } = useFormContext();

  const netPrice = useWatch({ name: netSource });
  const vatRate = useWatch({ name: vatSource });
  const precision = options?.maximumFractionDigits ?? 2;

  useEffect(() => {
    if (netPrice !== undefined && vatRate !== undefined) {
      try {
        const net = new Decimal(netPrice || 0);
        const vat = new Decimal(vatRate || 0).div(100).plus(1);

        // Calculate gross price: net * (1 + vatRate/100)
        const gross = net.mul(vat).toDecimalPlaces(precision);

        // Update the grossPrice field in the form
        setValue(props.source, gross.toNumber(), { shouldDirty: true });
      } catch (error) {
        // Silently fail if decimal conversion fails during typing
        console.error("Calculation error:", error);
      }
    }
  }, [netPrice, vatRate, setValue, precision, props.source]);

  return (
    <NumberInput
      {...props}
      disabled
      helperText="Calculated automatically based on net price and VAT"
    />
  );
};

GrossPriceInput.propTypes = {
  netSource: PropTypes.string,
  vatSource: PropTypes.string,
  label: PropTypes.string,
  options: PropTypes.object,
};

export default GrossPriceInput;
