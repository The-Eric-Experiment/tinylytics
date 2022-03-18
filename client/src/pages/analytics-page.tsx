import React, { FunctionComponent, Suspense, useMemo } from "react";
import { useSearchParams } from "react-router-dom";
import { Filters, Periods } from "../api/types";
import { FilterBar } from "../components/filter-bar/filer-bar";
import { Browsers } from "../components/sections/browsers";
import { Countries } from "../components/sections/countries";
import { OSs } from "../components/sections/os";
import { Referrers } from "../components/sections/referrers";
import { Summary } from "../components/sections/summary";
import { SuspenseCard } from "../components/sections/suspense-card";
import {
  Card,
  GridItem,
  PageGrid,
  PageHeader,
  PageLayout,
} from "../components/shared/page-layout";
import { DEPENDANT_FILTERS } from "../constants/filters";

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
    const dependant = DEPENDANT_FILTERS[key];

    if (dependant && dependant.length) {
      for (const depKey of dependant) {
        delete n[depKey];
      }
    }

    delete n[key];
    const params = new URLSearchParams(n);
    setSearchParams(params, { replace: true });
  };

  return (
    <PageLayout>
      <PageHeader>
        <FilterBar
          filters={filters}
          removeFilter={removeFilter}
          onFilter={updateFilters}
        />
      </PageHeader>
      <PageGrid>
        <Suspense fallback={<div>Loading...</div>}>
          <Summary domain={domain} filters={filters} />
        </Suspense>
        <GridItem take={2}>
          <SuspenseCard>
            <Browsers
              domain={domain}
              filters={filters}
              onFilter={updateFilters}
            />
          </SuspenseCard>
        </GridItem>
        <GridItem take={2}>
          <SuspenseCard>
            <OSs domain={domain} filters={filters} onFilter={updateFilters} />
          </SuspenseCard>
        </GridItem>
        <GridItem take={2}>
          <SuspenseCard>
            <Countries
              domain={domain}
              filters={filters}
              onFilter={updateFilters}
            />
          </SuspenseCard>
        </GridItem>
        <GridItem take={2}>
          <SuspenseCard>
            <Referrers
              domain={domain}
              filters={filters}
              onFilter={updateFilters}
            />
          </SuspenseCard>
        </GridItem>
      </PageGrid>
    </PageLayout>
  );
};
