import InventoryIcon from "@mui/icons-material/Inventory";
import { ProductList } from "./ProductList";
import { ProductEdit } from "./ProductEdit";
import { ProductCreate } from "./ProductCreate";

export default {
  list: ProductList,
  edit: ProductEdit,
  create: ProductCreate,
  icon: InventoryIcon,
  recordRepresentation: "name",
};
