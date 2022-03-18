export interface Summaries {
  pageViews: number;
  sessions: number;
}

export interface AnalyticsData {
  name: string;
  count: number;
  major: number | null;
  minor: number | null;
  patch: number | null;
}

export interface AnalyticsDataResponse {
  items: AnalyticsData[];
}

export enum Periods {
  TODAY = "today",
  YESTERDAY = "yesterday",
  P24H = "24h",
  WEEK = "week",
  LASTWEEK = "lastweek",
  P7D = "7d",
  MONTH = "month",
  LASTMONTH = "lastmonth",
  P30D = "30d",
  P90D = "90d",
  YEAR = "year",
  LASTYEAR = "lastyear",
  ALLTIME = "alltime",
}

export interface Filters {
  p: Periods;
  b?: string;
  bv?: string;
  c?: string;
  cv?: string; // get rid of this, country version
  os?: string;
  osv?: string;
  r?: string;
  rv?: string;
}

export interface FetcherRequest {
  url: string;
  filters: Filters;
}
