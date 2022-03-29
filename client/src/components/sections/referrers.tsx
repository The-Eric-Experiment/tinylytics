import React, { FunctionComponent } from "react";
import { useReferrers } from "../../api/analytics";
import { SuspenseCard } from "../shared/suspense-card";
import { TableWidget, TableWidgetWrapperProps } from "../shared/table-widget";

interface ReferrersProps extends TableWidgetWrapperProps {
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
      title="Referrers"
      filterPrimary="r"
      filterSecondary="rfp"
      showSelfWhenEmpty
    />
  );
};

export const Referrers: FunctionComponent<ReferrersProps> = (props) => {
  return (
    <SuspenseCard>
      <ReferrersContent {...props} />
    </SuspenseCard>
  );
};
