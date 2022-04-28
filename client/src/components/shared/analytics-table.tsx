import React, { FunctionComponent } from "react";
import {
  Table,
  TableBody,
  TableDataCell,
  TableHead,
  TableHeadCell,
  TableRow,
} from "react95";
import { AnalyticsData, AnalyticsDataResponse, Filters } from "../../api/types";

export interface AnalyticsTableProps {
  filterPrimary: keyof Filters;
  filterSecondary?: keyof Filters;
  data?: AnalyticsDataResponse;
  filters: Filters;
  showSelfWhenEmpty?: boolean;
  renderIcon?(item: AnalyticsData): React.ReactNode;
  formatName?(name: string): string;
  onFilter(filters: Partial<Filters>): void;
}

const AnalyticsTable: FunctionComponent<AnalyticsTableProps> = ({
  filterPrimary,
  filterSecondary,
  data,
  showSelfWhenEmpty,
  filters,
  renderIcon,
  formatName,
  onFilter,
}) => {
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

    return label;
  };

  return (
    <Table>
      <TableHead>
        <TableRow head>
          {renderIcon && <TableHeadCell />}
          <TableHeadCell>Name</TableHeadCell>
          <TableHeadCell disabled>Sessions</TableHeadCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {data!.items.map((item) => (
          <TableRow
            key={item.value + getLabel(item)}
            onClick={
              showAsLink(item)
                ? updateFilter(item)
                : () => {
                    /** */
                  }
            }
          >
            {renderIcon && (
              <TableDataCell style={{ textAlign: "center", width: "20px" }}>
                <span role="img" aria-label={getLabel(item)}>
                  {renderIcon(item)}
                </span>
              </TableDataCell>
            )}
            <TableDataCell>{renderLabel(item)}</TableDataCell>
            <TableDataCell style={{ textAlign: "right", width: "50px" }}>
              {item.count}
            </TableDataCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};

export default AnalyticsTable;
