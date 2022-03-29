import React, { FunctionComponent } from "react";
import { useOSs } from "../../api/analytics";
import { SuspenseCard } from "../shared/suspense-card";
import { TableWidget, TableWidgetWrapperProps } from "../shared/table-widget";

interface OSProps extends TableWidgetWrapperProps {
  domain: string;
}

export const OSContent: FunctionComponent<OSProps> = ({ domain, ...props }) => {
  const { data } = useOSs(domain, props.filters);

  return (
    <TableWidget
      {...props}
      title="Operating System"
      data={data}
      filterPrimary="os"
      filterSecondary="osv"
    />
  );
};

export const OS: FunctionComponent<OSProps> = (props) => {
  return (
    <SuspenseCard>
      <OSContent {...props} />
    </SuspenseCard>
  );
};
