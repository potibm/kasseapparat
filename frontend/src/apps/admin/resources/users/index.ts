import PersonIcon from "@mui/icons-material/Person";
import { UserList } from "./UserList";
import { UserCreate } from "./UserCreate";
import { UserEdit } from "./UserEdit";

export default {
  list: UserList,
  edit: UserEdit,
  create: UserCreate,
  icon: PersonIcon,
  recordRepresentation: "username",
};
