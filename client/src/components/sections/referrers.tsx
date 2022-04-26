import React, { FunctionComponent } from "react";
import { useReferrers } from "../../api/analytics";
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

  return (
    <TableWidget
      {...props}
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
