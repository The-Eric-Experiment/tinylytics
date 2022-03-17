import React, { FunctionComponent, Suspense, useMemo } from "react";
import { useSearchParams } from "react-router-dom";
import styled from "styled-components";
import { Filters, Periods } from "../api/types";
import { FilterBar } from "../components/filter-bar/filer-bar";
import { Browsers } from "../components/sections/browsers";
import { Countries } from "../components/sections/countries";
import { OSs } from "../components/sections/os";
import { Summary } from "../components/sections/summary";

export interface AnalyticsPageProps {}

export const AnalyticsPage: FunctionComponent<AnalyticsPageProps> = () => {
  const domain = "oldavista.com";
  const [searchParams, setSearchParams] = useSearchParams({
    p: Periods.P24H,
  });

  const filters: Filters = useMemo(() => {
    return {
      p: Periods.P24H,
      ...Object.fromEntries([...searchParams]),
    } as Filters;
  }, [searchParams]);

  const updateFilters = (newFilters: Partial<Filters>) => {
    const params = new URLSearchParams({ ...filters, ...newFilters });
    setSearchParams(params, { replace: true });
  };

  const removeFilter = (key: keyof Filters) => {
    const n = { ...filters };
    delete n[key];
    const params = new URLSearchParams(n);
    setSearchParams(params, { replace: true });
  };

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <FilterBar
        filters={filters}
        removeFilter={removeFilter}
        onFilter={updateFilters}
      />
      <Container>
        <Summary domain={domain} filters={filters} />
        <Browsers domain={domain} filters={filters} onFilter={updateFilters} />
        <OSs domain={domain} filters={filters} onFilter={updateFilters} />
        <Countries domain={domain} filters={filters} onFilter={updateFilters} />
      </Container>
    </Suspense>
  );
};

const Container = styled.div`
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  justify-content: space-between;
`;
