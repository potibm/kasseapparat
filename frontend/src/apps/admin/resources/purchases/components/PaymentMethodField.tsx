import { RaRecord, useRecordContext, FieldProps } from "react-admin";
import { useConfig } from "@core/config/providers/ConfigProvider";

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

  return <>{paymentMethodName}</>;
};
