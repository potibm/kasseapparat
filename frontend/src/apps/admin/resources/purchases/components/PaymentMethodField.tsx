import { RaRecord, useRecordContext, FieldProps } from "react-admin";
import { Box } from "@mui/material";
import { useConfig } from "@core/config/hooks/useConfig";
import { paymentMethodColor } from "@admin/ui/colors";

export const PaymentMethodField = <T extends RaRecord>(
  props: FieldProps<T>,
) => {
  const { source } = props;
  const record = useRecordContext<T>(props);
  const { paymentMethods } = useConfig();

  // If there is no record or the source field is missing, render nothing
  if (!record || !source || !record[source]) return null;

  const paymentMethoCodeFromRecord = String(record[source]);
  let paymentMethodName = paymentMethoCodeFromRecord;

  if (paymentMethods) {
    const paymentMethod = paymentMethods.find(
      (pm) => pm.code === paymentMethoCodeFromRecord,
    );
    if (paymentMethod) {
      paymentMethodName = paymentMethod.name;
    }
  }

  const index = (paymentMethods ?? []).findIndex(
    (pm) => pm.code === paymentMethoCodeFromRecord,
  );

  return (
    <Box component="span" sx={{ display: "inline-flex", alignItems: "center" }}>
      <Box
        component="span"
        aria-hidden="true"
        data-testid="payment-method-color"
        sx={{
          width: 10,
          height: 10,
          borderRadius: "50%",
          mr: 1,
          backgroundColor: paymentMethodColor(
            paymentMethoCodeFromRecord,
            Math.max(index, 0),
          ),
        }}
      />
      {paymentMethodName}
    </Box>
  );
};
