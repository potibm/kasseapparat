import ViewAgendaIcon from "@mui/icons-material/ViewAgenda";
import { SectionList } from "./SectionList";
import { SectionEdit } from "./SectionEdit";
import { SectionCreate } from "./SectionCreate";
import { sectionLabel } from "./section.utils";
import { Section } from "./types";

export default {
  list: SectionList,
  edit: SectionEdit,
  create: SectionCreate,
  icon: ViewAgendaIcon,
  recordRepresentation: (record: Section) => sectionLabel(record),
};
