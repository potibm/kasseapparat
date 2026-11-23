import { use } from "react";
import { SumupReaderContext } from "../context/SumupReaderContext";

/** The card reader of this device, shared by the checkout buttons and the menu. */
export const useSumupReader = () => use(SumupReaderContext);
