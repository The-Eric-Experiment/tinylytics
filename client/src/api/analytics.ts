import qs from "qs";
import useSWR, { SWRResponse } from "swr";
import {
  AnalyticsDataResponse,
  FetcherRequest,
  Filters,
  Summaries,
} from "./types";

export function fetcher(req: FetcherRequest) {
  const endpoint =
    req.url + qs.stringify(req.filters, { addQueryPrefix: true });

  return fetch(endpoint).then((o) => o.json());
}

export function useGet<T>(
  domain: string,
  endpoint: string,
  filters: any
): SWRResponse<T> {
  return useSWR({ url: `/api/${domain}/${endpoint}`, filters }, fetcher, {
    suspense: true,
  });
}

export function useSummaries(
  domain: string,
  filters: Filters
): SWRResponse<Summaries> {
  return useGet(domain, "summaries", filters);
}

export function useBrowsers(
  domain: string,
  filters: Filters
): SWRResponse<AnalyticsDataResponse> {
  return useGet(domain, "browsers", filters);
}

export function useOSs(
  domain: string,
  filters: Filters
): SWRResponse<AnalyticsDataResponse> {
  return useGet(domain, "os", filters);
}

export function useCountries(
  domain: string,
  filters: Filters
): SWRResponse<AnalyticsDataResponse> {
  return useGet(domain, "countries", filters);
}

export function useReferrers(
  domain: string,
  filters: Filters
): SWRResponse<AnalyticsDataResponse> {
  return useGet(domain, "referrers", filters);
}
