import GroupsIcon from "@mui/icons-material/Groups";
import { GuestlistList } from "./GuestlistList";
import { GuestlistCreate } from "./GuestlistCreate";
import { GuestlistEdit } from "./GuestlistEdit";

export default {
  list: GuestlistList,
  edit: GuestlistEdit,
  create: GuestlistCreate,
  icon: GroupsIcon,
  recordRepresentation: "name",
};
