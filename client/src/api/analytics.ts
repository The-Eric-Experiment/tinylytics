import useSWR, { SWRResponse } from "swr";

export interface Summaries {
  pageViews: number;
  sessions: number;
}

export interface AnalyticsData {
  name: string;
  count: number;
  major: string | null;
  minor: string | null;
  patch: string | null;
}

export interface AnalyticsDataResponse {
  items: AnalyticsData[];
}

export function get(endpoint: string) {
  return fetch(endpoint).then((o) => o.json());
}

export function useSummaries(domain: string): SWRResponse<Summaries> {
  return useSWR(`/api/${domain}/summaries?p=alltime`, get);
}

export function useBrowsers(
  domain: string
): SWRResponse<AnalyticsDataResponse> {
  return useSWR(`/api/${domain}/browsers?p=alltime`, get);
}

export function useOSs(domain: string): SWRResponse<AnalyticsDataResponse> {
  return useSWR(`/api/${domain}/os?p=alltime`, get);
}

export function useCountries(
  domain: string
): SWRResponse<AnalyticsDataResponse> {
  return useSWR(`/api/${domain}/countries?p=alltime`, get);
}
