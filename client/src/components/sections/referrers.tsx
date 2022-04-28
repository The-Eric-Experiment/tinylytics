import React, { FunctionComponent } from "react";
import { useReferrers } from "../../api/analytics";
import { AnalyticsData } from "../../api/types";
import { SuspenseWindow } from "../shared/suspense-window";
import { TableWidget, TableWidgetWrapperProps } from "../shared/table-widget";

interface ReferrersProps extends Omit<TableWidgetWrapperProps, "title"> {
  domain: string;
}

const ReferrersContent: FunctionComponent<ReferrersProps> = ({
  domain,
  ...props
}) => {
  const { data } = useReferrers(domain, props.filters);

  const renderIcon = (item: AnalyticsData) => {
    return (
      <img
        src={`https://www.google.com/s2/favicons?domain=${item.value}&sz=16`}
        alt={item.value}
      />
    );
  };

  return (
    <TableWidget
      {...props}
      renderIcon={renderIcon}
      data={data}
      filterPrimary="r"
      filterSecondary="rfp"
      showSelfWhenEmpty
    />
  );
};

export const Referrers: FunctionComponent<ReferrersProps> = (props) => {
  return (
    <SuspenseWindow title="Referrers">
      <ReferrersContent {...props} />
    </SuspenseWindow>
  );
};
