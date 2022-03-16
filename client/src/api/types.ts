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
  Today = "today",
  Yesterday = "yesterday",
  H24 = "24h",
  AllTime = "alltime",
}

export interface Filters {
  p: Periods;
  b?: string;
  bv?: string;
  c?: string;
  os?: string;
  osv?: string;
}

export interface FetcherRequest {
  url: string;
  filters: Filters;
}
