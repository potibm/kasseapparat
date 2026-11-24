import CategoryIcon from "@mui/icons-material/Category";
import { ProductGroupList } from "./ProductGroupList";
import { ProductGroupEdit } from "./ProductGroupEdit";
import { ProductGroupCreate } from "./ProductGroupCreate";
import { ProductGroup } from "./types";

export default {
  list: ProductGroupList,
  edit: ProductGroupEdit,
  create: ProductGroupCreate,
  icon: CategoryIcon,
  recordRepresentation: (record: ProductGroup) => record.name,
};
