import React, { FunctionComponent, Suspense, useMemo } from "react";
import { useSearchParams } from "react-router-dom";
import { Filters, Periods } from "../api/types";
import { Browsers } from "../components/sections/browsers";
import { Countries } from "../components/sections/countries";
import { OSs } from "../components/sections/os";
import { Summary } from "../components/sections/summary";

export interface AnalyticsPageProps {}

export const AnalyticsPage: FunctionComponent<AnalyticsPageProps> = () => {
  const domain = "oldavista.com";
  const [searchParams, setSearchParams] = useSearchParams({
    p: Periods.AllTime,
  });

  const filters: Filters = useMemo(() => {
    return {
      p: Periods.AllTime,
      ...Object.fromEntries([...searchParams]),
    } as Filters;
  }, [searchParams]);

  const updateFilters = (newFilters: Partial<Filters>) => {
    const params = new URLSearchParams({ ...filters, ...newFilters });
    setSearchParams(params, { replace: true });
  };

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Summary domain={domain} filters={filters} />
      <Browsers domain={domain} filters={filters} onFilter={updateFilters} />
      <OSs domain={domain} filters={filters} onFilter={updateFilters} />
      <Countries domain={domain} filters={filters} onFilter={updateFilters} />
    </Suspense>
  );
};
