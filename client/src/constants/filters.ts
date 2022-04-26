import { SelectOption } from "react95";
import { Filters, Periods } from "../api/types";

export const FILTER_NAMES: Partial<Record<keyof Filters, string>> = {
  p: "Period",
  b: "Browser",
  bv: "Browser Version",
  os: "OS",
  osv: "OS Version",
  c: "Country",
  r: "Referrer",
  rfp: "Referrer",
  pg: "Page",
};

export const PERIODS: Array<SelectOption<Periods>> = [
  { value: Periods.P24H, label: "Last 24 Hours" },
  { value: Periods.TODAY, label: "Today" },
  { value: Periods.YESTERDAY, label: "Yesterday" },
  { value: Periods.P7D, label: "Last 7 Days" },
  { value: Periods.WEEK, label: "This Week" },
  { value: Periods.LASTWEEK, label: "Last Week" },
  { value: Periods.P30D, label: "Last 30 Days" },
  { value: Periods.P90D, label: "Last 90 Days" },
  { value: Periods.MONTH, label: "This Month" },
  { value: Periods.LASTMONTH, label: "Last Month" },
  { value: Periods.YEAR, label: "This Year" },
  { value: Periods.LASTYEAR, label: "Last Year" },
  { value: Periods.ALLTIME, label: "All Time" },
];

export const DEPENDANT_FILTERS: Partial<
  Record<keyof Filters, Array<keyof Filters>>
> = {
  b: ["bv"],
  os: ["osv"],
  r: ["rfp"],
};

export const SHOW_AS_SAME_FILTER: Array<Array<keyof Filters>> = [["r", "rfp"]];

export const SHOW_PREVIOUS_FILTER_IF_EMPTY: Array<keyof Filters> = ["rfp"];
