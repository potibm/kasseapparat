import PersonIcon from "@mui/icons-material/Person";
import { GuestList } from "./GuestList";
import { GuestEdit } from "./GuestEdit";
import { GuestCreate } from "./GuestCreate";

export default {
  list: GuestList,
  edit: GuestEdit,
  create: GuestCreate,
  icon: PersonIcon,
  recordRepresentation: "name",
};
