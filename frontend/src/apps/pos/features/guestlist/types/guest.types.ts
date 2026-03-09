export interface Guest {
  id: number;
  name: string;
  code?: string;
  listName?: string;
  additionalGuests: number;
  arrivalNote?: string;
  attendedGuests: number;
}
