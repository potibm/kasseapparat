import { useMediaQuery, Theme } from "@mui/material";

/** True on phones: the admin lists switch from tables to simple lists there. */
export const useIsSmallScreen = () =>
  useMediaQuery((theme: Theme) => theme.breakpoints.down("sm"));
