import React, { FunctionComponent } from "react";
import styled from "styled-components";
import AnalyticsTable, { AnalyticsTableProps } from "./analytics-table";

interface TableWidgetProps extends AnalyticsTableProps {}

export interface TableWidgetWrapperProps
  extends Omit<
    TableWidgetProps,
    "filterPrimary" | "filterSecondary" | "data"
  > {}

export const TableWidget: FunctionComponent<TableWidgetProps> = ({
  ...tableProps
}) => {
  const { data } = tableProps;
  const hasPreviousFilters =
    data?.previousFilters && data.previousFilters.length > 0;

  return (
    <TableWidgetContainer>
      {hasPreviousFilters && <div>{data.previousFilters.join(", ")}</div>}
      <AnalyticsTable {...tableProps} />
    </TableWidgetContainer>
  );
};

export const TableWidgetContainer = styled.div`
  height: 100%;
  display: flex;
  flex-direction: column;
  & > div:last-child {
    box-shadow: inset 2px 2px 3px rgba(0, 0, 0, 0.2);
    background-color: #fff;
    flex: 1;

    tbody {
      background-color: transparent;
      box-shadow: unset;
    }
  }
`;
