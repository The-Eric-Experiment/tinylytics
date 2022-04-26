import React, { FunctionComponent } from "react";
import { useBrowsers } from "../../api/analytics";
import { SuspenseWindow } from "../shared/suspense-window";
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
      data={data}
      filterPrimary="b"
      filterSecondary="bv"
    />
  );
};

export const Browsers: FunctionComponent<BrowsersProps> = (props) => {
  return (
    <SuspenseWindow title="Browsers">
      <BrowsersContent {...props} />
    </SuspenseWindow>
  );
};
