import React, { FunctionComponent } from "react";
import { AnalyticsData, AnalyticsDataResponse, Filters } from "../../api/types";

interface TableWidgetProps {
  filterPrimary: keyof Filters;
  filterSecondary?: keyof Filters;
  data?: AnalyticsDataResponse;
  filters: Filters;
  showSelfWhenEmpty?: boolean;
  onFilter(filters: Partial<Filters>): void;
}

export interface TableWidgetWrapperProps
  extends Omit<
    TableWidgetProps,
    "filterPrimary" | "filterSecondary" | "data"
  > {}

export const TableWidget: FunctionComponent<TableWidgetProps> = ({
  filterPrimary,
  filterSecondary,
  data,
  showSelfWhenEmpty,
  filters,
  onFilter,
}) => {
  const hasPreviousFilters =
    data?.previousFilters && data.previousFilters.length > 0;

  function showAsLink(item: AnalyticsData) {
    return data && (item.drillable > 0 || data.items.length > 1);
  }

  function getLabel(item: AnalyticsData): string {
    if (filterSecondary && data!.previousFilters.length > 1) {
      return [...data!.previousFilters.slice(1), item.value]
        .filter((o) => !!o)
        .join(".");
    }

    if (!item.value && showSelfWhenEmpty) {
      return data!.previousFilters[data!.previousFilters.length - 1];
    }

    if (!item.value && !showSelfWhenEmpty) {
      return "(unknown)";
    }

    return item.value;
  }

  function getFilter(item: AnalyticsData): string {
    if (
      filterSecondary &&
      filters[filterSecondary] &&
      data!.previousFilters.length >= 1
    ) {
      return [...data!.previousFilters.slice(1), item.value]
        .filter((o) => !!o)
        .join("/");
    }

    return item.value;
  }

  function updateFilter(item: AnalyticsData) {
    return () => {
      let nf = filters;

      if (!filters[filterPrimary]) {
        nf = { ...nf, [filterPrimary]: item.value };
      } else if (filterSecondary) {
        nf = { ...nf, [filterSecondary]: getFilter(item) };
      }
      onFilter(nf);
    };
  }

  const renderLabel = (item: AnalyticsData) => {
    const label = getLabel(item);
    if (!showAsLink(item)) {
      return label;
    }

    return (
      <a onClick={updateFilter(item)} type="button" href="">
        {label}
      </a>
    );
  };

  return (
    <div>
      {hasPreviousFilters && <div>{data.previousFilters.join(", ")}</div>}
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Count</th>
          </tr>
        </thead>
        <tbody>
          {data!.items.map((item) => (
            <tr key={item.value + getLabel(item)}>
              <td>{renderLabel(item)}</td>
              <td>{item.count}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
