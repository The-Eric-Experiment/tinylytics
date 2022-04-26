import React, { FunctionComponent, Suspense } from "react";
import { useOSs } from "../../api/analytics";
import { SuspenseWindow } from "../shared/suspense-window";
import { TableWidget, TableWidgetWrapperProps } from "../shared/table-widget";

interface OSProps extends Omit<TableWidgetWrapperProps, "title"> {
  domain: string;
}

export const OSContent: FunctionComponent<OSProps> = ({ domain, ...props }) => {
  const { data } = useOSs(domain, props.filters);

  return (
    <TableWidget
      {...props}
      data={data}
      filterPrimary="os"
      filterSecondary="osv"
    />
  );
};

export const OS: FunctionComponent<OSProps> = (props) => {
  return (
    <SuspenseWindow title="Operating Systems">
      <OSContent {...props} />
    </SuspenseWindow>
  );
};
