import React, { FunctionComponent } from "react";
import { useCountries } from "../../api/analytics";
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

  return (
    <TableWidget {...props} data={data} filterPrimary="c" showSelfWhenEmpty />
  );
};

export const Countries: FunctionComponent<CountriesProps> = (props) => {
  return (
    <SuspenseCard>
      <CountriesContent {...props} />
    </SuspenseCard>
  );
};
