import React, { FunctionComponent } from "react";
import styled from "styled-components";
import { useBrowsers } from "../../api/analytics";
import { AnalyticsData } from "../../api/types";
import * as browsers from "../../assets/browsers";
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

  const getIcon = (item: AnalyticsData) => {
    const icon = ICONS.find((i) => i.regex(item.value));
    if (!icon) {
      return browsers.unknown;
    }

    return icon.icon;
  };

  const renderIcon = (item: AnalyticsData) => {
    const icon = getIcon(item);
    return <BrowserIcon src={icon} alt={item.value} />;
  };

  return (
    <TableWidget
      {...props}
      data={data}
      renderIcon={renderIcon}
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

const BrowserIcon = styled.img`
  width: 16px;
  height: 16px;
`;

const ICONS = [
  {
    regex: (input: string) => !!/(^|\s)(brave)(\s|$)/gim.exec(input),
    icon: browsers.brave,
  },
  {
    regex: (input: string) => !!/(^|\s)chrome(\s|$)/gim.exec(input),
    icon: browsers.chrome,
  },
  {
    regex: (input: string) => !!/(^|\s)chromium(\s|$)/gim.exec(input),
    icon: browsers.chromium,
  },
  {
    regex: (input: string) => !!/(^|\s)dreamkey(\s|$)/gim.exec(input),
    icon: browsers.dreamkey,
  },
  {
    regex: (input: string) => !!/(^|\s)dream[\s]*passport(\s|$)/gim.exec(input),
    icon: browsers.dreampassport,
  },
  {
    regex: (input: string) => !!/(^|\s)edge(\s|$)/gim.exec(input),
    icon: browsers.edge,
  },
  {
    regex: (input: string) => !!/(^|\s)firefox(\s|$)/gim.exec(input),
    icon: browsers.firefox,
  },
  {
    regex: (input: string) => !!/(^|\s)flock/gim.exec(input),
    icon: browsers.flock,
  },
  {
    regex: (input: string) => !!/(^|\s)konqueror/gim.exec(input),
    icon: browsers.konqueror,
  },
  {
    regex: (input: string) =>
      !!/(^|\s)(ie)|(internet[\s]*explorer)(\s|$)/gim.exec(input),
    icon: browsers.ie,
  },
  {
    regex: (input: string) => !!/(^|\s)maxthon(\s|$)/gim.exec(input),
    icon: browsers.maxthon,
  },
  {
    regex: (input: string) => !!/(^|\s)mosaic(\s|$)/gim.exec(input),
    icon: browsers.mosaic,
  },
  {
    regex: (input: string) => !!/(^|\s)mozilla(\s|$)/gim.exec(input),
    icon: browsers.mozilla,
  },
  {
    regex: (input: string) => !!/(^|\s)netpositive(\s|$)/gim.exec(input),
    icon: browsers.netpositive,
  },
  {
    regex: (input: string) => !!/(^|\s)netscape(\s|$)/gim.exec(input),
    icon: browsers.netscape,
  },
  {
    regex: (input: string) => !!/(^|\s)opera(\s|$)/gim.exec(input),
    icon: browsers.opera,
  },
  {
    regex: (input: string) => !!/(^|\s)otter(\s|$)/gim.exec(input),
    icon: browsers.otter,
  },
  {
    regex: (input: string) => !!/(^|\s)planetweb(\s|$)/gim.exec(input),
    icon: browsers.planetweb,
  },
  {
    regex: (input: string) => !!/(^|\s)retrozilla(\s|$)/gim.exec(input),
    icon: browsers.retrozilla,
  },
  {
    regex: (input: string) => !!/(^|\s)safari(\s|$)/gim.exec(input),
    icon: browsers.safari,
  },
  {
    regex: (input: string) => !!/(^|\s)samsung(\s|$)/gim.exec(input),
    icon: browsers.samsung,
  },
  {
    regex: (input: string) => !!/(^|\s)seamonkey(\s|$)/gim.exec(input),
    icon: browsers.seamonkey,
  },
  {
    regex: (input: string) => !!/(^|\s)uc([\s]*browser)?(\s|$)/gim.exec(input),
    icon: browsers.uc,
  },
  {
    regex: (input: string) => !!/(^|\s)vivaldi(\s|$)/gim.exec(input),
    icon: browsers.vivaldi,
  },
  {
    regex: (input: string) => !!/(^|\s)webpositive(\s|$)/gim.exec(input),
    icon: browsers.webpositive,
  },
];
