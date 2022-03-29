import React, { FunctionComponent } from "react";
import styled from "styled-components";
import { AnalyticsData, AnalyticsDataResponse, Filters } from "../../api/types";
import { TABLET } from "./page-layout";

interface TableWidgetProps {
  title: string;
  filterPrimary: keyof Filters;
  filterSecondary?: keyof Filters;
  data?: AnalyticsDataResponse;
  filters: Filters;
  showSelfWhenEmpty?: boolean;
  renderIcon?(item: AnalyticsData): React.ReactNode;
  formatName?(name: string): string;
  onFilter(filters: Partial<Filters>): void;
}

export interface TableWidgetWrapperProps
  extends Omit<
    TableWidgetProps,
    "filterPrimary" | "filterSecondary" | "data"
  > {}

export const TableWidget: FunctionComponent<TableWidgetProps> = ({
  title,
  filterPrimary,
  filterSecondary,
  data,
  showSelfWhenEmpty,
  filters,
  renderIcon,
  formatName,
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
    let label: string = getLabel(item);

    if (formatName) {
      label = formatName(label);
    }

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
    <Container>
      <Header>{title}</Header>
      {hasPreviousFilters && <div>{data.previousFilters.join(", ")}</div>}
      <Table>
        <thead>
          <tr>
            {renderIcon && <th className="icon"> </th>}
            <th>Name</th>
            <th className="count">Sessions</th>
          </tr>
        </thead>
        <tbody>
          {data!.items.map((item) => (
            <tr key={item.value + getLabel(item)}>
              {renderIcon && <td className="icon">{renderIcon(item)}</td>}
              <td>{renderLabel(item)}</td>
              <td className="value">{item.count}</td>
            </tr>
          ))}
        </tbody>
      </Table>
    </Container>
  );
};

const Container = styled.div``;

const Header = styled.h3`
  margin: 0;
`;

const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
  > thead,
  > tbody {
    > tr {
      > th,
      > td {
        border-bottom: 1px solid #ddd;
        text-align: left;
        padding: 8px 0;
      }
      > th.count,
      > td.value {
        text-align: right;
        font-weight: bold;
        border-left: 1px solid #ddd;
        width: 25%;
        @media all and (min-width: ${TABLET}) {
          width: 20%;
        }
      }
      > th.icon,
      > td.icon {
        width: 20px;
      }

      &:last-child {
        > td {
          border-bottom: 0;
        }
      }
    }
  }
`;
