import React, { FunctionComponent } from "react";
import { useBrowsers } from "../../api/analytics";
import { SuspenseCard } from "../shared/suspense-card";
import { TableWidget, TableWidgetWrapperProps } from "../shared/table-widget";

interface BrowsersProps extends Omit<TableWidgetWrapperProps, "title"> {
  domain: string;
}

const BrowsersContent: FunctionComponent<BrowsersProps> = ({
  domain,
  ...props
}) => {
  const { data } = useBrowsers(domain, props.filters);

  return (
    <TableWidget
      {...props}
      title="Browsers"
      data={data}
      filterPrimary="b"
      filterSecondary="bv"
    />
  );
};

export const Browsers: FunctionComponent<BrowsersProps> = (props) => {
  return (
    <SuspenseCard>
      <BrowsersContent {...props} />
    </SuspenseCard>
  );
};
