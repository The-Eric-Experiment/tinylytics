import { useMemo } from "react";
import { AnalyticsData, Filters } from "../api/types";

export function useAnalyticsFilters(
  filters: Filters,
  filterNames: { name: keyof Filters; version: keyof Filters },
  onFilter: (filters: Partial<Filters>) => void
) {
  const selectedFilter = useMemo(() => {
    if (filters[filterNames.version]) {
      return undefined;
    }

    return filters[filterNames.name];
  }, [filters]);

  function showAsLink(data: AnalyticsData) {
    if (filters[filterNames.version]) {
      let s = filters[filterNames.version]?.split(".");
      const parts = [data.major, data.minor, data.patch].filter(
        (o) => !!o || o === 0
      );

      return s!.length + 1 < parts.length;
    }

    return true;
  }

  function getLabel(data: AnalyticsData): string {
    if (filters[filterNames.version]) {
      let s = filters[filterNames.version]!.split(".");
      const parts = [data.major, data.minor, data.patch].filter(
        (o) => !!o || o === 0
      );
      if (s.length === parts.length) {
        return filters[filterNames.version]!;
      }

      return parts.slice(0, s.length + 1).join(".");
    }

    if (filters[filterNames.name] && (data.major || data.major === 0)) {
      return data.major.toString();
    }

    return data.name;
  }

  function updateFilter(data: AnalyticsData) {
    return () => {
      let nf = filters;

      if (!filters[filterNames.name]) {
        nf = { ...nf, [filterNames.name]: data.name };
      } else if (!filters[filterNames.version]) {
        nf = { ...nf, [filterNames.version]: getLabel(data) };
      }
      onFilter(nf);
    };
  }

  return { showAsLink, getLabel, updateFilter, selectedFilter };
}
