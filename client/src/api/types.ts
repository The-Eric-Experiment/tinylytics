export interface Summaries {
  pageViews: number;
  sessions: number;
  avgSessionDuration: number;
  bounceRate: number;
}

export interface AnalyticsData {
  value: string;
  count: number;
  drillable: number;
}

export interface AnalyticsDataResponse {
  previousFilters: string[];
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
  os?: string;
  osv?: string;
  c?: string;
  r?: string;
  rfp?: string;
  pg?: string;
}

export interface FetcherRequest {
  url: string;
  filters: Filters;
}

export interface Website {
  domain: string;
  title: string;
}
