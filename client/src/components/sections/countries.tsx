import React, { FunctionComponent } from "react";
import Flag from "react-world-flags";
import { useCountries } from "../../api/analytics";
import { AnalyticsData } from "../../api/types";
import { COUNTRIES } from "../../constants/countries";
import { SuspenseCard } from "../shared/suspense-card";
import { TableWidget, TableWidgetWrapperProps } from "../shared/table-widget";

interface CountriesProps extends TableWidgetWrapperProps {
  domain: string;
}

const CountriesContent: FunctionComponent<CountriesProps> = ({
  domain,
  ...props
}) => {
  const { data } = useCountries(domain, props.filters);

  const renderIcon = (item: AnalyticsData) => (
    <Flag code={item.value} width={16} />
  );

  const formatName = (name: string) => {
    return COUNTRIES[name] || name;
  };

  return (
    <TableWidget
      {...props}
      title="Countries"
      data={data}
      renderIcon={renderIcon}
      formatName={formatName}
      filterPrimary="c"
      showSelfWhenEmpty
    />
  );
};

export const Countries: FunctionComponent<CountriesProps> = (props) => {
  return (
    <SuspenseCard>
      <CountriesContent {...props} />
    </SuspenseCard>
  );
};
