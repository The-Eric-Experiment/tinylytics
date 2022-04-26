import React, { FunctionComponent } from "react";
import { usePages } from "../../api/analytics";
import { SuspenseWindow } from "../shared/suspense-window";
import { TableWidget, TableWidgetWrapperProps } from "../shared/table-widget";

interface PagesProps extends Omit<TableWidgetWrapperProps, "title"> {
  domain: string;
}

const PagesContent: FunctionComponent<PagesProps> = ({ domain, ...props }) => {
  const { data } = usePages(domain, props.filters);

  const formatName = (input: string) => {
    try {
      let url = input;
      if (!url.startsWith("http://") && !url.startsWith("https://")) {
        url = "http://" + input;
      }
      const parts = new URL(url);
      return parts.pathname + (parts.search || "");
    } catch (_) {
      return input;
    }
  };

  return (
    <TableWidget
      {...props}
      data={data}
      filterPrimary="pg"
      showSelfWhenEmpty
      formatName={formatName}
    />
  );
};

export const Pages: FunctionComponent<PagesProps> = (props) => {
  return (
    <SuspenseWindow title="Pages">
      <PagesContent {...props} />
    </SuspenseWindow>
  );
};
