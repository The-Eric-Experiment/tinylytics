import React, { FunctionComponent } from "react";
import { useBrowsers } from "../../api/analytics";
import { SuspenseCard } from "../shared/suspense-card";
import { TableWidget, TableWidgetWrapperProps } from "../shared/table-widget";

interface BrowsersProps extends TableWidgetWrapperProps {
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
