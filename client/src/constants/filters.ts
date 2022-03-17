import { Filters, Periods } from "../api/types";

export const FILTER_NAMES: Partial<Record<keyof Filters, string>> = {
  p: "Period",
  b: "Browser",
  bv: "Browser Version",
  os: "OS",
  osv: "OS Version",
  c: "Country",
};

export const PERIOD_NAMES: Record<Periods, string> = {
  [Periods.ALLTIME]: "All Time",
  [Periods.LASTMONTH]: "Last Month",
  [Periods.LASTWEEK]: "Last Week",
  [Periods.LASTYEAR]: "Last Year",
  [Periods.MONTH]: "This Month",
  [Periods.P24H]: "Last 24 Hours",
  [Periods.P30D]: "Last 30 Days",
  [Periods.P7D]: "Last 7 Days",
  [Periods.P90D]: "Last 90 Days",
  [Periods.TODAY]: "Today",
  [Periods.WEEK]: "This Week",
  [Periods.YEAR]: "This Year",
  [Periods.YESTERDAY]: "Yesterday",
};
