import React, { FunctionComponent, useEffect, useMemo } from "react";
import { useNavigate, useParams } from "react-router";
import { useSearchParams } from "react-router-dom";
import { useWebsites } from "../api/analytics";
import { Filters, Periods } from "../api/types";
import { FilterBar } from "../components/filter-bar/filer-bar";
import { Browsers } from "../components/sections/browsers";
import { Countries } from "../components/sections/countries";
import { OS } from "../components/sections/os";
import { Pages } from "../components/sections/pages";
import { Referrers } from "../components/sections/referrers";
import { Summary } from "../components/sections/summary";
import {
  GridItem,
  PageGrid,
  PageHeader,
  PageLayout,
} from "../components/shared/page-layout";
import { DEPENDANT_FILTERS } from "../constants/filters";

export interface AnalyticsPageProps {}

export const AnalyticsPage: FunctionComponent<AnalyticsPageProps> = () => {
  const params = useParams();
  const navigate = useNavigate();
  const { data } = useWebsites();

  useEffect(() => {
    if (!params.domain) {
      const site = data ? data[0] : undefined;
      navigate(`/${site ? site.domain : ""}`);
    }
  }, [params.domain]);
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

  if (!params.domain) {
    return null;
  }

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
        <Summary domain={params.domain} filters={filters} />
        <GridItem take={2}>
          <Browsers
            domain={params.domain}
            filters={filters}
            onFilter={updateFilters}
          />
        </GridItem>
        <GridItem take={2}>
          <OS
            domain={params.domain}
            filters={filters}
            onFilter={updateFilters}
          />
        </GridItem>
        <GridItem take={2}>
          <Pages
            domain={params.domain}
            filters={filters}
            onFilter={updateFilters}
          />
        </GridItem>
        <GridItem take={2}>
          <Referrers
            domain={params.domain}
            filters={filters}
            onFilter={updateFilters}
          />
        </GridItem>
        <GridItem take={4}>
          <Countries
            domain={params.domain}
            filters={filters}
            onFilter={updateFilters}
          />
        </GridItem>
      </PageGrid>
    </PageLayout>
  );
};
